package ratelimiter

import (
	"sync/atomic"
	"time"
)

func withMaxConcurrency(m *Manager) error {
	go func() {
		ticker := time.NewTicker(time.Second * 10)

		for {
			<-ticker.C
			now := time.Now().UTC()
			for _, token := range m.activeTokens {
				if token.ExpiresAt.Before(now) {
					m.Release(token)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-m.inChan:
				if !m.hasRemaining() {
					atomic.AddInt64(&m.awaiting, 1)
					continue
				}
				m.generateToken()
			case <-m.releaseChan:
				continue
			}
		}
	}()

	return nil
}
