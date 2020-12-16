package inmemory

import (
	"context"
	"github.com/KubaiDoLove/scalable-solutions/internal/app/datastore"
	"github.com/KubaiDoLove/scalable-solutions/internal/app/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"sort"
	"sync"
)

// Thread safe "in memory" store
type Store struct {
	mu   sync.Mutex
	asks map[uuid.UUID]models.Order
	bids map[uuid.UUID]models.Order
}

func New() datastore.DataStore {
	return &Store{
		asks: make(map[uuid.UUID]models.Order, 0),
		bids: make(map[uuid.UUID]models.Order, 0),
	}
}

// CreateOrder validates order on a very basic level and saves it
func (s *Store) CreateOrder(ctx context.Context, order *models.Order) error {
	if order == nil {
		return datastore.ErrEmptyStruct
	}

	if order.ID == uuid.Nil {
		return datastore.ErrZeroID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if order.Operation == models.Ask {
		s.asks[order.ID] = *order
		return nil
	}

	s.bids[order.ID] = *order
	return nil
}

// DisableOrder to remove it from market snapshot and other reads
func (s *Store) DisableOrder(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return datastore.ErrZeroID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if order, orderInAsks := s.asks[id]; orderInAsks {
		if order.IsProcessable() {
			s.asks[id].OrderGeneralInfo.IsEnabled = false
		}
		return nil
	}

	if order, orderInBids := s.bids[id]; orderInBids {
		if order.IsProcessable() {
			s.bids[id].OrderGeneralInfo.IsEnabled = false
		}
		return nil
	}

	return datastore.ErrOrderDoesNotExist
}

// OrderByID returns only enabled and not expired order
func (s *Store) OrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	if id == uuid.Nil {
		return nil, datastore.ErrZeroID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if order, orderInAsks := s.asks[id]; orderInAsks {
		if order.IsProcessable() {
			return &order, nil
		}
	}

	if order, orderInBids := s.bids[id]; orderInBids {
		if order.IsProcessable() {
			return &order, nil
		}
	}

	return nil, datastore.ErrOrderDoesNotExist
}

// MatchOrder to get available bids/asks for a given order
func (s *Store) MatchOrder(ctx context.Context, order *models.Order) ([]models.Order, error) {
	if order == nil {
		return nil, datastore.ErrEmptyStruct
	}

	if order.Operation == models.Bid {
		return s.matchBid(order.Price)
	}

	return s.matchAsk(order.Price)
}

func (s *Store) matchBid(bidPrice decimal.Decimal) ([]models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	matchingAsks := make([]models.Order, 0)

	for _, ask := range s.asks {
		if ask.IsProcessable() && ask.Price.LessThanOrEqual(bidPrice) && ask.Quantity > 0 {
			matchingAsks = append(matchingAsks, ask)
		}
	}

	return matchingAsks, nil
}

func (s *Store) matchAsk(askPrice decimal.Decimal) ([]models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	matchingBids := make([]models.Order, 0)

	for _, bid := range s.bids {
		if bid.IsProcessable() && bid.Price.GreaterThanOrEqual(askPrice) && bid.Quantity > 0 {
			matchingBids = append(matchingBids, bid)
		}
	}

	return matchingBids, nil
}

// MarketDataSnapshot to get actual market data ordered by price
func (s *Store) MarketDataSnapshot(ctx context.Context) (*models.MarketDataSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	asks := make([]models.OrderSnapshot, 0)
	for _, ask := range s.asks {
		if ask.IsProcessable() {
			asks = append(asks, *ask.Snapshot())
		}
	}

	bids := make([]models.OrderSnapshot, 0)
	for _, bid := range s.bids {
		if bid.IsProcessable() {
			bids = append(bids, *bid.Snapshot())
		}
	}

	sort.SliceStable(asks, func(i, j int) bool { return asks[i].Price.LessThan(asks[j].Price) })
	sort.SliceStable(bids, func(i, j int) bool { return bids[i].Price.LessThan(bids[j].Price) })

	return &models.MarketDataSnapshot{
		Asks: asks,
		Bids: bids,
	}, nil
}
