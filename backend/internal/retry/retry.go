package retry

import "math"

func Delay(retry_delay_seconds, retryCount int) int {
	return retry_delay_seconds * int(math.Pow(2, float64(retryCount)))
}
