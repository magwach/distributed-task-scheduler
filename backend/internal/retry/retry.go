package retry

func Delay(retry_delay_seconds, retryCount int) int {
	return retry_delay_seconds*2 ^ retryCount
}
