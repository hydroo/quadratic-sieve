package main

import (
	"fmt"
	"math"
	"math/big"
	"runtime"
	"time"

	"github.com/hydroo/quadratic-sieve/misc"
)


func factorBase(n *big.Int) []*big.Int {

	/* calculate 'S' upper bound for the primes to collect */
	lnn := float64(n.BitLen()) * math.Log(2)

	if lnn < 1.0 {
		/* if this is not done. if n is 1 everything explodes sqrt(<0) = NaN */
		lnn = 1.0
	}

	lnlnn := math.Log(lnn)
	exp := math.Sqrt(lnn*lnlnn) * 0.5

	if (exp >= 43) {
		/* this is reached when trying to factorize about 2^1500 or larger */
		panic("factorBase(): exponent too large ... reimplement this using big ints")
	}

	S := int64(math.Ceil(math.Pow(math.E, exp))) // magic parameter (wikipedia)

	primes := make([]*big.Int, 1)
	primes[0] = misc.MinusOne

	for p := int64(2); p <= S; p += 1 {
		if misc.IsPrimeBruteForceSmallInt(p) {
			
			P := big.NewInt(int64(p))

			nModP := big.NewInt(0)
			nModP.Mod(n, P)

			if p == 2 || nModP.Cmp(misc.Zero) == 0 {
				/* n is always a square rest (mod 2), and 0 is always a squarerest of 0 */
				primes = append(primes, P)
			} else {
				/* euler criterium. given ggt(a,p)=1: n is square rest mod p, iff n**((p-1)/2) \equiv 1 (mod p) */
				result := big.NewInt(0)
				result.Exp(n, big.NewInt((p-1)/2), P)
				result.Mod(result,P)

				if result.Cmp(misc.One) == 0 {
					primes = append(primes, P)
				}
			}
		}
	}

	return primes
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
	L.Exp(misc.Two, big.NewInt(exp), nil) // magic parameter (wikipedia)

	sqrtN := misc.SquareRootCeil(n)

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
	for ; ci.Cmp(cMax) == -1; ci.Add(ci, misc.One) {

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
						rest.Mul(rest, misc.MinusOne)
					} else {
						exponents[0] = 0
					}
					repeat = false
					continue
				}

				tmpQuotient.DivMod(rest, p, tmpRest)

				if tmpRest.Cmp(misc.Zero) == 0 {
					exponents[i] += 1
					rest.Set(tmpQuotient)
				} else {
					repeat = false
				}

				if tmpQuotient.Cmp(misc.Zero) == 0 {
					repeat = false
				}

			}
		}

		/* if rest is 1 c(i)^2 - n has been successfully broken down into a number that can be represented through
		the factor base -> save c(i)^2 and the exponents for the prime factors */
		if rest.Cmp(misc.One) == 0 {
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
	intervalBig.Add(intervalBig, misc.One)

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

			if x.Cmp(misc.Zero) == 0 || x.Cmp(misc.One) == 0 || y.Cmp(misc.Zero) == 0 || y.Cmp(misc.One) == 0 {
				/* no trivial divisors */
				continue
			}

			xTimesY := big.NewInt(0)
			test := big.NewInt(0)
			testMod := big.NewInt(0)

			xTimesY.Mul(x, y)
			test.DivMod(xTimesY, n, testMod)

			if testMod.Cmp(misc.Zero) != 0 {
				continue
			}

			gcd := big.NewInt(0)

			for test.Cmp(misc.One) == 1 {

				gcd.GCD(nil, nil, x, test)

				if gcd.Cmp(misc.One) == 1 {
					x.Div(x, gcd)
					test.Div(test, gcd)
				}

				gcd.GCD(nil, nil, y, test)

				if gcd.Cmp(misc.One) == 1 {
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
func factorize(n *big.Int, benchmark bool) (*big.Int, *big.Int) {

	t1 := time.Now()

	factorBase := factorBase(n)

	//t2 := time.Now()

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

	if x != nil && y != nil && x.Cmp(y) == 1 {
		x, y = y, x
	}

	if x != nil && y != nil {
		fmt.Print(n, " + ", x, y)
		if benchmark == true {
			fmt.Print(" wall ", nanoSecondsToString(t5.Sub(t1).Nanoseconds()),
			" sieve ", nanoSecondsToString(t4.Sub(t3).Nanoseconds()),
			" combing ", nanoSecondsToString(t5.Sub(t4).Nanoseconds()))
		}
	} else {
		fmt.Print(n, " - - -")
		if benchmark == true {
			fmt.Print(" wall ", nanoSecondsToString(t5.Sub(t1).Nanoseconds()),
			" sieve ", nanoSecondsToString(t4.Sub(t3).Nanoseconds()),
			" combing ", nanoSecondsToString(t5.Sub(t4).Nanoseconds()))
		}
	}
	fmt.Println()

	//fmt.Println("n:", n,  "sieve interval: [", min, "..", max, "] =", max.Int64() - min.Int64(), "factorbase:", factorBase, "factoredsievenums(",len(factoredSieveNums),"):", factoredSieveNums, "c(i)", cis, "exponents:", exponents, "result:", x, "*", y)

	return x, y
}

