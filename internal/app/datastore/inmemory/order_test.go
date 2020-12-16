package inmemory

import (
	"context"
	"github.com/KubaiDoLove/scalable-solutions/internal/app/datastore"
	"github.com/KubaiDoLove/scalable-solutions/internal/app/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	validStore := &Store{
		asks: make(map[uuid.UUID]models.Order, 0),
		bids: make(map[uuid.UUID]models.Order, 0),
	}

	assert.Equal(t, validStore, New())
}

func TestStore_CreateOrder(t *testing.T) {
	store := New()

	type errTestCase struct {
		order       *models.Order
		expectedErr error
	}

	errTestCases := []errTestCase{
		{
			expectedErr: datastore.ErrEmptyStruct,
		},
		{
			order: &models.Order{
				OrderGeneralInfo: &models.OrderGeneralInfo{},
			},
			expectedErr: datastore.ErrZeroID,
		},
	}

	for _, testCase := range errTestCases {
		assert.Equal(t, testCase.expectedErr, store.CreateOrder(context.Background(), testCase.order))
	}

	testBid, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(32),
		Quantity:     5,
		Operation:    models.Bid,
		CounterParty: "BidCounterParty",
	})

	err := store.CreateOrder(context.Background(), testBid)
	assert.Nil(t, err)

	bidFromStore, _ := store.OrderByID(context.Background(), testBid.ID)
	assert.Equal(t, testBid, bidFromStore)

	testAsk, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(21),
		Quantity:     3,
		Operation:    models.Ask,
		CounterParty: "AskCounterParty",
	})

	err = store.CreateOrder(context.Background(), testAsk)
	assert.Nil(t, err)

	askFromStore, err := store.OrderByID(context.Background(), testAsk.ID)
	assert.Equal(t, testAsk, askFromStore)
}

func TestStore_DisableOrder(t *testing.T) {
	store := New()

	assert.Equal(t, datastore.ErrZeroID, store.DisableOrder(context.Background(), uuid.Nil))

	testBid, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(32),
		Quantity:     5,
		Operation:    models.Bid,
		CounterParty: "BidCounterParty",
	})
	_ = store.CreateOrder(context.Background(), testBid)

	err := store.DisableOrder(context.Background(), testBid.ID)
	assert.Nil(t, err)

	_, err = store.OrderByID(context.Background(), testBid.ID)
	assert.Equal(t, datastore.ErrOrderDoesNotExist, err)

	testAsk, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(21),
		Quantity:     3,
		Operation:    models.Ask,
		CounterParty: "AskCounterParty",
	})
	_ = store.CreateOrder(context.Background(), testAsk)

	err = store.DisableOrder(context.Background(), testAsk.ID)
	assert.Nil(t, err)

	_, err = store.OrderByID(context.Background(), testAsk.ID)
	assert.Equal(t, datastore.ErrOrderDoesNotExist, err)
}

func TestStore_OrderByID(t *testing.T) {
	store := New()

	_, err := store.OrderByID(context.Background(), uuid.Nil)
	assert.Equal(t, datastore.ErrZeroID, err)

	testBidOne, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(1),
		Quantity:     1,
		Operation:    models.Bid,
		CounterParty: "testBidOne",
	})
	_ = store.CreateOrder(context.Background(), testBidOne)
	testBidTwo, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(2),
		Quantity:     2,
		Operation:    models.Bid,
		CounterParty: "testBidTwo",
	})
	_ = store.CreateOrder(context.Background(), testBidTwo)

	bidOne, err := store.OrderByID(context.Background(), testBidOne.ID)
	assert.Nil(t, err)
	assert.Equal(t, testBidOne, bidOne)
	bidTwo, err := store.OrderByID(context.Background(), testBidTwo.ID)
	assert.Nil(t, err)
	assert.Equal(t, testBidTwo, bidTwo)

	testAskOne, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(1),
		Quantity:     1,
		Operation:    models.Ask,
		CounterParty: "testAskOne",
	})
	_ = store.CreateOrder(context.Background(), testAskOne)
	testAskTwo, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(2),
		Quantity:     2,
		Operation:    models.Ask,
		CounterParty: "testAskTwo",
	})
	_ = store.CreateOrder(context.Background(), testAskTwo)

	askOne, err := store.OrderByID(context.Background(), testAskOne.ID)
	assert.Nil(t, err)
	assert.Equal(t, testAskOne, askOne)
	askTwo, err := store.OrderByID(context.Background(), testAskTwo.ID)
	assert.Nil(t, err)
	assert.Equal(t, testAskTwo, askTwo)
}

// We automatically test matchBid and matchAsk functions when we test MatchOrder
func TestStore_MatchOrder(t *testing.T) {
	store := New()

	_, err := store.MatchOrder(context.Background(), nil)
	assert.Equal(t, datastore.ErrEmptyStruct, err)

	validAskOne, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(100),
		Quantity:     10,
		Operation:    models.Ask,
		CounterParty: "validAskOne",
	})
	_ = store.CreateOrder(context.Background(), validAskOne)
	validAskTwo, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(120),
		Quantity:     100,
		Operation:    models.Ask,
		CounterParty: "validAskTwo",
	})
	_ = store.CreateOrder(context.Background(), validAskTwo)
	invalidAskOne, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(999),
		Quantity:     2,
		Operation:    models.Ask,
		CounterParty: "invalidAskOne",
	})
	_ = store.CreateOrder(context.Background(), invalidAskOne)
	invalidAskTwo, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(1),
		Quantity:     0,
		Operation:    models.Ask,
		CounterParty: "invalidAskTwo",
	})
	_ = store.CreateOrder(context.Background(), invalidAskTwo)

	testBid := &models.Order{
		OrderGeneralInfo: &models.OrderGeneralInfo{
			Price:     decimal.NewFromInt(150),
			Operation: models.Bid,
		},
	}

	matchingAsks, err := store.MatchOrder(context.Background(), testBid)
	assert.Nil(t, err)
	assert.Len(t, matchingAsks, 2)

	matchingAsksIDs := make([]uuid.UUID, 0, len(matchingAsks))
	for _, ask := range matchingAsks {
		assert.Equal(t, models.Ask, ask.Operation)
		matchingAsksIDs = append(matchingAsksIDs, ask.ID)
	}
	assert.Contains(t, matchingAsksIDs, validAskOne.ID)
	assert.Contains(t, matchingAsksIDs, validAskTwo.ID)
	assert.NotContains(t, matchingAsksIDs, invalidAskOne.ID)
	assert.NotContains(t, matchingAsksIDs, invalidAskTwo.ID)

	validBid, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(10),
		Quantity:     1,
		Operation:    models.Bid,
		CounterParty: "validBid",
	})
	_ = store.CreateOrder(context.Background(), validBid)
	invalidBidOne, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(1),
		Quantity:     2,
		Operation:    models.Bid,
		CounterParty: "invalidBidOne",
	})
	_ = store.CreateOrder(context.Background(), invalidBidOne)
	invalidBidTwo, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(12),
		Quantity:     0,
		Operation:    models.Bid,
		CounterParty: "invalidBidTwo",
	})
	_ = store.CreateOrder(context.Background(), invalidBidTwo)

	testAsk := &models.Order{
		OrderGeneralInfo: &models.OrderGeneralInfo{
			Price:     decimal.NewFromInt(5),
			Operation: models.Ask,
		},
	}

	matchingBids, err := store.MatchOrder(context.Background(), testAsk)
	assert.Nil(t, err)
	assert.Len(t, matchingBids, 1)

	matchingBidsIDs := make([]uuid.UUID, 0, len(matchingBids))
	for _, bid := range matchingBids {
		assert.Equal(t, models.Bid, bid.Operation)
		matchingBidsIDs = append(matchingBidsIDs, bid.ID)
	}
	assert.Contains(t, matchingBidsIDs, validBid.ID)
	assert.NotContains(t, matchingBidsIDs, invalidBidOne.ID)
	assert.NotContains(t, matchingBidsIDs, invalidBidTwo.ID)
}

func TestStore_MarketDataSnapshot(t *testing.T) {
	store := New()
	notValidDate := time.Now().UTC().Add(time.Hour * -8)

	validAskOne, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(100),
		Quantity:     10,
		Operation:    models.Ask,
		CounterParty: "validAskOne",
	})
	_ = store.CreateOrder(context.Background(), validAskOne)
	validAskTwo, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(120),
		Quantity:     100,
		Operation:    models.Ask,
		CounterParty: "validAskTwo",
	})
	_ = store.CreateOrder(context.Background(), validAskTwo)
	invalidAskOne, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(999),
		Quantity:     2,
		Operation:    models.Ask,
		CounterParty: "invalidAskOne",
		ValidUntil:   &notValidDate,
	})
	_ = store.CreateOrder(context.Background(), invalidAskOne)

	validBidOne, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(10),
		Quantity:     1,
		Operation:    models.Bid,
		CounterParty: "validBidOne",
	})
	_ = store.CreateOrder(context.Background(), validBidOne)
	validBidTwo, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(12),
		Quantity:     1,
		Operation:    models.Bid,
		CounterParty: "validBidTwo",
	})
	_ = store.CreateOrder(context.Background(), validBidTwo)
	invalidBidOne, _ := models.NewGoodTillCancelledOrder(&models.OrderGeneralInfo{
		TradeCode:    uuid.New(),
		Price:        decimal.NewFromInt(1),
		Quantity:     2,
		Operation:    models.Bid,
		CounterParty: "invalidBidOne",
		ValidUntil:   &notValidDate,
	})
	_ = store.CreateOrder(context.Background(), invalidBidOne)

	marketData, err := store.MarketDataSnapshot(context.Background())
	assert.Nil(t, err)
	assert.Len(t, marketData.Asks, 2)
	assert.Len(t, marketData.Bids, 2)

	// Checks asks price sorting
	for i := 0; i < len(marketData.Asks)-1; i++ {
		assert.True(t, marketData.Asks[i].Price.LessThanOrEqual(marketData.Asks[i+1].Price))
	}
	// Checks bids price sorting
	for i := 0; i < len(marketData.Bids)-1; i++ {
		assert.True(t, marketData.Bids[i].Price.LessThanOrEqual(marketData.Bids[i+1].Price))
	}
}
