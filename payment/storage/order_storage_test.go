//go:build integration_test
// +build integration_test

package storage

import (
	"context"
	"time"

	"testing"

	"github.com/kaweel/workshop-tdd/payment/clock"
	"github.com/kaweel/workshop-tdd/payment/constant"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mssql"
	"gorm.io/gorm"
)

func TestOrderStorage(t *testing.T) {
	var ctx context.Context
	var ot OrderStorage
	var c CustomerProfile
	var m MerchantProfile
	var o *Order
	var container *mssql.MSSQLServerContainer
	var db *gorm.DB
	var cl clock.Clock
	var tn time.Time

	setup := func() {
		cl = clock.NewClock()
		tn = cl.Now()
		ctx = context.Background()
		container, db = SetupMSSQL(ctx, t)
		db.Debug().AutoMigrate(&CustomerProfile{}, &MerchantProfile{}, &Order{})
		ot = NewOrderStorage(db)
		c = CustomerProfile{
			Model: gorm.Model{
				CreatedAt: tn,
			},
			Name:   "Madmax Drinkcola",
			Status: constant.CustomerStatusActive,
			Amount: 1000,
		}
		m = MerchantProfile{
			Model: gorm.Model{
				CreatedAt: tn,
			},
			Name:   "Rabit Cart",
			Status: constant.MerchantStatusActive,
			Amount: 1100,
		}
		o = &Order{
			Model: gorm.Model{
				CreatedAt: tn,
			},
			Customer: c,
			Merchant: m,
			Amount:   1200,
		}
		err := ot.Save(o)
		if err != nil {
			t.Fatalf("Failed to setup data [%v]", err.Error())
		}
	}

	cleanup := func() {
		defer CleanUpMSSQL(container, ctx, t)
	}

	t.Run("get order by id not found should return error", func(t *testing.T) {
		//Arrange
		setup()
		defer cleanup()

		//Action
		_, err := ot.GetOrder(2)

		//Assert
		assert.Equal(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("get order found should return order", func(t *testing.T) {
		//Arrange
		setup()
		defer cleanup()

		//Action
		expected, err := ot.GetOrder(o.ID)

		//Assert
		assert.Nil(t, err)
		assert.Equal(t, expected.ID, o.ID)
		assert.Equal(t, expected.MerchantID, o.MerchantID)
		assert.Equal(t, expected.Amount, o.Amount)
		assert.Equal(t, expected.Status, o.Status)
		assert.Equal(t, expected.CreatedAt.In(time.UTC), o.CreatedAt)

		assert.Equal(t, expected.CustomerID, o.CustomerID)
		assert.Equal(t, expected.Customer.Name, o.Customer.Name)
		assert.Equal(t, expected.Customer.Amount, o.Customer.Amount)
		assert.Equal(t, expected.Customer.Status, o.Customer.Status)

		assert.Equal(t, expected.MerchantID, o.MerchantID)
		assert.Equal(t, expected.Merchant.Name, o.Merchant.Name)
		assert.Equal(t, expected.Merchant.Amount, o.Merchant.Amount)
		assert.Equal(t, expected.Merchant.Status, o.Merchant.Status)
	})
}
