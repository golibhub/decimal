package decimal

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/apd/v3"
)

// DivisionPrecision is the default precision for division used for all the `Decimal.Div`
// operations. It can be changed by setting `decimal.DivisionPrecision` to a different value.
var DivisionPrecision uint32 = 32

// Decimal is an arbitrary-precision decimal.
type Decimal struct {
	v *apd.Decimal
}

func (d *Decimal) ensureInitialized() {
	if d.v == nil {
		d.v = apd.New(0, 0)
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(b []byte) error {
	dec := json.NewDecoder(bytes.NewReader(b))

	var s json.Number
	if err := dec.Decode(&s); err != nil {
		return err
	}

	v, err := FromStr(s.String())
	if err != nil {
		return err
	}

	*d = v

	return nil
}

// FromInt creates a new Decimal from an int64.
func FromInt(i int64) Decimal {
	return Decimal{v: apd.New(i, 0)}
}

// FromStr creates a new Decimal from a string.
func FromStr(str string) (Decimal, error) {
	d, c, err := apd.NewFromString(str)
	if err != nil {
		return Decimal{}, fmt.Errorf("failed to parse %s as decimal; Err: %w; apd.Condition: %v", str, err, c)
	}

	return Decimal{v: d}, nil
}

// MustFromStr creates a new Decimal from a string.
// It panics if the string is not a valid decimal.
func MustFromStr(str string) Decimal {
	d, c, err := apd.NewFromString(str)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %s as decimal; Err: %v; apd.Condition: %v", str, err, c))
	}

	return Decimal{v: d}
}

// FromFloat64 creates a new Decimal from a float64.
func FromFloat64(f float64) (Decimal, error) {
	d := apd.New(0, 0)
	d, err := d.SetFloat64(f)
	if err != nil {
		return Decimal{}, fmt.Errorf("failed to convert %f to decimal; Err: %w", f, err)
	}

	return Decimal{v: d}, nil
}

// Add adds two decimals and returns the result.
func (d Decimal) Add(d2 Decimal) Decimal {
	d.ensureInitialized()
	d2.ensureInitialized()

	var res apd.Decimal
	c, err := apd.BaseContext.Add(&res, d.v, d2.v)
	if err != nil {
		panic(fmt.Sprintf("fialed to [%s + %s]; Err: %v; apd.Condition: %v", d.v.String(), d2.v.String(), err, c))
	}

	return Decimal{v: &res}
}

// Sub subtracts two decimals and returns the result.
func (d Decimal) Sub(d2 Decimal) Decimal {
	d.ensureInitialized()
	d2.ensureInitialized()

	var res apd.Decimal
	c, err := apd.BaseContext.Sub(&res, d.v, d2.v)
	if err != nil {
		panic(fmt.Sprintf("fialed to [%s - %s]; Err: %v; apd.Condition: %v", d.v.String(), d2.v.String(), err, c))
	}

	return Decimal{v: &res}
}

// Mul multiplies two decimals and returns the result.
func (d Decimal) Mul(d2 Decimal) Decimal {
	d.ensureInitialized()
	d2.ensureInitialized()

	var res apd.Decimal
	c, err := apd.BaseContext.Mul(&res, d.v, d2.v)
	if err != nil {
		panic(fmt.Sprintf("fialed to [%s * %s]; Err: %v; apd.Condition: %v", d.v.String(), d2.v.String(), err, c))
	}

	return Decimal{v: &res}
}

// MulInt multiplies a decimal by an int64 and returns the result.
func (d Decimal) MulInt(i int64) Decimal {
	d2 := FromInt(i)
	return d.Mul(d2)
}

// Div divides the current decimal by a given one and returns the result.
func (d Decimal) Div(d2 Decimal) Decimal {
	return d.DivPrecision(d2, DivisionPrecision)
}

// DivPrecision divides the current decimal by a given one and returns the result with a given precision.
func (d Decimal) DivPrecision(d2 Decimal, precision uint32) Decimal {
	d.ensureInitialized()
	d2.ensureInitialized()

	var res apd.Decimal

	c, err := apd.BaseContext.WithPrecision(precision).Quo(&res, d.v, d2.v)
	if err != nil {
		panic(fmt.Sprintf("fialed to [%s / %s]; Err: %v; apd.Condition: %v", d.v.String(), d2.v.String(), err, c))
	}
	res.Reduce(&res)

	return Decimal{v: &res}
}

// Round rounds the decimal to n digits after 0.
func (d Decimal) Round(n uint16) Decimal {
	d.ensureInitialized()

	var res apd.Decimal
	c, err := apd.BaseContext.Quantize(&res, d.v, -int32(n))
	if err != nil {
		panic(fmt.Sprintf("fialed to round %s to %d digits after 0; Err: %v; apd.Condition: %v", d.v.String(), n, err, c))
	}

	return Decimal{v: &res}
}

// Reduce removes all the trailing zeroes from the decimal.
func (d Decimal) Reduce() Decimal {
	var x apd.Decimal
	d.v.Reduce(&x)

	return Decimal{v: &x}
}

// RoundOrNil returns nil if the Decimal is nil or rounds it to n digits after 0.
func (d *Decimal) RoundOrNil(n uint16) *Decimal {
	if d == nil {
		return nil
	}
	valueToRound := *d
	valueRound := valueToRound.Round(n)

	return &valueRound
}

// Cmp compares two decimals and returns:
// -1 if d < d2
// 0 if d == d2
// 1 if d > d2
func (d Decimal) Cmp(d2 Decimal) int {
	d.ensureInitialized()
	d2.ensureInitialized()

	return d.v.Cmp(d2.v)
}

// Less returns true if d < d2.
func (d Decimal) Less(d2 Decimal) bool {
	return d.Cmp(d2) == -1
}

// LessOrEqual returns true if d <= d2.
func (d Decimal) LessOrEqual(d2 Decimal) bool {
	return d.Cmp(d2) <= 0
}

// Greater returns true if d > d2.
func (d Decimal) Greater(d2 Decimal) bool {
	return d.Cmp(d2) == 1
}

// GreaterOrEqual returns true if d >= d2.
func (d Decimal) GreaterOrEqual(d2 Decimal) bool {
	return d.Cmp(d2) >= 0
}

// Equal returns true if d == d2.
func (d Decimal) Equal(d2 Decimal) bool {
	return d.Cmp(d2) == 0
}

// InRangeInt returns true if d is in the range of [min, max] integers.
func (d Decimal) InRangeInt(min, max int64) bool {
	d.ensureInitialized()

	minD := FromInt(min)
	maxD := FromInt(max)

	return d.Cmp(minD) >= 0 && d.Cmp(maxD) <= 0
}

// InRange returns true if d is in the range of [min, max] decimals.
func (d Decimal) InRange(min, max Decimal) bool {
	d.ensureInitialized()

	return d.Cmp(min) >= 0 && d.Cmp(max) <= 0
}

// IsZero returns true if d == 0.
func (d Decimal) IsZero() bool {
	d.ensureInitialized()

	return d.v.IsZero()
}

// IsNegative returns true if d < 0.
func (d *Decimal) IsNegative() bool {
	d.ensureInitialized()

	return d.v.Negative
}

// String returns the string representation of the decimal.
func (d Decimal) String() string {
	d.ensureInitialized()

	return d.v.String()
}

// Float64 returns the float64 representation of the decimal.
func (d Decimal) Float64() (float64, error) {
	d.ensureInitialized()

	return d.v.Float64()
}

// MustFloat64 returns the float64 representation of the decimal.
// It panics if the decimal can't be converted to float64.
func (d Decimal) MustFloat64() float64 {
	d.ensureInitialized()

	f, err := d.Float64()
	if err != nil {
		panic(fmt.Sprintf("failed to convert %s to float64; Err: %v", d.String(), err))
	}

	return f
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Decimal) UnmarshalText(text []byte) error {
	d.ensureInitialized()

	return d.v.UnmarshalText(text)
}

// MarshalText implements encoding.TextMarshaler.
func (d Decimal) MarshalText() (text []byte, err error) {
	d.ensureInitialized()

	return d.v.MarshalText()
}

// Scan implements sql.Scanner.
func (d *Decimal) Scan(src interface{}) error {
	d.ensureInitialized()

	return d.v.Scan(src)
}

// Value implements driver.Valuer.
func (d Decimal) Value() (driver.Value, error) {
	d.ensureInitialized()

	return d.v.Value()
}
