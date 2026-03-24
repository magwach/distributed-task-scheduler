package queue

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

var RedisClient *redis.Client

const TaskQueueKey = "task:queue"

func InitRedis(redisUrl string) {

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	_, err := RedisClient.Ping(Ctx).Result()

	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	log.Println("Connected to Redis")

}

func Enqueue(taskID string) error {

	err := RedisClient.LPush(Ctx, TaskQueueKey, taskID).Err()

	if err != nil {
		log.Println("Failed to set the value: ", err)
		return err
	}

	return nil

}

func Dequeue(timeout time.Duration) (string, error) {

	result, err := RedisClient.BRPop(Ctx, timeout, TaskQueueKey).Result()

	if err != nil {
		return "", err
	}

	taskId := result[1]

	return taskId, nil
}
