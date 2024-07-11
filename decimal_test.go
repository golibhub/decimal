package decimal_test

import (
	"encoding/json"
	"testing"

	"github.com/golibhub/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDecimal(t *testing.T) {
	t.Parallel()

	t.Run("Scan", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		assert.Equal(t, "1.23", d.String())

		err = d.Scan("4.56")
		assert.NoError(t, err)
		assert.Equal(t, "4.56", d.String())
	})

	t.Run("Value", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		v, err := d.Value()
		assert.NoError(t, err)
		assert.Equal(t, "1.23", v)
	})

	t.Run("MarshalText", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		v, err := d.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, "1.23", string(v))
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		d := decimal.FromInt(0)
		err := d.UnmarshalText([]byte("1.23"))
		assert.NoError(t, err)
		assert.Equal(t, "1.23", d.String())
	})

	t.Run("String", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		assert.Equal(t, "1.23", d.String())
	})

	t.Run("IsNegative", func(t *testing.T) {
		d, err := decimal.FromStr("-1.23")
		assert.NoError(t, err)
		assert.True(t, d.IsNegative())
	})

	t.Run("IsZero", func(t *testing.T) {
		d, err := decimal.FromStr("0")
		assert.NoError(t, err)
		assert.True(t, d.IsZero())
	})

	t.Run("Equal", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		d2, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		assert.True(t, d.Equal(d2))
	})

	t.Run("Cmp - Equal", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		d2, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		assert.Equal(t, 0, d.Cmp(d2))
	})

	t.Run("Cmp - Less", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		d2, err := decimal.FromStr("1.24")
		assert.NoError(t, err)
		assert.Equal(t, -1, d.Cmp(d2))
	})

	t.Run("Cmp - Greater", func(t *testing.T) {
		d, err := decimal.FromStr("1.24")
		assert.NoError(t, err)
		d2, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		assert.Equal(t, 1, d.Cmp(d2))
	})

	t.Run("Add", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		d2, err := decimal.FromStr("1.24")
		assert.NoError(t, err)
		d3 := d.Add(d2)
		assert.Equal(t, "2.47", d3.String())
	})

	t.Run("Sub", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		d2, err := decimal.FromStr("1.24")
		assert.NoError(t, err)
		d3 := d.Sub(d2)
		assert.Equal(t, "-0.01", d3.String())
	})

	t.Run("Mul", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		d2, err := decimal.FromStr("1.24")
		assert.NoError(t, err)
		d3 := d.Mul(d2)
		assert.Equal(t, "1.5252", d3.String())
	})

	t.Run("Div", func(t *testing.T) {
		d, err := decimal.FromStr("1.23")
		assert.NoError(t, err)
		d2, err := decimal.FromStr("1.24")
		assert.NoError(t, err)
		d3 := d.Div(d2)
		assert.Equal(t, "0.99193548", d3.Round(8).String())
	})

	t.Run("Round", func(t *testing.T) {
		d, err := decimal.FromStr("1.23456789")
		assert.NoError(t, err)
		assert.Equal(t, "1.23", d.Round(2).String())
		assert.Equal(t, "1.235", d.Round(3).String())
		assert.Equal(t, "1.2346", d.Round(4).String())
		assert.Equal(t, "1.23457", d.Round(5).String())
		assert.Equal(t, "1.234568", d.Round(6).String())
		assert.Equal(t, "1.2345679", d.Round(7).String())
		assert.Equal(t, "1.23456789", d.Round(8).String())
		assert.Equal(t, "1.23456789", d.Round(9).String())
	})

	t.Run("Round - no digits after period", func(t *testing.T) {
		d := decimal.MustFromStr("1")

		assert.Equal(t, "1", d.Round(10).String())
	})

	t.Run("RoundOrNil", func(t *testing.T) {
		var dNil *decimal.Decimal
		d, err := decimal.FromStr("1.23456789")
		assert.NoError(t, err)

		assert.Equal(t, dNil, dNil.RoundOrNil(2))
		assert.Equal(t, "1.23", d.RoundOrNil(2).String())
		assert.Equal(t, "1.235", d.RoundOrNil(3).String())
		assert.Equal(t, "1.2346", d.RoundOrNil(4).String())
		assert.Equal(t, "1.23457", d.RoundOrNil(5).String())
		assert.Equal(t, "1.234568", d.RoundOrNil(6).String())
		assert.Equal(t, "1.2345679", d.RoundOrNil(7).String())
		assert.Equal(t, "1.23456789", d.RoundOrNil(8).String())
		assert.Equal(t, "1.23456789", d.RoundOrNil(9).String())
	})

	t.Run("InRangeInt", func(t *testing.T) {
		d, err := decimal.FromStr("1")
		assert.NoError(t, err)
		assert.True(t, d.InRangeInt(1, 10))
		assert.True(t, d.InRangeInt(1, 1))
		assert.True(t, d.InRangeInt(-10, 1))
		assert.True(t, d.InRangeInt(-10, 10))
		assert.False(t, d.InRangeInt(-10, -1))
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var testCases = []struct {
			text     string
			expected decimal.Decimal
		}{
			{
				text:     `"1.23"`,
				expected: decimal.MustFromStr("1.23"),
			},
			{
				text:     `1`,
				expected: decimal.MustFromStr("1"),
			},
			{
				text:     `1.123456789`,
				expected: decimal.MustFromStr("1.123456789"),
			},
		}

		for _, testCase := range testCases {
			var actual decimal.Decimal
			err := json.Unmarshal([]byte(testCase.text), &actual)
			assert.NoError(t, err)
			assert.True(t, testCase.expected.Equal(actual))
		}
	})
}

func TestDecimal_Uninitialized(t *testing.T) {
	t.Parallel()
	t.Run("String", func(t *testing.T) {
		var d decimal.Decimal
		assert.Equal(t, "0", d.String())
	})

	t.Run("Equal", func(t *testing.T) {
		var d decimal.Decimal
		assert.True(t, d.Equal(decimal.MustFromStr("0")))
	})

	t.Run("Cmp", func(t *testing.T) {
		var d decimal.Decimal
		assert.Equal(t, 0, d.Cmp(decimal.MustFromStr("0")))
	})

	t.Run("Add", func(t *testing.T) {
		var d decimal.Decimal
		assert.Equal(t, "1.23", d.Add(decimal.MustFromStr("1.23")).String())
	})

	t.Run("Sub", func(t *testing.T) {
		var d decimal.Decimal
		assert.Equal(t, "-1.23", d.Sub(decimal.MustFromStr("1.23")).String())
	})

	t.Run("Mul", func(t *testing.T) {
		var d decimal.Decimal
		assert.Equal(t, "0", d.Mul(decimal.MustFromStr("1.23")).Reduce().String())
	})

	t.Run("Div", func(t *testing.T) {
		var d decimal.Decimal
		assert.Equal(t, "0", d.Div(decimal.MustFromStr("1.23")).Reduce().String())
	})

	t.Run("Round", func(t *testing.T) {
		var d decimal.Decimal
		assert.Equal(t, "0", d.Round(0).String())
	})

	t.Run("InRangeInt", func(t *testing.T) {
		var d decimal.Decimal
		assert.False(t, d.InRangeInt(1, 10))
	})

	t.Run("Value", func(t *testing.T) {
		var d decimal.Decimal
		val, err := d.Value()
		assert.NoError(t, err)
		assert.Equal(t, "0", val)
	})

	t.Run("FromFloat64", func(t *testing.T) {
		d, err := decimal.FromFloat64(1.23)
		assert.NoError(t, err)
		assert.Equal(t, "1.23", d.String())
	})
}
