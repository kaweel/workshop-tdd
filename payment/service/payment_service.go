package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/kaweel/workshop-tdd/payment/clock"
	"github.com/kaweel/workshop-tdd/payment/constant"
	"github.com/kaweel/workshop-tdd/payment/messaging"
	"gorm.io/gorm"

	"github.com/kaweel/workshop-tdd/payment/storage"
)

type RequestPayment struct {
	OrderID uint                    `json:"orderID"`
	Channel constant.PaymentChannel `json:"channel"`
	Amount  float64                 `json:"amount"`
}

type Service interface {
	Payment(r RequestPayment) error
}

type service struct {
	o storage.OrderStorage
	p storage.PaymentTranasctionStorage
	m messaging.KafkaProducer
	c clock.Clock
}

func NewService(o storage.OrderStorage, p storage.PaymentTranasctionStorage, m messaging.KafkaProducer, c clock.Clock) Service {
	return &service{
		o: o,
		p: p,
		m: m,
		c: c,
	}
}

type PaymentMessage struct {
	OrderID   uint                              `json:"orderID"`
	Status    constant.PaymentTranasctionStatus `json:"status"`
	Amount    float64                           `json:"amount"`
	Reason    string                            `json:"reason"`
	CreatedAt time.Time                         `json:"createdAt"`
}

func validateOrderPayment(r RequestPayment, getOrderByID func(id uint) (*storage.Order, error)) error {
	v := constant.IsValidPaymentChannel(r.Channel)
	if !v {
		return errors.New("invalid payment channel")
	}
	o, err := getOrderByID(r.OrderID)
	if err != nil {
		return err
	}
	v = constant.IsOrderRequestPayment(o.Status)
	if !v {
		return errors.New("order status is not request payment")
	}
	v = constant.IsActiveCustomer(o.Customer.Status)
	if !v {
		return errors.New("customer status is not active")
	}
	if o.Customer.Amount < o.Amount {
		return errors.New("customer amount is not enough")
	}
	v = constant.IsActiveMerchant(o.Merchant.Status)
	if !v {
		return errors.New("merchant status is not active")
	}
	return nil
}

func (s *service) Payment(r RequestPayment) error {
	n := s.c.Now()
	t := &storage.PaymentTranasction{
		Model: gorm.Model{
			UpdatedAt: n,
		},
		OrderID: r.OrderID,
		Amount:  r.Amount,
		Channel: r.Channel,
		Status:  constant.PaymentTranasctionStatusConfirm,
	}

	l := PaymentMessage{
		OrderID:   r.OrderID,
		Amount:    r.Amount,
		CreatedAt: n,
		Status:    constant.PaymentTranasctionStatusConfirm,
	}

	u := messaging.RequestPublish{
		Topic:   constant.KafkaTopicPaymentTransaction,
		Key:     strconv.FormatUint(uint64(r.OrderID), 10),
		Message: l,
	}

	validateOrderErr := validateOrderPayment(r, s.o.GetOrder)
	if validateOrderErr != nil {
		t.Status = constant.PaymentTranasctionStatusReject
		t.Reason = validateOrderErr.Error()
		l.Status = t.Status
		l.Reason = t.Reason
		u.Message = l
	}

	err := s.p.Save(t)
	if err != nil {
		return err
	}

	if err = s.m.Publish(u); err != nil {
		return err
	}

	return validateOrderErr
}
