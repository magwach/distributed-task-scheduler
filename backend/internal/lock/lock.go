package lock

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/magwach/distributed-task-scheduler/backend/internal/queue"
)

func AquireLock(taskID string) (bool, string) {

	value := uuid.New().String()

	ok, err := queue.RedisClient.SetNX(
		context.Background(),
		fmt.Sprintf("lock:task:%v", taskID),
		value,
		30*time.Second,
	).Result()

	if err != nil {
		log.Println("Error trying to acquire lock")
		return false, ""
	}

	if !ok {
		log.Println("Lock is held by someone else")
		return false, ""
	}

	return ok, value
}

func ReleaseLock(taskID string, lockValue string) error {
	ctx := context.Background()

	key := fmt.Sprintf("lock:task:%v", taskID)

	luaScript := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
	`

	result, err := queue.RedisClient.Eval(
		ctx,
		luaScript,
		[]string{key},
		lockValue,
	).Result()

	if err != nil {
		return err
	}

	if result.(int64) == 0 {
		log.Println("Lock not released (not owner or already expired):", taskID)
	}

	return nil
}
