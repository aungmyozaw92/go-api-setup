package worker

import (
	"log"
	"time"
)

// EmailWorker handles email processing
type EmailWorker struct {
	ticker *time.Ticker
	done   chan bool
}

// NewEmailWorker creates a new email worker
func NewEmailWorker() *EmailWorker {
	return &EmailWorker{
		done: make(chan bool),
	}
}

// Start begins email processing (implements Worker interface)
func (w *EmailWorker) Start() {
	w.ticker = time.NewTicker(30 * time.Second)
	
	log.Println("📧 Starting email worker (every 30 seconds)")
	
	for {
		select {
		case <-w.ticker.C:
			// Example: Process pending emails
			log.Printf("📧 Processing pending emails...")
			// Add your email logic here
		case <-w.done:
			log.Println("🛑 Stopping email worker")
			return
		}
	}
}

// Stop gracefully stops the email worker (implements Worker interface)
func (w *EmailWorker) Stop() {
	if w.ticker != nil {
		w.ticker.Stop()
	}
	close(w.done)
}

// Name returns the worker name (implements Worker interface)
func (w *EmailWorker) Name() string {
	return "EmailWorker"
} 