//go:build integration
// +build integration

package lock

import (
	"context"
	"os"
	"testing"

	"github.com/magwach/distributed-task-scheduler/backend/internal/queue"
	"github.com/stretchr/testify/assert"
)

func setupRedis() {
	redisAddress := os.Getenv("REDIS_ADDR")
	queue.InitRedis(redisAddress)
}

func TestAcquireLock(t *testing.T) {
	setupRedis()

	taskID := "test-task-id"

	ok, lockValue := AquireLock(taskID)

	assert.True(t, ok)
	assert.NotEmpty(t, lockValue)

	queue.RedisClient.Del(context.Background(), "lock:task:"+taskID)
}

func TestAcquireLock_AlreadyLocked(t *testing.T) {
	setupRedis()

	taskID := "test-task-id"

	ok1, _ := AquireLock(taskID)
	assert.True(t, ok1)

	ok2, _ := AquireLock(taskID)
	assert.False(t, ok2)

	queue.RedisClient.Del(context.Background(), "lock:task:"+taskID)
}

func TestReleaseLock(t *testing.T) {
	setupRedis()

	taskID := "test-task-id"

	_, lockValue := AquireLock(taskID)

	err := ReleaseLock(taskID, lockValue)
	assert.NoError(t, err)

	ok, _ := AquireLock(taskID)
	assert.True(t, ok)

	queue.RedisClient.Del(context.Background(), "lock:task:"+taskID)
}
