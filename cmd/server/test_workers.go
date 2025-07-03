package main

import (
	"context"
	"log"
	"time"

	"github.com/aungmyozaw92/go-api-setup/internal/repository/mocks"
	"github.com/aungmyozaw92/go-api-setup/internal/worker"
	"github.com/stretchr/testify/mock"
)

// testSimpleApproach demonstrates Approach A
func testSimpleApproach() {
	log.Println("\n🔹 TESTING APPROACH A: Simple Worker")
	
	// Mock repository for testing
	mockRepo := &mocks.MockUserRepository{}
	mockRepo.On("Count", mock.MatchedBy(func(ctx context.Context) bool { return true })).Return(int64(42), nil)
	
	// Simple approach - each worker started manually
	userMonitor := worker.NewUserMonitor(mockRepo)
	go userMonitor.StartUserCountMonitoring()
	
	emailWorker := worker.NewEmailWorker()
	go emailWorker.Start()
	
	log.Println("✅ Started workers manually")
	
	// Let it run for a bit
	time.Sleep(5 * time.Second)
	
	// No easy way to stop workers in simple approach
	log.Println("❌ Difficult to stop workers gracefully")
}

// testManagerApproach demonstrates Approach B
func testManagerApproach() {
	log.Println("\n🔹 TESTING APPROACH B: Manager Pattern")
	
	// Mock repository for testing
	mockRepo := &mocks.MockUserRepository{}
	mockRepo.On("Count", mock.MatchedBy(func(ctx context.Context) bool { return true })).Return(int64(42), nil)
	
	// Manager approach - all workers managed centrally
	workerManager := worker.SetupDefaultWorkers(mockRepo)
	workerManager.StartAll()
	
	// Let it run for a bit
	time.Sleep(5 * time.Second)
	
	// Easy to stop all workers
	workerManager.StopAll()
	log.Println("✅ All workers stopped gracefully")
}

// Uncomment main function to test
/*
func main() {
	testSimpleApproach()
	testManagerApproach()
}
*/ 