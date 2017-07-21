package main

import (
	"math/big"
	//"fmt"
	"math"
)


var minusOne *big.Int
var zero *big.Int
var one *big.Int
var two *big.Int
var three *big.Int
var four *big.Int
var oneMillion *big.Int
var firstFewPrimes []*big.Int


func initPrimes() {
	minusOne = big.NewInt(-1)
	zero = big.NewInt(0)
	one = big.NewInt(1)
	two = big.NewInt(2)
	three = big.NewInt(3)
	four = big.NewInt(4)
	oneMillion = big.NewInt(1000000)
	firstFewPrimes = generateFirstFewPrimes()
}


func even (n *big.Int) bool {
	rest := big.NewInt(0)
	rest.Mod(n, two)
	return rest.Cmp(zero) == 0
}



func isPrime(n *big.Int) bool {

	if n.Sign() == -1 {
		panic("negative numbers disallowed")
	}

	cmpTwo := n.Cmp(two)
	cmpMillion := n.Cmp(oneMillion)

	if cmpTwo == -1 {
		return false
	} else if cmpTwo == 0 {
		return true
	} else {
		if cmpMillion <= 0 {
			ret := isPrimeBruteForce(n)
			if ret == false {
				//fmt.Println(n, " - brute force")
			} else {
				//fmt.Println(n, " +")
			}
			return ret
		} else {
			if isPrimeFirstFew(n) == false {
				//fmt.Println(n, " - first few")
				return false
			} else if isPrimeMillerRabin(n) == false {
				//fmt.Println(n, " - miller rabin")
				return false
			} else if isPrimeBruteForce(n) == false {
				//fmt.Println(n , " - brute force")
				return false
			} else {
				//fmt.Println(n, " +")
				return true
			}
		}
	}

	panic("impossible")
	return false
}


func isPrimeFirstFew(n *big.Int) bool {

	rest := big.NewInt(0)

	for _, mod := range firstFewPrimes {
		rest.Mod(n, mod)
		if rest.Cmp(zero) == 0 {
			return false
		}
	}

	return true
}


func isPrimeBruteForce(n *big.Int) bool {

	cmp := n.Cmp(two)
	if cmp == 0 {
		return true
	} else if cmp == -1 {
		// less than two
		//fmt.Println(n, " - less than two")
		return false
	} else if even(n) == true {
		// even
		//fmt.Println(n, " - even")
		return false
	}

	approxSqrt := big.NewInt(2)
	approxSqrt.Exp(approxSqrt, big.NewInt(int64((n.BitLen()+1)/2)), nil)

	rest := big.NewInt(0)

	for mod := big.NewInt(3); mod.Cmp(approxSqrt) <= 0; mod.Add(mod, two) {
		rest.Mod(n, mod)
		if rest.Cmp(zero) == 0 {
			//fmt.Println(n, " - is disisible by ", mod)
			return false
		}
	}

	return true
}


func isPrimeBruteForceSmallInt(n int64) bool {
	if n == 2 {
		return true
	} else if n < 2 || n % 2 == 0 {
		return false
	}

	max := int64(math.Sqrt(float64(n))) + 1

	for i := int64(3); i <= max; i+=2 {
		if n % i == 0 {
			return false
		}
	}
		
	return true
}


func isPrimeMillerRabin(n *big.Int) bool {
	return big.ProbablyPrime(n, n.BitLen())
}


func generateFirstFewPrimes() []*big.Int {

	var primes [200000]*big.Int

	primes[0] = big.NewInt(2)

	count := 1
	for n := int64(3); n <= 1000000; n += 2 {
		if isPrimeBruteForceSmallInt(n) == true {
			primes[count] = big.NewInt(n)
			count += 1
		}
	}

	bigPrimes := make([]*big.Int, count)

	for i := 0; i < count; i += 1 {
		bigPrimes[i] = primes[i]
	}

	return bigPrimes
}



func generatePrimes(min *big.Int, count int, returnChannel chan<- *big.Int) {

	i := big.NewInt(1)
	i.Set(min)
	for ;; i.Add(i, one) {
		if isPrime(i) == true {
			ret := big.NewInt(0)
			ret.Set(i)
			returnChannel <- ret
		}
	}
}

