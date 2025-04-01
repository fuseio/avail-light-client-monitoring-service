package uptime

import (
	"context"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// sampling interval (5 minutes)
	IntervalDuration = 5 * time.Minute
	// 60 days
	
	MaxAllTimeHistoryDuration = 60 * 24 * time.Hour
	//  7 days
	WeeklyHistoryDuration = 7 * 24 * time.Hour
)

type Calculator struct {
	heartbeats *mongo.Collection
}

// NewCalculator creates a new uptime calculator
func NewCalculator(heartbeatsCollection *mongo.Collection) *Calculator {
	return &Calculator{
		heartbeats: heartbeatsCollection,
	}
}

// GetUptimePercentages calculates both all-time and weekly uptime percentages
func (c *Calculator) GetUptimePercentages(ctx context.Context, clientAddress string, createdAt time.Time) (allTime, weekly float64, err error) {
	allTime, err = c.CalculateAllTimeUptime(ctx, clientAddress, createdAt)
	if err != nil {
		return 0, 0, err
	}

	weekly, err = c.CalculateWeeklyUptime(ctx, clientAddress, createdAt)
	if err != nil {
		return allTime, 0, err
	}

	return allTime, weekly, nil
}

// CalculateAllTimeUptime calculates the uptime percentage based on 5-minute intervals
func (c *Calculator) CalculateAllTimeUptime(ctx context.Context, clientAddress string, createdAt time.Time) (float64, error) {
	// Cap all-time history at 60 days max
	startTime := createdAt
	maxStart := time.Now().Add(-MaxAllTimeHistoryDuration)
	if startTime.Before(maxStart) {
		startTime = maxStart
	}

	return c.calculateUptimeForPeriod(ctx, clientAddress, startTime)
}

// CalculateWeeklyUptime calculates the uptime percentage for the last 7 days
func (c *Calculator) CalculateWeeklyUptime(ctx context.Context, clientAddress string, createdAt time.Time) (float64, error) {
	// For weekly, use the later of: a week ago or client creation time
	startTime := time.Now().Add(-WeeklyHistoryDuration)
	if createdAt.After(startTime) {
		startTime = createdAt
	}

	return c.calculateUptimeForPeriod(ctx, clientAddress, startTime)
}

// calculateUptimeForPeriod counts the number of intervals with at least one heartbeat
// and calculates the uptime percentage
func (c *Calculator) calculateUptimeForPeriod(ctx context.Context, clientAddress string, startTime time.Time) (float64, error) {
	// Calculate expected intervals
	duration := time.Since(startTime)
	expectedIntervals := int64(duration / IntervalDuration)
	if expectedIntervals == 0 {
		expectedIntervals = 1 // Avoid division by zero
	}

	pipeline := mongo.Pipeline{
		// Stage 1: Match documents by client and time range
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "client_address", Value: clientAddress},
			{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
		}}},
		
		// Stage 2: Project only necessary fields with pre-calculated interval bucket
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "interval_bucket", Value: bson.D{
				{Key: "$floor", Value: bson.D{
					{Key: "$divide", Value: bson.A{
						bson.D{{Key: "$toLong", Value: "$timestamp"}},
						int64(IntervalDuration / time.Millisecond),
					}},
				}},
			}},
		}}},
		
		// Stage 3: Group by the interval bucket to count unique intervals
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$interval_bucket"},
		}}},
		
		// Stage 4: Count total intervals
		bson.D{{Key: "$count", Value: "total_intervals"}},
	}

	opts := options.Aggregate().SetMaxTime(3 * time.Second)

	cursor, err := c.heartbeats.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalIntervals int64 `bson:"total_intervals"`
	}

	totalIntervals := int64(0)
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		totalIntervals = result.TotalIntervals
	}

	// Calculate uptime percentage
	uptimePercentage := float64(totalIntervals) / float64(expectedIntervals) * 100
	return math.Min(uptimePercentage, 100), nil
} 