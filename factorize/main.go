package main

import (
	"math/big"
	"fmt"
	"os"
)


func main() {

	var helpText string
	helpText += "factor [options] <min> <step>                            \n"
	helpText += "                                                         \n"
	helpText += "    breaks n down into two factors starting at min       \n"
	helpText += "                                                         \n"
	helpText += "  --benchmark     print out additonal timing information \n"
	helpText += "                                                         \n"
	helpText += "    default is 1 1                                      \n"

	args := os.Args[1:]

	var min *big.Int
	var step *big.Int
	benchmark := false

	for i := 0; i < len(args); i++ {

		if args[i] == "-h" || args[i] == "--help" {

			fmt.Print(helpText)
			os.Exit(0)

		} else if args[i] == "--benchmark" {

			benchmark = true

		} else {

			if args[i][0] != '-' {

				var boolerr bool

				if min == nil {
					min = big.NewInt(0)
					min, boolerr = min.SetString(args[i], 10)
				} else if step == nil {
					step = big.NewInt(0)
					step, boolerr = step.SetString(args[i], 10)
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

	for i := min;; i.Add(i,step) {

		factorize(i, benchmark)
	}

}

