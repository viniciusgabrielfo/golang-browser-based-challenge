package internal

import "fmt"

type StockQuote struct {
	Symbol string
	Quote  float64
}

func (s *StockQuote) String() string {
	return fmt.Sprintf("%s quote is $%.2f per share", s.Symbol, s.Quote)
}
