package utils

import (
	"errors"
	"fmt"
	"math"
	"time"
)

func RetryNTimes(maxRetries int, wait int, fn func() (interface{}, error)) (interface{}, error) {
	for i := 0; i < maxRetries; i++ {
		val, err := fn()
		if err == nil {
			return val, nil
		}

		fmt.Printf("Attempt %d failed: %s\n", i+1, err)

		// Wait for 2^i seconds before the next attempt
		time.Sleep(time.Second * time.Duration(math.Pow(float64(wait), float64(i))))
	}

	return nil, errors.New("maximum number of retries reached")
}
