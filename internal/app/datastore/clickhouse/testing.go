package clickhouse

import (
	"testing"
)

func TestDB(t *testing.T) (*OrderBook, func()) {
	t.Helper()

	store, err := New()
	if err != nil {
		t.Fatal("no db connection: ", err)
	}
	orderBook := store.(*OrderBook)

	return orderBook, func() {
		if _, err := orderBook.db.Exec("TRUNCATE TABLE IF EXISTS orders"); err != nil {
			t.Fatal(err)
		}

		if err := orderBook.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
