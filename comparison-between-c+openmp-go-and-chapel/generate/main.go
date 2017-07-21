package main

import (
	"big"
	"fmt"
	"strconv"
	"os"
)

func main() {

	initPrimes()

	args := os.Args[1:]

	primes := false
	composites := false

	min := big.NewInt(1)
	count := 100

	for i := 0; i < len(args); i++ {

		if args[i] == "--primes" || args[i] == "--composites" {

			var err os.Error
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

				i += 1
			}

			if i < len(args) && args[i][0] != '-' {
				count, err = strconv.Atoi(args[i])

				if err != nil {
					panic(err)
				}

				i += 1
			}

		} else if args[i] == "-h" || args[i] == "--help" {

			fmt.Println("generate [options]")
			fmt.Println("")
			fmt.Println("  --primes <min> <count>       generates 'count' primes starting")
			fmt.Println("                               at 'min'                         ")
			fmt.Println("  --composites <min> <count>   generates 'count' composites     ")
			fmt.Println("                               larger or equal than 'min'       ")
			fmt.Println("                               composed of two primes           ")
			fmt.Println("  default for 'min' is 1                                        ")
			fmt.Println("  default for 'count' is 100                                    ")

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

	if min.Cmp(one) == -1 {
		panic("min < 1")
	}

	if count < 1 {
		panic("count < 1")
	}



	if primes == true {

		channel := make(chan *big.Int)

		go generatePrimes(min, count, channel)

		for n := range channel {
			fmt.Println(n)
		}

	} else {
		channel := make(chan []*big.Int)

		go generateComposites(min, count, channel)

		for x := range channel {
			fmt.Println(x[0],"=", x[1],"*", x[2])
		}

	}


}

