package worker

import (
	"log"
	"sync"

	"github.com/aungmyozaw92/go-api-setup/internal/repository"
)

// Worker interface for all background workers
type Worker interface {
	Start()
	Stop()
	Name() string
}

// Manager handles all background workers
type Manager struct {
	workers []Worker
	wg      sync.WaitGroup
}

// NewManager creates a new worker manager
func NewManager() *Manager {
	return &Manager{
		workers: make([]Worker, 0),
	}
}

// AddWorker adds a worker to the manager
func (m *Manager) AddWorker(worker Worker) {
	m.workers = append(m.workers, worker)
	log.Printf("➕ Added worker: %s", worker.Name())
}

// StartAll starts all registered workers
func (m *Manager) StartAll() {
	log.Println("🚀 Starting all workers...")
	
	for _, worker := range m.workers {
		m.wg.Add(1)
		go func(w Worker) {
			defer m.wg.Done()
			log.Printf("▶️  Starting worker: %s", w.Name())
			w.Start()
		}(worker)
	}
	
	log.Printf("✅ Started %d workers", len(m.workers))
}

// StopAll stops all workers gracefully
func (m *Manager) StopAll() {
	log.Println("🛑 Stopping all workers...")
	
	for _, worker := range m.workers {
		worker.Stop()
		log.Printf("⏹️  Stopped worker: %s", worker.Name())
	}
	
	m.wg.Wait()
	log.Println("✅ All workers stopped")
}

// SetupDefaultWorkers creates default workers for the application
func SetupDefaultWorkers(userRepo repository.UserRepository) *Manager {
	manager := NewManager()
	
	// Add user monitoring worker
	userMonitor := NewUserMonitor(userRepo)
	manager.AddWorker(userMonitor)
	
	// Add email processing worker
	emailWorker := NewEmailWorker()
	manager.AddWorker(emailWorker)
	
	// Add more workers here as needed:
	// cleanupWorker := NewCleanupWorker(db)
	// manager.AddWorker(cleanupWorker)
	
	// analyticsWorker := NewAnalyticsWorker(analyticsRepo)
	// manager.AddWorker(analyticsWorker)
	
	return manager
} 