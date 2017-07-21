package main

import (
	"big"
	//"fmt"
)


func generateComposites(min *big.Int, count int, returnChannel chan<- []*big.Int) {

	firstPrimeChannel := make(chan *big.Int)
	secondPrimeChannel := make(chan *big.Int)


	leftStart := big.NewInt(2)
	leftStart.Exp(leftStart, big.NewInt(int64(min.BitLen())), nil)
	leftStart.Rsh(leftStart, uint(min.BitLen()*3/4))

	rightStart := big.NewInt(2)
	rightStart.Exp(rightStart, big.NewInt(int64(min.BitLen())), nil)
	rightStart.Rsh(rightStart, uint(min.BitLen()/4))

	//fmt.Println("l", leftStart)
	//fmt.Println("r", rightStart)
	//fmt.Println()

	go generatePrimes(leftStart, count, firstPrimeChannel)
	go generatePrimes(rightStart, count, secondPrimeChannel)

	for {
		ret := make([]*big.Int, 3)

		ret[2] = <-firstPrimeChannel
		ret[1] = <-secondPrimeChannel


		ret[0] = big.NewInt(0)
		ret[0].Mul(ret[2], ret[1])

		returnChannel <- ret
	}


}

