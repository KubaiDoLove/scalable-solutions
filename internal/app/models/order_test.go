package models

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewGoodTillCancelledOrder(t *testing.T) {
	type testCase struct {
		generalInfo   OrderGeneralInfo
		expectedOrder *Order
	}

	testValidUntil := time.Now().UTC().Add(time.Hour * 24 * 7)
	testCounterParty := "testCounterParty"

	testCases := []testCase{
		{
			generalInfo:   OrderGeneralInfo{},
			expectedOrder: &Order{},
		},
		{
			generalInfo: OrderGeneralInfo{
				ValidUntil:   &testValidUntil,
				Price:        decimal.NewFromInt(20),
				Quantity:     1,
				Operation:    Bid,
				CounterParty: testCounterParty,
			},
			expectedOrder: &Order{
				OrderGeneralInfo: OrderGeneralInfo{
					ValidUntil:   &testValidUntil,
					Price:        decimal.NewFromInt(20),
					Quantity:     1,
					Operation:    Bid,
					CounterParty: testCounterParty,
				},
			},
		},
		{
			generalInfo: OrderGeneralInfo{
				ValidUntil:   &testValidUntil,
				Price:        decimal.NewFromFloat32(5.5),
				Quantity:     3,
				Operation:    Ask,
				CounterParty: testCounterParty,
			},
			expectedOrder: &Order{
				OrderGeneralInfo: OrderGeneralInfo{
					ValidUntil:   &testValidUntil,
					Price:        decimal.NewFromFloat32(5.5),
					Quantity:     3,
					Operation:    Ask,
					CounterParty: testCounterParty,
				},
			},
		},
	}

	for _, testCase := range testCases {
		order := NewGoodTillCancelledOrder(testCase.generalInfo)
		assert.NotNil(t, order)

		assert.NotNil(t, order.ValidUntil)
		if testCase.generalInfo.ValidUntil == nil {
			defaultExpireTime := time.Now().UTC().Add(time.Hour * 24 * 90)
			defaultExpireYear, defaultExpireMonth, defaultExpireDay := defaultExpireTime.Date()

			expireYear, expireMonth, expireDay := order.ValidUntil.Date()
			assert.Equal(t, defaultExpireYear, expireYear)
			assert.Equal(t, defaultExpireMonth, expireMonth)
			assert.Equal(t, defaultExpireDay, expireDay)
		}
		assert.Equal(t, testCase.expectedOrder.Price, order.Price)
		assert.Equal(t, testCase.expectedOrder.Quantity, order.Quantity)
		assert.Equal(t, testCase.expectedOrder.Operation, order.Operation)
		assert.Equal(t, testCase.expectedOrder.CounterParty, order.CounterParty)
		assert.True(t, order.IsEnabled)

		assert.NotZero(t, order.TradeCode)
		assert.Equal(t, GoodTillCancelled, order.Type)
	}
}
