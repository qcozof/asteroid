package utils

import (
	"fmt"
	"time"
)

func Countdown(seconds int64) {
	seconds--
	fmt.Printf("%d\r", seconds)
	time.Sleep(time.Second)
	if seconds > 0 {
		Countdown(seconds)
	}
}
