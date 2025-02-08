//go:build unit_test
// +build unit_test

package service

import (
	"errors"
	"testing"
	"time"

	"github.com/kaweel/workshop-tdd/payment/constant"
	"github.com/kaweel/workshop-tdd/payment/messaging"
	"github.com/kaweel/workshop-tdd/payment/storage"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockOrderStorage struct {
	o   *storage.Order
	err error
}

func (m *mockOrderStorage) SetOrder(o *storage.Order, err error) {
	m.o = o
	m.err = err
}

func (m *mockOrderStorage) GetOrder(id uint) (*storage.Order, error) {
	return m.o, m.err
}

func (m *mockOrderStorage) Save(o *storage.Order) error {
	return m.err
}

type mockPaymentTranasctionStorage struct {
	Calls []*storage.PaymentTranasction
	err   error
}

func (m *mockPaymentTranasctionStorage) SetSave(err error) {
	m.err = err
}

func (m *mockPaymentTranasctionStorage) Save(o *storage.PaymentTranasction) error {
	m.Calls = append(m.Calls, o)
	return m.err
}

type mockKafkaProducer struct {
	Calls []messaging.RequestPublish
	err   error
}

func (m *mockKafkaProducer) SetPublish(err error) {
	m.err = err
}

func (m *mockKafkaProducer) Publish(r messaging.RequestPublish) error {
	m.Calls = append(m.Calls, r)
	return m.err
}

type mockClock struct {
	t time.Time
}

func (m *mockClock) SetNow(t time.Time) {
	m.t = t
}

func (m *mockClock) Now() time.Time {
	return m.t
}

func TestPaymentService(t *testing.T) {
	var s Service
	var m *mockOrderStorage
	var mp *mockPaymentTranasctionStorage
	var mk *mockKafkaProducer
	var mt *mockClock
	var o *storage.Order
	var err error
	var prr error
	var krr error
	var r RequestPayment
	var pm PaymentMessage

	setup := func() {
		m = &mockOrderStorage{}
		mp = &mockPaymentTranasctionStorage{}
		mk = &mockKafkaProducer{}
		mt = &mockClock{}
		o = &storage.Order{
			Model: gorm.Model{
				ID: 1,
			},
			Amount:     100,
			Status:     constant.OrderStatusRequestPayment,
			CustomerID: 1,
			Customer: storage.CustomerProfile{
				Model: gorm.Model{
					ID: 1,
				},
				Status: constant.CustomerStatusActive,
				Amount: 1000,
			},
			MerchantID: 1,
			Merchant: storage.MerchantProfile{
				Model: gorm.Model{
					ID: 1,
				},
				Status: constant.MerchantStatusActive,
			},
		}
		err = nil
		prr = nil
		krr = nil
		m.SetOrder(o, err)
		mp.SetSave(prr)
		mk.SetPublish(krr)
		mt.SetNow(time.Now())
		s = NewService(m, mp, mk, mt)
		r = RequestPayment{
			OrderID: 1,
			Channel: constant.PaymentChannelDebit,
		}
		pm = PaymentMessage{
			OrderID:   r.OrderID,
			Amount:    r.Amount,
			CreatedAt: mt.Now(),
			Status:    constant.PaymentTranasctionStatusConfirm,
		}
	}

	t.Run("payment channel is not support should reject transaction and publish reject event", func(t *testing.T) {
		//Arrange
		setup()
		r.Channel = "Zebit"
		pm.Status = constant.PaymentTranasctionStatusReject
		pm.Reason = "invalid payment channel"

		//Action
		actual := s.Payment(r)

		//Assert
		assertTransactionRejected(t, pm, actual, mp, mk)
	})

	t.Run("order not found should reject transaction and publish reject event", func(t *testing.T) {
		//Arrange
		setup()
		pm.Status = constant.PaymentTranasctionStatusReject
		pm.Reason = "order status not found"
		err = errors.New("order status not found")
		m.SetOrder(nil, err)

		//Action
		actual := s.Payment(r)

		//Assert
		assertTransactionRejected(t, pm, actual, mp, mk)
	})

	t.Run("order status is not request payment should reject transaction and publish reject event", func(t *testing.T) {
		//Arrange
		setup()
		pm.Status = constant.PaymentTranasctionStatusReject
		pm.Reason = "order status is not request payment"
		o.Status = constant.OrderStatusOpen
		m.SetOrder(o, nil)

		//Action
		actual := s.Payment(r)

		//Assert
		assertTransactionRejected(t, pm, actual, mp, mk)
	})

	t.Run("customer status is not active should reject transaction and publish reject event", func(t *testing.T) {
		//Arrange
		setup()
		pm.Status = constant.PaymentTranasctionStatusReject
		pm.Reason = "customer status is not active"
		o.Customer.Status = constant.CustomerStatusInActive
		m.SetOrder(o, nil)

		//Action
		actual := s.Payment(r)

		//Assert
		assertTransactionRejected(t, pm, actual, mp, mk)
	})

	t.Run("customer amount is not enough should reject transaction and publish reject event", func(t *testing.T) {
		//Arrange
		setup()
		pm.Status = constant.PaymentTranasctionStatusReject
		pm.Reason = "customer amount is not enough"
		o.Customer.Amount = 0
		m.SetOrder(o, nil)

		//Action
		actual := s.Payment(r)

		//Assert
		assertTransactionRejected(t, pm, actual, mp, mk)
	})

	t.Run("merchant status is not active should reject transaction and publish reject event", func(t *testing.T) {
		//Arrange
		setup()
		pm.Status = constant.PaymentTranasctionStatusReject
		pm.Reason = "merchant status is not active"
		o.Merchant.Status = constant.MerchantStatusInActive
		m.SetOrder(o, nil)

		//Action
		actual := s.Payment(r)

		//Assert
		assertTransactionRejected(t, pm, actual, mp, mk)
	})

	t.Run("make transaction fail should not publish reject event", func(t *testing.T) {
		//Arrange
		setup()
		prr = errors.New("unknown error")
		mp.SetSave(prr)

		//Action
		expected := s.Payment(r)

		//Assert
		assert.EqualError(t, expected, "unknown error")
	})

	t.Run("publish transaction fail should not publish reject event", func(t *testing.T) {
		//Arrange
		setup()
		krr = errors.New("publish complete payment transaction fail")
		mk.SetPublish(krr)
		//Action
		expected := s.Payment(r)

		//Assert
		assert.EqualError(t, expected, "publish complete payment transaction fail")
	})

	t.Run("success transaction should publish completed event", func(t *testing.T) {
		//Arrange
		setup()
		mk.SetPublish(krr)
		pt := &storage.PaymentTranasction{
			Model: gorm.Model{
				UpdatedAt: mt.t,
			},
			OrderID: r.OrderID,
			Amount:  r.Amount,
			Channel: r.Channel,
			Status:  constant.PaymentTranasctionStatusConfirm,
		}

		//Action
		expected := s.Payment(r)

		//Assert save txn
		assert.Nil(t, expected)
		assert.Equal(t, 1, len(mp.Calls))
		assert.Equal(t, pt, mp.Calls[0])

		//Assert publish msg
		assert.Equal(t, 1, len(mk.Calls))
		assert.Equal(t, constant.KafkaTopicPaymentTransaction, mk.Calls[0].Topic)
		assert.Equal(t, "1", mk.Calls[0].Key)
		assert.Equal(t, pm, mk.Calls[0].Message)
	})
}

func assertTransactionRejected(t *testing.T, pm PaymentMessage, actual error, mp *mockPaymentTranasctionStorage, mk *mockKafkaProducer) {
	assert.Equal(t, pm.Reason, actual.Error())
	assert.Equal(t, pm.Status, mp.Calls[0].Status)
	assert.Equal(t, pm.Reason, mp.Calls[0].Reason)
	assert.Equal(t, pm.Status, mk.Calls[0].Message.(PaymentMessage).Status)
	assert.Equal(t, pm.Reason, mk.Calls[0].Message.(PaymentMessage).Reason)
}
