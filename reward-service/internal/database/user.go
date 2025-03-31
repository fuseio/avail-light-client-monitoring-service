package database

import (
	"time"
)

type UserType string

const (
	UserTypeOperator UserType = "OPERATOR"
	UserTypeDelegate UserType = "DELEGATOR"
)

// User represents either an operator or delegator
type User struct {
	Address              string             `bson:"address"`
	UserType             UserType           `bson:"user_type"`
	TotalClaimedPoints   int64              `bson:"total_claimed_points"`
	LastClaimTime        time.Time          `bson:"last_claim_time"`
	CreatedAt            time.Time          `bson:"created_at"`
	UpdatedAt            time.Time          `bson:"updated_at"`
	Delegators           map[string]int64   `bson:"delegators,omitempty"` // map[delegatorAddress]delegatedAmount
	Operators            map[string]int64   `bson:"operators,omitempty"`  // map[operatorAddress]delegatedAmount
	AllUptimePercentage  float64            `bson:"all_uptime_percentage,omitempty"`
	Status               string             `bson:"status,omitempty"`
}