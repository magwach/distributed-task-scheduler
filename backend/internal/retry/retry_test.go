package retry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelay(t *testing.T) {
	delay := Delay(60, 0)

	assert.Equal(t, 60, delay)

	delay = Delay(60, 1)

	assert.Equal(t, 120, delay)

	delay = Delay(60, 2)

	assert.Equal(t, 240, delay)
}
