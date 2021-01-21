package rate_limiter

import (
	"proxytest/pkg/config"
	"time"
)

type RateLimiter interface {
	Check(clientID string) bool
}

type accessData struct {
	count     int
	createdAt time.Time
}

func newAccessData(now time.Time) *accessData {
	return &accessData{
		count:     1,
		createdAt: now,
	}
}

type inMemoryRateLimiter struct {
	cfg       config.RateLimitConfig
	countData map[string]*accessData
}

func (irl *inMemoryRateLimiter) Check(clientID string) bool {
	now := time.Now()

	ad, ok := irl.countData[clientID]
	if !ok || now.Sub(ad.createdAt).Seconds() > float64(irl.cfg.Duration()) {
		irl.countData[clientID] = newAccessData(now)
		return true
	}

	if ad.count < irl.cfg.MaxRequest() {
		irl.countData[clientID].count++
		return true
	}

	return false
}

func NewRateLimiter(cfg config.RateLimitConfig) RateLimiter {
	return &inMemoryRateLimiter{
		cfg:       cfg,
		countData: make(map[string]*accessData),
	}
}
