package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/hydroo/quadratic-sieve/misc"
)

func main() {

	var helpText string
	helpText += "generate [options]                                             \n"
	helpText += "                                                               \n"
	helpText += "  --primes <min>            generates primes starting at 'min' \n"
	helpText += "                            at 'min'                           \n"
	helpText += "  --composites <min>        generates composites larger or     \n"
	helpText += "                            equal than 'min'                   \n"
	helpText += "                            composed of two primes             \n"
	helpText += "  default for 'min' is 1                                       \n"

	args := os.Args[1:]

	primes := false
	composites := false

	min := big.NewInt(1)

	if len(args) == 0 {
		fmt.Print(helpText)
		os.Exit(-1)
	}

	for i := 0; i < len(args); i++ {

		if args[i] == "--primes" || args[i] == "--composites" {

			var boolerr bool

			if args[i] == "--primes" {
				primes = true
			} else if args[i] == "--composites" {
				composites = true
			} else {
				panic("impossible")
			}

			i += 1

			if i < len(args) && args[i][0] != '-' {
				min, boolerr = min.SetString(args[i], 10)

				if boolerr == false {
					panic(boolerr)
				}
			}

		} else if args[i] == "-h" || args[i] == "--help" {

			fmt.Print(helpText)
			os.Exit(-1)

		} else {

			fmt.Println("unknown argument: ", args[i])
			os.Exit(-1)

		}
	}

	if primes == true && composites == true {
		panic("cannot do both at the same time")
	}

	if primes == false && composites == false {
		panic("choose primes or composites")
	}

	if min.Cmp(misc.One) == -1 {
		panic("min < 1")
	}

	if primes == true {

		channel := make(chan *big.Int)

		go generatePrimes(min, channel)

		for n := range channel {
			fmt.Println(n)
		}

	} else {
		channel := make(chan []*big.Int)

		go generateComposites(min, channel)

		for x := range channel {
			fmt.Println(x[0], "=", x[1], "*", x[2])
		}

	}

}


func generateComposites(min *big.Int, returnChannel chan<- []*big.Int) {

	firstPrimeChannel := make(chan *big.Int)
	secondPrimeChannel := make(chan *big.Int)

	leftStart := big.NewInt(2)
	leftStart.Exp(leftStart, big.NewInt(int64(min.BitLen())), nil)
	leftStart.Rsh(leftStart, uint(min.BitLen()*3/4))

	rightStart := big.NewInt(2)
	rightStart.Exp(rightStart, big.NewInt(int64(min.BitLen())), nil)
	rightStart.Rsh(rightStart, uint(min.BitLen()/4))

	go generatePrimes(leftStart, firstPrimeChannel)
	go generatePrimes(rightStart, secondPrimeChannel)

	for {
		ret := make([]*big.Int, 3)

		ret[2] = <-firstPrimeChannel
		ret[1] = <-secondPrimeChannel


		ret[0] = big.NewInt(0)
		ret[0].Mul(ret[2], ret[1])

		returnChannel <- ret
	}


}


func generatePrimes(min *big.Int, returnChannel chan<- *big.Int) {

	i := big.NewInt(1)
	i.Set(min)
	for ;; i.Add(i, misc.One) {
		if misc.IsPrime(i) == true {
			ret := big.NewInt(0)
			ret.Set(i)
			returnChannel <- ret
		}
	}
}

