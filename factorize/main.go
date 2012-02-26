package main

import (
	"math/big"
	"fmt"
	"os"

	"github.com/hydroo/quadratic-sieve/misc"
)


func main() {

	args := os.Args[1:]

	var min *big.Int
	var step *big.Int
	var count *big.Int

	for i := 0; i < len(args); i++ {

		if args[i] == "-h" || args[i] == "--help" {

			fmt.Println("factor <min> <step> <count>")
			fmt.Println("")
			fmt.Println("    breaks n down into two factors")
			fmt.Println("")
			fmt.Println("    default is 1 1 +inf")
			os.Exit(0)

		} else {

			if args[i][0] != '-' {

				var boolerr bool

				if min == nil {
					min = big.NewInt(0)
					min, boolerr = min.SetString(args[i], 10)
				} else if step == nil {
					step = big.NewInt(0)
					step, boolerr = step.SetString(args[i], 10)
				} else if count == nil {
					count = big.NewInt(0)
					count, boolerr = count.SetString(args[i], 10)
				} else {
					panic("too many arguments")
				}

				if boolerr == false {
					fmt.Println("not a valid number: ", args[i])
					os.Exit(-1)
				}

			} else {

				fmt.Println("unknown argument: ", args[i])
				os.Exit(-1)

			}
		}	
	}

	if min == nil {
		min = big.NewInt(1)
	}

	if step == nil {
		step = big.NewInt(1)
	}

	if count == nil {
		count = big.NewInt(-1)
	}


	for i := min; count.Cmp(misc.Zero) != 0; i.Add(i,step) {

		x, y := factorize(i)

		if x != nil && y != nil {
			fmt.Println(i, "+", x, y)
		} else {
			fmt.Println(i, "-")
		}

		count.Sub(count, misc.One)
	}

}
