package main

import (
	"fmt"
	"math"
	"math/big"
	"runtime"
	"time"
)

/* helper */
func nanoSecondsToString(n int64) string {
	return fmt.Sprintf("%f", float64(n)/1000000000.0)
}

func SquareRootCeil(n *big.Int) *big.Int {

	if n.Cmp(one) == -1 {
		panic("cannot get square root of a number smaller than one")
	}

	upperLimit := big.NewInt(2)
	lowerLimit := big.NewInt(2)

	upperLimitExp := big.NewInt(int64(math.Ceil(float64(n.BitLen()) / 2)))
	lowerLimitExp := big.NewInt(int64(math.Floor(float64(n.BitLen()-1) / 2)))

	upperLimit.Exp(upperLimit, upperLimitExp, nil)
	lowerLimit.Exp(lowerLimit, lowerLimitExp, nil)

	middle := big.NewInt(0)

	middleSquared := big.NewInt(0)

	/* binary search */
	for upperLimit.Cmp(lowerLimit) != 0 {

		if upperLimit.Cmp(lowerLimit) == -1 {
			panic("upperlimit < lowerlimit shouldnt happen")
		}

		middle.Add(upperLimit, lowerLimit)
		middle.Div(middle, two)

		middleSquared.Exp(middle, two, nil)

		if middleSquared.Cmp(n) == -1 {

			if lowerLimit.Cmp(middle) == 0 {
				return upperLimit
			}

			lowerLimit.Set(middle)
		} else {
			upperLimit.Set(middle)
		}
	}

	return upperLimit
}

func factorBase(n *big.Int) []*big.Int {

	/* calculate 'S' upper bound for the primes to collect */
	lnn := float64(n.BitLen()) * math.Log(2)

	if lnn < 1.0 {
		/* if this is not done. if n is 1 everything explodes sqrt(<0) = NaN */
		lnn = 1.0
	}

	lnlnn := math.Log(lnn)
	exp := math.Sqrt(lnn*lnlnn) * 0.5
	S := int(math.Ceil(math.Pow(math.E, exp))) // magic parameter (wikipedia)

	primes := make([]*big.Int, S)

	nModP := big.NewInt(0)

	count := 0
	for p := 2; p <= S; p += 1 {
		if isPrimeBruteForceSmallInt(int64(p)) {

			P := big.NewInt(int64(p))

			nModP.Mod(n, P)

			if nModP.BitLen() > 63 {
				panic("oh noez, number to large. rewrite this code")
			}

			/* iterate through the whole ring of Z_p and test wether some i^2 equals n -> it is square mod p */
			isSquare := false
			iSquaredModP := big.NewInt(0)
			i := big.NewInt(0)
			for ; i.Cmp(P) == -1; i.Add(i, one) {
				iSquaredModP.Exp(i, two, P)
				if iSquaredModP.Cmp(nModP) == 0 {
					isSquare = true
					break
				}
			}

			if isSquare == true {
				primes[count] = P
				count += 1
			}
		}
	}

	/* copy results into a new array to save space */
	ret := make([]*big.Int, count+1)
	ret[0] = big.NewInt(-1)
	for i := 0; i < count; i += 1 {
		ret[i+1] = primes[i]
	}

	return ret
}

func sieveInterval(n *big.Int) (*big.Int, *big.Int) {

	lnn := float64(n.BitLen()) * math.Log(2)

	if lnn < 1.0 {
		/* if this is not done. if n is 1 everything explodes sqrt(<0) = NaN */
		lnn = 1.0
	}

	lnlnn := math.Log(lnn)
	exp := int64(math.Ceil(math.Sqrt(lnn*lnlnn) * math.Log2(math.E)))

	L := big.NewInt(0)
	L.Exp(two, big.NewInt(exp), nil) // magic parameter (wikipedia)

	sqrtN := SquareRootCeil(n)

	sieveMin := big.NewInt(0)
	sieveMin.Sub(sqrtN, L)
	sieveMax := big.NewInt(0)
	sieveMax.Add(sqrtN, L)

	return sieveMin, sieveMax
}

func sieveWork(n *big.Int, factorBase []*big.Int, cMin, cMax *big.Int, retCis, retDs []*big.Int, retExponents [][]int, done chan<- bool) {

	retIterator := 0

	ci := big.NewInt(0)
	ci.Set(cMin)
	d := big.NewInt(0)
	rest := big.NewInt(0)
	tmpRest := big.NewInt(0)
	tmpQuotient := big.NewInt(0)
	exponents := make([]int, len(factorBase))

	/* foreach c(i) */
	for ; ci.Cmp(cMax) == -1; ci.Add(ci, one) {

		/* calculate c(i)^1 - n */
		d.Mul(ci, ci)
		d.Sub(d, n)

		/* copy the result because we need to modify it and save c(i)^2 - n as well */
		rest.Set(d)

		for i, p := range factorBase {

			exponents[i] = 0

			repeat := true

			/* repeat as long as rest % p == 0 -> add 1 to the exponent for each division by p */
			for repeat == true {

				if i == 0 {
					/* needs special handling: p = -1 */
					if rest.Sign() == -1 {
						exponents[0] = 1
						rest.Mul(rest, minusOne)
					} else {
						exponents[0] = 0
					}
					repeat = false
					continue
				}

				tmpQuotient.DivMod(rest, p, tmpRest)

				if tmpRest.Cmp(zero) == 0 {
					exponents[i] += 1
					rest.Set(tmpQuotient)
				} else {
					repeat = false
				}

				if tmpQuotient.Cmp(zero) == 0 {
					repeat = false
				}

			}
		}

		/* if rest is 1 c(i)^2 - n has been successfully broken down into a number that can be represented through
		the factor base -> save c(i)^2 and the exponents for the prime factors */
		if rest.Cmp(one) == 0 {
			retDs[retIterator] = big.NewInt(0)
			retDs[retIterator].Set(d)
			retCis[retIterator] = big.NewInt(0)
			retCis[retIterator].Set(ci)
			retExponents[retIterator] = make([]int, len(factorBase))
			copy(retExponents[retIterator], exponents)
			retIterator += 1
		}
	}

	done <- true
}

func sieve(n *big.Int, factorBase []*big.Int, cMin, cMax *big.Int) ([]*big.Int, []*big.Int, [][]int) {

	intervalBig := big.NewInt(0)
	intervalBig.Sub(cMax, cMin)
	intervalBig.Add(intervalBig, one)

	if intervalBig.BitLen() > 31 {
		panic("fufufu sieve interval too large. code newly.")
	}

	interval := int(intervalBig.Int64())

	retCiTmp := make([]*big.Int, interval)
	retDsTmp := make([]*big.Int, interval)
	retExponentsTmp := make([][]int, interval)

	threads := runtime.GOMAXPROCS(-1)
	tasksRest := interval % threads
	tasksStep := interval / threads

	tasksStepBig := big.NewInt(int64(tasksStep))
	tasksStepBigPlusOne := big.NewInt(int64(tasksStep + 1))

	ci := big.NewInt(0)
	ci.Set(cMin)

	doneChannel := make(chan bool)

	/* foreach c(i) */
	for i := 0; ci.Cmp(cMax) <= 0; {

		cMinIt := big.NewInt(0)
		cMaxIt := big.NewInt(0)
		cMinIt.Set(ci)

		j := 0

		if tasksRest > 0 {
			cMaxIt.Add(cMinIt, tasksStepBigPlusOne)
			j = i + tasksStep + 1
		} else {
			cMaxIt.Add(cMinIt, tasksStepBig)
			j = i + tasksStep
		}

		go sieveWork(n, factorBase, cMinIt, cMaxIt, retCiTmp[i:j], retDsTmp[i:j], retExponentsTmp[i:j], doneChannel)

		if tasksRest > 0 {
			ci.Add(ci, tasksStepBigPlusOne)
			tasksRest -= 1
			i += tasksStep + 1
		} else {
			ci.Add(ci, tasksStepBig)
			i += tasksStep
		}
	}

	max := threads
	if tasksStep == 0 {
		max = tasksRest
	}

	for i := 0; i < max; i += 1 {
		<-doneChannel
	}

	retIterator := 0
	for _, ci := range retCiTmp {
		if ci != nil {
			retIterator += 1
		}
	}

	/* copy the result into new, shorter arrays */
	retDs := make([]*big.Int, retIterator)
	retCi := make([]*big.Int, retIterator)
	retExponents := make([][]int, retIterator)

	i := 0
	j := 0
	for ; j < interval; j += 1 {

		if retCiTmp[j] != nil {
			retDs[i] = retDsTmp[j] // c(i)^2 - n
			retCi[i] = retCiTmp[j] // c(i)
			retExponents[i] = retExponentsTmp[j]
			i += 1
		}
	}

	return retDs, retCi, retExponents
}

func combineRecursively(start int, currentExponents []int, cMul *big.Int, cSquaredList, cis []*big.Int, exponents [][]int, factorBase []*big.Int, n *big.Int) (*big.Int, *big.Int) {

	for i := start; i < len(cSquaredList); i += 1 {

		newCurrentExponents := make([]int, len(currentExponents))

		for j, k := range exponents[i] {
			newCurrentExponents[j] = k + currentExponents[j]
		}

		newCMul := big.NewInt(1)
		newCMul.Mul(cMul, cis[i])
		newCMul.Mod(newCMul, n)

		even := true
		for _, k := range newCurrentExponents {
			if k%2 == 1 {
				even = false
				break
			}
		}

		if even == true {

			b := big.NewInt(1)
			for i, k := range newCurrentExponents {

				if i == 0 {
					continue /* minus one makes screw ups */
				}

				for j := 0; j < k/2; j += 1 {
					b.Mul(b, factorBase[i])
				}
			}

			a := newCMul

			x := big.NewInt(0)
			x.Add(a, b)
			x.Mod(x, n)

			y := big.NewInt(0)
			y.Sub(a, b)
			y.Mod(y, n)

			if x.Cmp(zero) == 0 || x.Cmp(one) == 0 || y.Cmp(zero) == 0 || y.Cmp(one) == 0 {
				/* no trivial divisors */
				continue
			}

			xTimesY := big.NewInt(0)
			test := big.NewInt(0)
			testMod := big.NewInt(0)

			xTimesY.Mul(x, y)
			test.DivMod(xTimesY, n, testMod)

			if testMod.Cmp(zero) != 0 {
				continue
			}

			gcd := big.NewInt(0)

			for test.Cmp(one) == 1 {

				big.GcdInt(gcd, nil, nil, x, test)

				if gcd.Cmp(one) == 1 {
					x.Div(x, gcd)
					test.Div(test, gcd)
				}

				big.GcdInt(gcd, nil, nil, y, test)

				if gcd.Cmp(one) == 1 {
					y.Div(y, gcd)
					test.Div(test, gcd)
				}

			}

			xTimesY.Mul(x, y)

			if xTimesY.Cmp(n) == 0 {
				//fmt.Println(n, "=", x, "*", y, "(", newUsedIndizes, a, b, ")")
				return x, y
			}
		}

		x, y := combineRecursively(i+1, newCurrentExponents, newCMul, cSquaredList, cis, exponents, factorBase, n)
		if x != nil && y != nil {
			return x, y
		}

	}

	return nil, nil
}

/* usually this step should be performed by a solver for a system of linear equations in Z_2
but this is too much work for now, so it just tests all (all subsets of the powerset of all exponent vectors)
the linear combinations of exponent vectors for evenness

returns nil, nil if nothing is found */
func combine(cSquaredList, cis []*big.Int, exponents [][]int, factorBase []*big.Int, n *big.Int) (*big.Int, *big.Int) {

	if len(exponents) == 0 {
		return nil, nil
	}

	currentExponents := make([]int, len(exponents[0]))
	for i, _ := range currentExponents {
		currentExponents[i] = 0
	}

	cMul := big.NewInt(1)

	return combineRecursively(0, currentExponents, cMul, cSquaredList, cis, exponents, factorBase, n)
}

/* returns nil, nil upon failure */
func factorize(n *big.Int) (*big.Int, *big.Int) {

	t1 := time.Now()

	factorBase := factorBase(n)

	t2 := time.Now()

	min, max := sieveInterval(n)

	t3 := time.Now()

	factoredSieveNums, cis, exponents := sieve(n, factorBase, min, max)

	t4 := time.Now()

	const combinationsLog2 = 20

	if len(cis) > combinationsLog2 {
		/* avoid building too large powersets */
		factoredSieveNums = factoredSieveNums[:combinationsLog2]
		cis = cis[:combinationsLog2]
		exponents = exponents[:combinationsLog2]
	}

	/* usually you want to solve this through an LGS ... TODO? */
	x, y := combine(factoredSieveNums, cis, exponents, factorBase, n)

	t5 := time.Now()
	fmt.Sprint(t1, t2, t3, t4, t5) // Q_UNUSED(...), (void) ...

	if x != nil && y != nil && x.Cmp(y) == 1 {
		x, y = y, x
	}

	if x != nil && y != nil {
		fmt.Println(n, "+", x, y, "wall", nanoSecondsToString(t5.Sub(t1).Nanoseconds()), "sieve", nanoSecondsToString(t4.Sub(t3).Nanoseconds()), "combing", nanoSecondsToString(t5.Sub(t4).Nanoseconds()))
	} else {
		fmt.Println(n, "- - -", "wall", nanoSecondsToString(t5.Sub(t1).Nanoseconds()), "sieve", nanoSecondsToString(t4.Sub(t3).Nanoseconds()), "combing", nanoSecondsToString(t5.Sub(t4).Nanoseconds()))
	}

	//fmt.Println("n:", n,  "sieve interval: [", min, "..", max, "] =", max.Int64() - min.Int64(), "factorbase:", factorBase, "factoredsievenums(",len(factoredSieveNums),"):", factoredSieveNums, "c(i)", cis, "exponents:", exponents, "result:", x, "*", y)

	return x, y

}
