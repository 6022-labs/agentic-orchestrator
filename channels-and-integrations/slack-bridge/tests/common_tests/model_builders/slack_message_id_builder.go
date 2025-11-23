package model_builders

import (
	"fmt"
	"time"
)

func GenerateSlackMessageID(timestamp time.Time) string {
	timestampNano := timestamp.UnixNano()
	sec := timestampNano / int64(time.Second)
	nsec := timestampNano % int64(time.Second)
	return fmt.Sprintf("%d.%06d", sec, nsec/1000)
}

func GenerateRandomSlackMessageID() string {
	timestamp := time.Now()
	return GenerateSlackMessageID(timestamp)
}
