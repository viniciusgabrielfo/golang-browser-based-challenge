package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoteString(t *testing.T) {
	t.Parallel()

	someStockQuote := StockQuote{
		Symbol: "MSFT",
		Quote:  10.558,
	}

	assert.Equal(t, "MSFT quote is $10.56 per share", someStockQuote.String())
}
