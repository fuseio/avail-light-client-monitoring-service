package service

import (
	"log"
	"time"

	"reward-service/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type RewardService struct {
	db     *mongo.Database
	logger *log.Logger
	done   chan struct{}
}

func NewRewardService(db *mongo.Database, logger *log.Logger) *RewardService {
	return &RewardService{
		db:     db,
		logger: logger,
		done:   make(chan struct{}),
	}
}

// Stop signals the service to shut down
func (s *RewardService) Stop() {
	close(s.done)
}

// ScheduleRewards sets up a scheduler to run rewards at 10am Israel time every day
func (s *RewardService) ScheduleRewards() {
	s.logger.Println("Setting up reward scheduler to run at 10am Israel time daily")
	
	// Load Israel timezone
	israelLoc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		s.logger.Printf("Error loading Israel timezone: %v, falling back to UTC", err)
		israelLoc = time.UTC
	}
	
	go func() {
		for {
			now := time.Now().In(israelLoc)
			
			// Calculate next 10am Israel time
			next10am := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, israelLoc)
			if now.After(next10am) || now.Equal(next10am) {
				next10am = next10am.Add(24 * time.Hour)
			}
			
			// Wait until next 10am
			timeToWait := next10am.Sub(now)
			s.logger.Printf("Next reward calculation scheduled for %s (waiting %s)", 
				next10am.Format("2006-01-02 15:04:05 MST"), timeToWait)
			
			select {
			case <-time.After(timeToWait):
				// Create a new processor
				processor := database.NewRewardProcessor(s.db, s.logger)
				
				// Process rewards
				s.logger.Println("Starting scheduled reward processing")
				if err := processor.ProcessRewards(); err != nil {
					s.logger.Printf("Error processing rewards: %v", err)
				} else {
					s.logger.Println("Scheduled reward processing completed successfully")
				}
				
			case <-s.done:
				s.logger.Println("Reward scheduler shutting down")
				return
			}
		}
	}()
}

// ProcessRewardsManually allows for manual reward processing (useful for testing)
func (s *RewardService) ProcessRewardsManually(customCycleID string) error {
	processor := database.NewRewardProcessor(s.db, s.logger)
	
	if customCycleID != "" {
		processor.SetCycleID(customCycleID)
	}
	
	return processor.ProcessRewards()
}