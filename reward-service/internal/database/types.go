package database

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// RewardRecord represents a reward given to a user
type RewardRecord struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserAddress    string             `bson:"user_address" json:"user_address"`
    Points         int64              `bson:"points" json:"points"`
    CycleID        string             `bson:"cycle_id" json:"cycle_id"`
    IsClaimed      bool               `bson:"is_claimed" json:"is_claimed"`
    CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
    ClaimedAt      time.Time          `bson:"claimed_at,omitempty" json:"claimed_at,omitempty"`
    OperatorAddress string            `bson:"operator_address,omitempty" json:"operator_address,omitempty"`
    Type           string             `bson:"type" json:"type"`
}

// RewardSummary contains summary information about rewards
type RewardSummary struct {
    CycleID         string    `json:"cycle_id" bson:"cycle_id"`
    TotalRewards    int       `json:"total_rewards" bson:"total_rewards"`
    TotalPoints     int64     `json:"total_points" bson:"total_points"`
    OperatorRewards int       `json:"operator_rewards" bson:"operator_rewards"`
    DelegatorRewards int      `json:"delegator_rewards" bson:"delegator_rewards"`
    Timestamp       time.Time `json:"timestamp" bson:"timestamp"`
}

// MonitoringClientInfo mirrors the ClientInfo struct from monitoring service
type MonitoringClientInfo struct {
    Address                string    `bson:"address"`
    TotalTime              int64     `bson:"total_time"`
    LastHeartbeat          time.Time `bson:"last_heartbeat"`
    CreatedAt              time.Time `bson:"created_at"`
    NFTAmount              int64     `bson:"nft_amount"`
    CommissionRate         float64   `bson:"commission_rate"`
    Status                 string    `bson:"status"`
    AllUptimePercentage    float64   `bson:"all_uptime_percentage"`
    WeeklyUptimePercentage float64   `bson:"weekly_uptime_percentage"`
    OperatorName           string    `bson:"operator_name"`
    RewardCollectorAddress string    `bson:"reward_collector_address"`
}

// DelegationRecord represents a delegation from one address to another
type DelegationRecord struct {
    ID             primitive.ObjectID `bson:"_id,omitempty"`
    FromAddress    string             `bson:"from_address"`
    ToAddress      string             `bson:"to_address"`
    Amount         int64              `bson:"amount"`
    CommissionRate float64            `bson:"commission_rate"`
    Timestamp      time.Time          `bson:"timestamp"`
}

// User type from monitoring service
type User struct {
    Address              string             `bson:"address"`
    UserType             string             `bson:"user_type"`
    TotalClaimedPoints   int64              `bson:"total_claimed_points"`
    LastClaimTime        time.Time          `bson:"last_claim_time"`
    CreatedAt            time.Time          `bson:"created_at"`
    UpdatedAt            time.Time          `bson:"updated_at"`
    Delegators           map[string]int64   `bson:"delegators,omitempty"`
    Operators            map[string]int64   `bson:"operators,omitempty"`
    AllUptimePercentage  float64            `bson:"all_uptime_percentage,omitempty"`
    Status               string             `bson:"status,omitempty"`
    CommissionRate       float64            `bson:"commission_rate,omitempty"`
    RewardAddress        string             `bson:"reward_address,omitempty"`
} 