# A package for general numeric arithmetic in Go

This package, `goarith`, implements mixed mode arithmetic
of `int32`, `int64`, `float64` and `*big.Int`.

The package defines four concrete types:

```Go
type Int32 int32
type Int64 int64
type Float64 float64
type BigInt big.Int
```

`Int32`, `Int64`, `Float64` and `*BigInt` implement `Number`:

```Go
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

	// RQuo returns the rounded quotient this/b.
	RQuo(b Number) Float64

	// QuoRem returns the quotient and the remainder of this/b.
	QuoRem(b Number) (quotient Number, remainder Number)
}
```

## Example

The following example computes factorial numbers.

```Go
package main

import (
	"fmt"
	"github.com/nukata/goarith"
)

func main() {
	var a goarith.Number = goarith.Int32(1)
	for i := goarith.Int32(2); i <= 30; i++ {
		a = a.Mul(i)
		if i%10 == 0 {
			fmt.Printf("%3d %T\t%s\n", i, a, a.String())
		}
	}
}
```

It will print the output as follows.

```
 10 goarith.Int32	3628800
 20 goarith.Int64	2432902008176640000
 30 *goarith.BigInt	265252859812191058636308480000000
```
