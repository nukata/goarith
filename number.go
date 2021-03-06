// H31.02.10/R01.10.27 by SUZUKI Hisao

// Package goarith implements general numeric arithmetic.
package goarith

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"strconv"
	"strings"
)

// Number is a general numeric type.
type Number interface {
	// String returns a string representation of the number.
	String() string

	// Int returns the int value for this and a bool indicating whether
	// the int value represents this exactly.
	Int() (i int, exact bool)

	// Add adds this and b (i.e. it return this + b).
	Add(b Number) Number

	// Sub subtracts b from this (i.e. it returns this - b).
	Sub(b Number) Number

	// Cmp compares this and b and returns:
	//
	// -1 if this <  b
	//  0 if this == b
	//  1 if this >  b
	//
	Cmp(b Number) int

	// Mul multiplies this by b (i.e. it returns this * b).
	Mul(b Number) Number

	// RQuo returns the rounded quotient of this and b.
	RQuo(b Number) Float64

	// QuoRem returns the quotient and the remainder of this and b.
	// The quotient will be an Int32, Int64 or BigInt.
	QuoRem(b Number) (quotient Number, remainder Number)
}

// Int32 implements Number.
type Int32 int32

// Int64 implements Number.
type Int64 int64

// Float64 implements Number.
type Float64 float64

// *BigInt implements Number.
type BigInt big.Int

const (
	MaxInt = (1<<bits.UintSize)/2 - 1
	MinInt = (1 << bits.UintSize) / -2
)

// String methods

func (a Int32) String() string {
	return strconv.FormatInt(int64(a), 10)
}

func (a Int64) String() string {
	return strconv.FormatInt(int64(a), 10)
}

func (a Float64) String() string {
	s := strconv.FormatFloat(float64(a), 'g', -1, 64)
	if !strings.ContainsAny(s, ".e") {
		s += ".0"
	}
	return s
}

func (a *BigInt) String() string {
	return (*big.Int)(a).String()
}

// AsNumber converts a numeric value into a Number.
// The numeric value may be int32, int64, int, float32, float64 or *big.Int.
// For Int32, Int64, Float64 and *BigInt, it behaves as an identity function.
// For the other types, it returns nil.
func AsNumber(a interface{}) Number {
	switch x := a.(type) {
	case Int32:
		return x
	case Int64:
		return x
	case Float64:
		return x
	case *BigInt:
		return x
	case int32:
		return Int32(x)
	case int64:
		return Int64(x).reduce()
	case int:
		return Int64(x).reduce()
	case float32:
		return Float64(x)
	case float64:
		return Float64(x)
	case *big.Int:
		return (*BigInt)(x).reduce()
	}
	return nil
}

// Int methods

func (a Int32) Int() (int, bool) {
	return int(a), true
}

func (a Int64) Int() (int, bool) {
	if bits.UintSize >= 64 {
		return int(a), true
	} else if MinInt <= a && a <= MaxInt {
		return int(a), true
	} else if a < 0 {
		return MinInt, false
	} else {
		return MaxInt, false
	}
}

func (a Float64) Int() (int, bool) {
	return int(a), false
}

func (a *BigInt) Int() (int, bool) {
	x := (*big.Int)(a)
	if x.IsInt64() {
		i := x.Int64()
		return Int64(i).Int()
	} else if x.Sign() < 0 {
		return MinInt, false
	} else {
		return MaxInt, false
	}
}

// Utilities

func (a *BigInt) toFloat64() Float64 {
	z := new(big.Rat).SetInt((*big.Int)(a))
	f, _ := z.Float64() // f may be infinity.
	return Float64(f)
}

func (a Int64) reduce() Number {
	if math.MinInt32 <= a && a <= math.MaxInt32 {
		return Int32(a)
	}
	return a
}

func (a *BigInt) reduce() Number {
	if x := (*big.Int)(a); x.IsInt64() {
		i := x.Int64()
		return Int64(i).reduce()
	}
	return a
}

func (a Int64) addInt64(b Int64) Number {
	c := a + b
	if (a >= 0 && b >= 0 && c < 0) || (a < 0 && b < 0 && c >= 0) { // overflow
		z := big.NewInt(int64(a))
		z.Add(z, big.NewInt(int64(b)))
		return (*BigInt)(z)
	}
	return c.reduce()
}

func (a *BigInt) addBigInt(b *big.Int) Number {
	z := new(big.Int)
	z.Add((*big.Int)(a), b)
	return (*BigInt)(z).reduce()
}

func (a Int64) subInt64(b Int64) Number {
	neg := -b
	if neg != b { // b != 0x800...00
		return a.addInt64(neg)
	}
	if a < 0 {
		return (a - b).reduce()
	}
	z := big.NewInt(int64(a))
	z.Sub(z, big.NewInt(int64(b)))
	return (*BigInt)(z)
}

func (a *BigInt) subBigInt(b *big.Int) Number {
	z := new(big.Int)
	z.Sub((*big.Int)(a), b)
	return (*BigInt)(z).reduce()
}

func (a Int64) cmpInt64(b Int64) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

func (a Float64) cmpFloat64(b Float64) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

func (a Int64) mulInt64(b Int64) Number {
	z := big.NewInt(int64(a))
	z.Mul(z, big.NewInt(int64(b)))
	return (*BigInt)(z).reduce()
}

func (a *BigInt) mulBigInt(b *big.Int) Number {
	z := new(big.Int)
	z.Mul((*big.Int)(a), b)
	return (*BigInt)(z).reduce()
}

func (a Int64) quoRemInt64(b Int64) (Number, Number) {
	return (a / b).reduce(), (a % b).reduce()
}

func (a Float64) quoRemFloat64(b Float64) (Number, Float64) {
	q := math.Trunc(float64(a) / float64(b))
	r := math.Mod(float64(a), float64(b))
	if !math.IsInf(q, 0) && !math.IsNaN(q) {
		s := fmt.Sprintf("%.0f", q)
		z := new(big.Int)
		if _, ok := z.SetString(s, 10); ok {
			return (*BigInt)(z).reduce(), Float64(r)
		}
	}
	return Float64(q), Float64(r)
}

func (a *BigInt) quoRemBigInt(b *big.Int) (Number, Number) {
	q := new(big.Int)
	r := new(big.Int)
	q.QuoRem((*big.Int)(a), b, r)
	return (*BigInt)(q).reduce(), (*BigInt)(r).reduce()
}

// Add methods

func (a Int32) Add(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return (Int64(a) + Int64(y)).reduce()
	case Int64:
		return Int64(a).addInt64(y)
	case Float64:
		return Float64(a) + y
	case *BigInt:
		x := big.NewInt(int64(a))
		x.Add(x, (*big.Int)(y))
		return (*BigInt)(x).reduce()
	}
	panic(fmt.Sprintf("%s.Add(%s)", a.String(), b.String()))
}

func (a Int64) Add(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return a.addInt64(Int64(y))
	case Int64:
		return a.addInt64(y)
	case Float64:
		return Float64(a) + y
	case *BigInt:
		x := big.NewInt(int64(a))
		x.Add(x, (*big.Int)(y))
		return (*BigInt)(x).reduce()
	}
	panic(fmt.Sprintf("%s.Add(%s)", a.String(), b.String()))
}

func (a Float64) Add(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return a + Float64(y)
	case Int64:
		return a + Float64(y)
	case Float64:
		return a + y
	case *BigInt:
		return a + y.toFloat64()
	}
	panic(fmt.Sprintf("%s.Add(%s)", a.String(), b.String()))
}

func (a *BigInt) Add(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return a.addBigInt(big.NewInt(int64(y)))
	case Int64:
		return a.addBigInt(big.NewInt(int64(y)))
	case Float64:
		return a.toFloat64() + y
	case *BigInt:
		return a.addBigInt((*big.Int)(y))
	}
	panic(fmt.Sprintf("%s.Add(%s)", a.String(), b.String()))
}

// Sub methods

func (a Int32) Sub(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return (Int64(a) - Int64(y)).reduce()
	case Int64:
		return Int64(a).subInt64(y)
	case Float64:
		return Float64(a) - y
	case *BigInt:
		x := big.NewInt(int64(a))
		x.Sub(x, (*big.Int)(y))
		return (*BigInt)(x).reduce()
	}
	panic(fmt.Sprintf("%s.Sub(%s)", a.String(), b.String()))
}

func (a Int64) Sub(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return a.subInt64(Int64(y))
	case Int64:
		return a.subInt64(y)
	case Float64:
		return Float64(a) - y
	case *BigInt:
		x := big.NewInt(int64(a))
		x.Sub(x, (*big.Int)(y))
		return (*BigInt)(x).reduce()
	}
	panic(fmt.Sprintf("%s.Sub(%s)", a.String(), b.String()))
}

func (a Float64) Sub(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return a - Float64(y)
	case Int64:
		return a - Float64(y)
	case Float64:
		return a - y
	case *BigInt:
		return a - y.toFloat64()
	}
	panic(fmt.Sprintf("%s.Sub(%s)", a.String(), b.String()))
}

func (a *BigInt) Sub(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return a.subBigInt(big.NewInt(int64(y)))
	case Int64:
		return a.subBigInt(big.NewInt(int64(y)))
	case Float64:
		return a.toFloat64() - y
	case *BigInt:
		return a.subBigInt((*big.Int)(y))
	}
	panic(fmt.Sprintf("%s.Sub(%s)", a.String(), b.String()))
}

// Cmp methods

func (a Int32) Cmp(b Number) int {
	switch y := b.(type) {
	case Int32:
		if a < y {
			return -1
		} else if a > y {
			return 1
		} else {
			return 0
		}
	case Int64:
		return Int64(a).cmpInt64(y)
	case Float64:
		return Float64(a).cmpFloat64(y)
	case *BigInt:
		x := big.NewInt(int64(a))
		return x.Cmp((*big.Int)(y))
	}
	panic(fmt.Sprintf("%s.Cmp(%s)", a.String(), b.String()))
}

func (a Int64) Cmp(b Number) int {
	switch y := b.(type) {
	case Int32:
		return a.cmpInt64(Int64(y))
	case Int64:
		return a.cmpInt64(y)
	case Float64:
		return Float64(a).cmpFloat64(y)
	case *BigInt:
		x := big.NewInt(int64(a))
		return x.Cmp((*big.Int)(y))
	}
	panic(fmt.Sprintf("%s.Cmp(%s)", a.String(), b.String()))
}

func (a Float64) Cmp(b Number) int {
	switch y := b.(type) {
	case Int32:
		return a.cmpFloat64(Float64(y))
	case Int64:
		return a.cmpFloat64(Float64(y))
	case Float64:
		return a.cmpFloat64(y)
	case *BigInt:
		return a.cmpFloat64(y.toFloat64())
	}
	panic(fmt.Sprintf("%s.Cmp(%s)", a.String(), b.String()))
}

func (a *BigInt) Cmp(b Number) int {
	switch y := b.(type) {
	case Int32:
		return (*big.Int)(a).Cmp(big.NewInt(int64(y)))
	case Int64:
		return (*big.Int)(a).Cmp(big.NewInt(int64(y)))
	case Float64:
		return a.toFloat64().cmpFloat64(y)
	case *BigInt:
		return (*big.Int)(a).Cmp((*big.Int)(y))
	}
	panic(fmt.Sprintf("%s.Cmp(%s)", a.String(), b.String()))
}

// Mul methods

func (a Int32) Mul(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return (Int64(a) * Int64(y)).reduce()
	case Int64:
		return Int64(a).mulInt64(y)
	case Float64:
		return Float64(a) * y
	case *BigInt:
		x := big.NewInt(int64(a))
		x.Mul(x, (*big.Int)(y))
		return (*BigInt)(x).reduce()
	}
	panic(fmt.Sprintf("%s.Mul(%s)", a.String(), b.String()))
}

func (a Int64) Mul(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return a.mulInt64(Int64(y))
	case Int64:
		return a.mulInt64(y)
	case Float64:
		return Float64(a) * y
	case *BigInt:
		x := big.NewInt(int64(a))
		x.Mul(x, (*big.Int)(y))
		return (*BigInt)(x).reduce()
	}
	panic(fmt.Sprintf("%s.Mul(%s)", a.String(), b.String()))
}

func (a Float64) Mul(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return a * Float64(y)
	case Int64:
		return a * Float64(y)
	case Float64:
		return a * y
	case *BigInt:
		return a * y.toFloat64()
	}
	panic(fmt.Sprintf("%s.Mul(%s)", a.String(), b.String()))
}

func (a *BigInt) Mul(b Number) Number {
	switch y := b.(type) {
	case Int32:
		return a.mulBigInt(big.NewInt(int64(y)))
	case Int64:
		return a.mulBigInt(big.NewInt(int64(y)))
	case Float64:
		return a.toFloat64() + y
	case *BigInt:
		return a.mulBigInt((*big.Int)(y))
	}
	panic(fmt.Sprintf("%s.Mul(%s)", a.String(), b.String()))
}

// RQuo methods

func (a Int32) RQuo(b Number) Float64 {
	return Float64(a).RQuo(b)
}

func (a Int64) RQuo(b Number) Float64 {
	return Float64(a).RQuo(b)
}

func (a Float64) RQuo(b Number) Float64 {
	switch y := b.(type) {
	case Int32:
		return a / Float64(y)
	case Int64:
		return a / Float64(y)
	case Float64:
		return a / y
	case *BigInt:
		return a / y.toFloat64()
	}
	panic(fmt.Sprintf("%s.RQuo(%s)", a.String(), b.String()))
}

func (a *BigInt) RQuo(b Number) Float64 {
	return a.toFloat64().RQuo(b)
}

// QuoRem methods

func (a Int32) QuoRem(b Number) (Number, Number) {
	switch y := b.(type) {
	case Int32:
		return a / y, a % y
	case Int64:
		return Int64(a).quoRemInt64(y)
	case Float64:
		return Float64(a).quoRemFloat64(y)
	case *BigInt:
		x := big.NewInt(int64(a))
		return (*BigInt)(x).quoRemBigInt((*big.Int)(y))
	}
	panic(fmt.Sprintf("%s.RQuoRem(%s)", a.String(), b.String()))
}

func (a Int64) QuoRem(b Number) (Number, Number) {
	switch y := b.(type) {
	case Int32:
		return a.quoRemInt64(Int64(y))
	case Int64:
		return a.quoRemInt64(y)
	case Float64:
		return Float64(a).quoRemFloat64(y)
	case *BigInt:
		x := big.NewInt(int64(a))
		return (*BigInt)(x).quoRemBigInt((*big.Int)(y))
	}
	panic(fmt.Sprintf("%s.RQuoRem(%s)", a.String(), b.String()))
}

func (a Float64) QuoRem(b Number) (Number, Number) {
	switch y := b.(type) {
	case Int32:
		return a.quoRemFloat64(Float64(y))
	case Int64:
		return a.quoRemFloat64(Float64(y))
	case Float64:
		return a.quoRemFloat64(y)
	case *BigInt:
		return a.quoRemFloat64(y.toFloat64())
	}
	panic(fmt.Sprintf("%s.RQuoRem(%s)", a.String(), b.String()))
}

func (a *BigInt) QuoRem(b Number) (Number, Number) {
	switch y := b.(type) {
	case Int32:
		return a.quoRemBigInt(big.NewInt(int64(y)))
	case Int64:
		return a.quoRemBigInt(big.NewInt(int64(y)))
	case Float64:
		return a.toFloat64().quoRemFloat64(y)
	case *BigInt:
		return a.quoRemBigInt((*big.Int)(y))
	}
	panic(fmt.Sprintf("%s.RQuoRem(%s)", a.String(), b.String()))
}
