package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewGoodTillCancelledOrder(t *testing.T) {
	type testCase struct {
		generalInfo   *OrderGeneralInfo
		expectedOrder *Order
	}

	testValidUntil := time.Now().UTC().Add(time.Hour * 24 * 7)
	testCounterParty := "testCounterParty"
	testTradeCode := uuid.New()

	testCases := []testCase{
		{},
		{
			generalInfo: &OrderGeneralInfo{
				TradeCode: testTradeCode,
			},
			expectedOrder: &Order{
				OrderGeneralInfo: &OrderGeneralInfo{
					TradeCode: testTradeCode,
				},
			},
		},
		{
			generalInfo: &OrderGeneralInfo{
				TradeCode:    testTradeCode,
				ValidUntil:   &testValidUntil,
				Price:        decimal.NewFromInt(20),
				Quantity:     1,
				Operation:    Bid,
				CounterParty: testCounterParty,
			},
			expectedOrder: &Order{
				OrderGeneralInfo: &OrderGeneralInfo{
					TradeCode:    testTradeCode,
					ValidUntil:   &testValidUntil,
					Price:        decimal.NewFromInt(20),
					Quantity:     1,
					Operation:    Bid,
					CounterParty: testCounterParty,
				},
			},
		},
		{
			generalInfo: &OrderGeneralInfo{
				TradeCode:    testTradeCode,
				ValidUntil:   &testValidUntil,
				Price:        decimal.NewFromFloat32(5.5),
				Quantity:     3,
				Operation:    Ask,
				CounterParty: testCounterParty,
			},
			expectedOrder: &Order{
				OrderGeneralInfo: &OrderGeneralInfo{
					TradeCode:    testTradeCode,
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
		order, err := NewGoodTillCancelledOrder(testCase.generalInfo)
		if err != nil {
			assert.EqualError(t, err, ErrNoEmptyGeneralInfo.Error())
			continue
		}

		assert.NotNil(t, order)
		assert.NotZero(t, order.ID)
		assert.Equal(t, testCase.expectedOrder.TradeCode, order.TradeCode)
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
		assert.Equal(t, GoodTillCancelled, order.Type)
	}
}

func TestOrder_IsProcessable(t *testing.T) {
	type testCase struct {
		order    Order
		expected bool
	}

	notValidDate := time.Now().UTC().Add(time.Hour * -8)
	validDate := time.Now().UTC().Add(time.Hour * 8)
	testCases := []testCase{
		{
			order: Order{
				OrderGeneralInfo: &OrderGeneralInfo{},
			},
			expected: false,
		},
		{
			order: Order{
				OrderGeneralInfo: &OrderGeneralInfo{
					ValidUntil: &notValidDate,
					IsEnabled:  true,
				},
			},
			expected: false,
		},
		{
			order: Order{
				OrderGeneralInfo: &OrderGeneralInfo{
					ValidUntil: &validDate,
					IsEnabled:  false,
				},
			},
			expected: false,
		},
		{
			order: Order{
				OrderGeneralInfo: &OrderGeneralInfo{
					ValidUntil: &validDate,
					IsEnabled:  true,
				},
			},
			expected: true,
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, testCase.order.IsProcessable())
	}
}

func TestOrder_Snapshot(t *testing.T) {
	type testCase struct {
		order            Order
		expectedSnapshot *OrderSnapshot
	}

	testCases := []testCase{
		{
			order: Order{
				OrderGeneralInfo: &OrderGeneralInfo{},
			},
			expectedSnapshot: &OrderSnapshot{},
		},
		{
			order: Order{
				OrderGeneralInfo: &OrderGeneralInfo{
					Price:    decimal.NewFromInt(2),
					Quantity: 3,
				},
			},
			expectedSnapshot: &OrderSnapshot{
				Price:    decimal.NewFromInt(2),
				Quantity: 3,
			},
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectedSnapshot, testCase.order.Snapshot())
	}
}
