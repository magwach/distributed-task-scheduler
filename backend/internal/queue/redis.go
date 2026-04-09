package queue

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/magwach/distributed-task-scheduler/backend/pkg/utils"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

var RedisClient *redis.Client

const TaskQueueKey = "task:queue"

func InitRedis(redisUrl string) {
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalf("Invalid Redis URL: %v", err)
	}

	if opt.TLSConfig == nil && len(redisUrl) >= 8 && redisUrl[:8] == "rediss://" {
		opt.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	RedisClient = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
}

func Enqueue(taskID, taskPriority string) error {

	score, err := utils.PriorityScores(taskPriority)

	if err != nil {
		return err
	}

	err = RedisClient.ZAdd(context.Background(), "task:queue", redis.Z{
		Score:  float64(score),
		Member: taskID,
	}).Err()

	if err != nil {
		log.Println("Failed to set the value: ", err)
		return err
	}

	return nil

}

func Dequeue(timeout time.Duration) (string, error) {

	result, err := RedisClient.BZPopMax(Ctx, timeout, TaskQueueKey).Result()

	if err != nil {
		return "", err
	}

	taskId, ok := result.Member.(string)
	if !ok {
		return "", fmt.Errorf("failed to cast task ID")
	}
	
	return taskId, nil
}

func GetRedisClient() *redis.Client {
	return RedisClient
}
