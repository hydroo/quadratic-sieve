package main


import (
	//"fmt"
	"math/big"
	"testing"
)


func TestFactorize(t *testing.T) {

	type Number struct {
		n, x, y int64
	}


	nums := []Number{
			{1649, 17, 97},
			//{7429, 19, 391},
			//{7429, 17, 437},
			{7429, 23, 323},
			{40198364677, 599, 67109123},
			{18923626564873, 2203, 8589934891},
			//{362684905587521, 66847, 5425597343},
			{362684905587521, 1133, 320110243237},
			{2626849055875131, 4549, 577456376319},
			//{2626849055875147, 25025783, 104965709},
			{2626849055875147, 128477, 20446064711},
			}

	xShould := big.NewInt(0)
	yShould := big.NewInt(0)


	for _, num := range nums {

		x, y := factorize(big.NewInt(num.n), false)

		if x.Cmp(y) > 0 {
			x, y = y, x
		}

		xShould.SetInt64(num.x)
		yShould.SetInt64(num.y)

		if xShould.Cmp(x) != 0 || yShould.Cmp(y) != 0 {
			t.Error(num.n, "=", xShould, "*", yShould, "but " ,num.n ,"=", x, "*", y)
		} else {
			//fmt.Println(num.n, "=",x,"*",y)
		}
	}
}

