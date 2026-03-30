package utils

import (
	"errors"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func ParseCron(schedule string) (time.Time, error) {
	parser := cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)

	cronSchedule, err := parser.Parse(schedule)

	if err != nil {
		log.Println("Failed to parse the cron")
		return time.Time{}, errors.New("failed to parse the cron")
	}

	return cronSchedule.Next(time.Now()), nil
}

func PriorityScores(priority string) (int, error) {
	switch priority {
	case "low":
		return 1, nil
	case "normal":
		return 2, nil
	case "high":
		return 3, nil
	default:
		return 0, errors.New("invalid priority")
	}
}
