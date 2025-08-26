package valueobjects

import (
	"strings"
)

type TaxNumber struct {
	taxNumber string
}

func NewTaxNumber(taxNumber string) (TaxNumber, error) {
	taxNumber = strings.TrimSpace(taxNumber)
	return TaxNumber{taxNumber: taxNumber}, nil
}

func (t TaxNumber) Get() string {
	return t.taxNumber
}
