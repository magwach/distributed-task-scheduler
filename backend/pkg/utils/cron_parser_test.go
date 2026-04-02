package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseCron_Valid(t *testing.T) {
	parsedTime, err := ParseCron("* * * * *")

	assert.NoError(t, err)

	assert.Greater(t, parsedTime, time.Now())

}

func TestParseCron_Invalid(t *testing.T) {
	_, err := ParseCron("invalid cron")

	assert.Error(t, err)
}
