package rate_limiter_test

import (
	"github.com/stretchr/testify/suite"
	"proxytest/pkg/config"
	"proxytest/pkg/rate_limiter"
	"proxytest/pkg/test"
	"testing"
	"time"
)

type rateLimiterSite struct {
	rl rate_limiter.RateLimiter
	suite.Suite
}

func (rls *rateLimiterSite) SetupSuite() {
	mockRateLimiterConfig := &config.MockRateLimitConfig{}
	mockRateLimiterConfig.On("MaxRequest").Return(50)
	mockRateLimiterConfig.On("Duration").Return(2)

	rls.rl = rate_limiter.NewRateLimiter(mockRateLimiterConfig)
}

func TestRateLimiter(t *testing.T) {
	suite.Run(t, new(rateLimiterSite))
}

func (rls *rateLimiterSite) TestCheckBelowMaxLimit() {
	rls.Assert().True(rls.rl.Check(test.RandString(8)))
}

func (rls *rateLimiterSite) TestCheckExceededMaxLimit() {
	clientID := test.RandString(8)

	for i := 0; i < 50; i++ {
		rls.Assert().True(rls.rl.Check(clientID))
	}

	rls.Assert().False(rls.rl.Check(clientID))
}

func (rls *rateLimiterSite) TestCheckResetAfterTimeLimit() {
	clientID := test.RandString(8)

	for i := 0; i < 50; i++ {
		rls.Assert().True(rls.rl.Check(clientID))
	}

	time.Sleep(time.Second * 2)

	rls.Assert().True(rls.rl.Check(clientID))
}
