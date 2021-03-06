package main

import (
	"big"
	"fmt"
	"testing"
)

func TestFactorize(t *testing.T) {

	initPrimes()

	type Number struct {
		n, x, y int64
	}


	nums := []Number{
			{1649, 17, 97},
			/*{7429, 19, 391}, */
			{7429, 17, 437},
			{40198364677, 599, 67109123},
			{18923626564873, 2203, 8589934891},
			{362684905587521, 66847, 5425597343},
			{2626849055875131, 4549, 577456376319},
			{2626849055875147, 25025783, 104965709}}

	xShould := big.NewInt(0)
	yShould := big.NewInt(0)

	for _, num := range nums {

		x, y := factorize(big.NewInt(num.n))

		xShould.SetInt64(num.x)
		yShould.SetInt64(num.y)

		if xShould.Cmp(x) != 0 || yShould.Cmp(y) != 0 {
			t.Errorf(fmt.Sprint(num.n, "=", xShould, "*", yShould, "but " ,num.n ,"=", x, "*", y))
		} else {
			fmt.Println(num.n, "=",x,"*",y)
		}
	}
}

