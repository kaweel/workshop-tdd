package storage

import (
	"github.com/kaweel/workshop-tdd/payment/constant"
	"gorm.io/gorm"
)

type CustomerProfile struct {
	gorm.Model
	Name   string                  `gorm:"type:varchar(100);not null;"`
	Status constant.CustomerStatus `gorm:"type:varchar(10);not null;"`
	Amount float64                 `gorm:"not null"`
}
