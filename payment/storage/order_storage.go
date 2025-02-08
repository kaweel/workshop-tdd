package storage

import (
	"github.com/kaweel/workshop-tdd/payment/constant"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID uint                 `gorm:"not null"` // Foreign Key to CustomerProfile
	MerchantID uint                 `gorm:"not null"` // Foreign Key to MerchantProfile
	Amount     float64              `gorm:"not null"`
	Status     constant.OrderStatus `gorm:"type:varchar(20);not null;"`

	// Relations
	Customer CustomerProfile `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;saveAssociation:true"`
	Merchant MerchantProfile `gorm:"foreignKey:MerchantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;saveAssociation:true"`
}

type OrderStorage interface {
	GetOrder(id uint) (*Order, error)
	Save(*Order) error
}

type orderStorage struct {
	db *gorm.DB
}

func NewOrderStorage(db *gorm.DB) OrderStorage {
	return &orderStorage{
		db: db,
	}
}

func (s *orderStorage) Save(o *Order) error {
	r := s.db.Debug().Save(o)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *orderStorage) GetOrder(id uint) (*Order, error) {
	o := &Order{}
	r := s.db.Debug().Preload("Customer").Preload("Merchant").Where("ID = ?", id).First(o)
	if r.Error != nil {
		return nil, r.Error
	}
	return o, nil
}
