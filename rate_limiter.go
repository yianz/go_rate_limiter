package ratelimiter

import (
	"time"
)

type RateLimiter interface {
	Acquire() (*Token, error)
	Release(*Token) error
}

type Config struct {
	// Type of rate limiter
	Type RateLimiterType
	// Max number of tokens
	Limit int
	// Throttles the number of resources that can be issues within a window interval
	Throttle time.Duration
	// The duration of how long the rate limit exists before being reset
	ResetsIn time.Duration
}

type RateLimiterType uint32

const (
	ThrottleType RateLimiterType = iota + 1
	ConcurrencyType
	WindowIntervalType
)

func New(conf *Config) (RateLimiter, error) {
	m, err := NewManager(conf)

	if err != nil {
		return m, err
	}

	switch conf.Type {
	case ThrottleType:
		err = withThrottle(m)
	case ConcurrencyType:
		err = withMaxConcurrency(m)
	case WindowIntervalType:
		err = withWindowInterval(m)
	default:
		err = ErrInvalidLimiterType
	}

	if err != nil {
		return nil, err
	}

	return m, err
}
