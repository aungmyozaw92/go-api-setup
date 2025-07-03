package worker

import (
	"context"
	"log"
	"time"

	"github.com/aungmyozaw92/go-api-setup/internal/repository"
)

// UserMonitor handles user count monitoring
type UserMonitor struct {
	userRepo repository.UserRepository
	ticker   *time.Ticker
	done     chan bool
}

// NewUserMonitor creates a new user monitor
func NewUserMonitor(userRepo repository.UserRepository) *UserMonitor {
	return &UserMonitor{
		userRepo: userRepo,
		done:     make(chan bool),
	}
}

// Start begins the user count monitoring (implements Worker interface)
func (m *UserMonitor) Start() {
	m.ticker = time.NewTicker(10 * time.Second)
	
	log.Println("📊 Starting user count monitoring (every 10 seconds)")
	
	for {
		select {
		case <-m.ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			count, err := m.userRepo.Count(ctx)
			cancel()
			
			if err != nil {
				log.Printf("❌ Error getting user count: %v", err)
			} else {
				log.Printf("👥 Current user count: %d", count)
			}
		case <-m.done:
			log.Println("🛑 Stopping user count monitoring")
			return
		}
	}
}

// Stop gracefully stops the user monitoring (implements Worker interface)
func (m *UserMonitor) Stop() {
	if m.ticker != nil {
		m.ticker.Stop()
	}
	close(m.done)
}

// Name returns the worker name (implements Worker interface)
func (m *UserMonitor) Name() string {
	return "UserCountMonitor"
}

// StartUserCountMonitoring runs a background goroutine that logs user count every 10 seconds
// Deprecated: Use Start() method instead for better control
func (m *UserMonitor) StartUserCountMonitoring() {
	m.Start()
} 