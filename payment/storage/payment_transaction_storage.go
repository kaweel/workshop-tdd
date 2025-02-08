package storage

import (
	"github.com/kaweel/workshop-tdd/payment/constant"
	"gorm.io/gorm"
)

type PaymentTranasction struct {
	gorm.Model
	OrderID uint                              `gorm:"not null"`
	Channel constant.PaymentChannel           `gorm:"type:varchar(10);not null;"`
	Amount  float64                           `gorm:"not null;"`
	Status  constant.PaymentTranasctionStatus `gorm:"type:varchar(30);not null;"`
	Reason  string                            `gorm:"type:varchar(255);"`

	// Relation
	Order Order `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
type PaymentTranasctionStorage interface {
	Save(p *PaymentTranasction) error
}

type paymentTranasctionStorage struct {
	db *gorm.DB
}

func NewPaymentTranasctionStorage(db *gorm.DB) PaymentTranasctionStorage {
	return &paymentTranasctionStorage{
		db: db,
	}
}

func (s *paymentTranasctionStorage) Save(p *PaymentTranasction) error {
	r := s.db.Debug().Save(p)
	if r.Error != nil {
		return r.Error
	}
	return nil
}
