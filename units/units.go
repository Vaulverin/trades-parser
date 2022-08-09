package units

import (
	"github.com/shopspring/decimal"
	"math/big"
)

func ToDecimal(ivalue interface{}, decimals uint8) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

func ToUnits(value decimal.Decimal, decimals uint8) *big.Int {
	mul := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))
	return value.Mul(mul).BigInt()
}
