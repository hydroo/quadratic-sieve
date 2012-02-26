package main

import (
	"fmt"
)

func nanoSecondsToString(n int64) string {
	return fmt.Sprintf("%f", float64(n)/1000000000.0)
}


