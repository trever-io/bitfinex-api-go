package rest

import (
	"net/url"
	"path"
	"strconv"

	"github.com/trever-io/bitfinex-api-go/pkg/convert"
	"github.com/trever-io/bitfinex-api-go/pkg/models/book"
	"github.com/trever-io/bitfinex-api-go/pkg/models/common"
)

type BitfinexBookService interface {
	All(symbol string, precision common.BookPrecision, priceLevels int) (*book.Snapshot, error)
}

type BookService struct {
	Synchronous
}

// All - retrieve all books for the given symbol with the given precision at the given price level
// see https://docs.bitfinex.com/reference#rest-public-books for more info
func (b *BookService) All(symbol string, precision common.BookPrecision, priceLevels int) (*book.Snapshot, error) {
	req := NewRequestWithMethod(path.Join("book", symbol, string(precision)), "GET")
	req.Params = make(url.Values)
	req.Params.Add("len", strconv.Itoa(priceLevels))

	raw, err := b.Request(req)
	if err != nil {
		return nil, err
	}

	return book.SnapshotFromRaw(symbol, string(precision), convert.ToInterfaceArray(raw), raw)
}
