// H31.02.10/H31.02.11 by SUZUKI Hisao

package goarith

import (
	"fmt"
	"math/big"
)

func ExampleBigInt_fibonacci() {
	var a Number = Int32(0)
	var b Number = Int32(1)
	for i := 0; i <= 100; i++ {
		if i%10 == 0 {
			fmt.Printf("%3d %T \t%s\n", i, a, a.String())
		}
		a = a.Add(b)
		a, b = b, a
	}
	// Output:
	//   0 goarith.Int32 	0
	//  10 goarith.Int32 	55
	//  20 goarith.Int32 	6765
	//  30 goarith.Int32 	832040
	//  40 goarith.Int32 	102334155
	//  50 goarith.Int64 	12586269025
	//  60 goarith.Int64 	1548008755920
	//  70 goarith.Int64 	190392490709135
	//  80 goarith.Int64 	23416728348467685
	//  90 goarith.Int64 	2880067194370816120
	// 100 *goarith.BigInt 	354224848179261915075
}

func ExampleBigInt_factorial() {
	a := AsNumber(1)
	for i := Int64(2); i <= 40; i++ {
		a = a.Mul(i)
		if i%10 == 0 {
			fmt.Printf("%3d %T \t%s\n", i, a, a.String())
		}
	}
	// Output:
	//  10 goarith.Int32 	3628800
	//  20 goarith.Int64 	2432902008176640000
	//  30 *goarith.BigInt 	265252859812191058636308480000000
	//  40 *goarith.BigInt 	815915283247897734345611269596115894272000000000
}

func ExampleFloat64_String() {
	var a Float64 = 1.234
	fmt.Println(a.String())
	fmt.Println(Float64(5.000).String())
	// Output:
	// 1.234
	// 5.0
}

func ExampleInt64_QuoRem() {
	q, r := Int64(13).QuoRem(Int64(4))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	q, r = Int64(-13).QuoRem(Int64(4))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	q, r = Int64(13).QuoRem(Int64(-4))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	// Output:
	// goarith.Int32 3, goarith.Int32 1
	// goarith.Int32 -3, goarith.Int32 -1
	// goarith.Int32 -3, goarith.Int32 1
}

func ExampleFloat64_QuoRem() {
	q, r := Float64(13).QuoRem(Float64(4))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	q, r = Float64(-13).QuoRem(Float64(4))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	q, r = Float64(13).QuoRem(Float64(-4))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	q, r = Float64(13.4).QuoRem(Float64(1))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	q, r = Float64(-13.4).QuoRem(Float64(1))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	// Output:
	// goarith.Float64 3.0, goarith.Float64 1.0
	// goarith.Float64 -3.0, goarith.Float64 -1.0
	// goarith.Float64 -3.0, goarith.Float64 1.0
	// goarith.Float64 13.0, goarith.Float64 0.40000000000000036
	// goarith.Float64 -13.0, goarith.Float64 -0.40000000000000036
}

func ExampleBigInt_QuoRem() {
	q, r := (*BigInt)(big.NewInt(13)).QuoRem(Int64(4))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	q, r = (*BigInt)(big.NewInt(-13)).QuoRem(Int64(4))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	q, r = (*BigInt)(big.NewInt(13)).QuoRem(Int64(-4))
	fmt.Printf("%T %s, %T %s\n", q, q.String(), r, r.String())
	// Output:
	// goarith.Int32 3, goarith.Int32 1
	// goarith.Int32 -3, goarith.Int32 -1
	// goarith.Int32 -3, goarith.Int32 1
}
