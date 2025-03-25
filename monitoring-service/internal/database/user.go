package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserType string

const (
	UserTypeOperator UserType = "OPERATOR"
	UserTypeDelegate UserType = "DELEGATOR"
)

type User struct {
	Address            string             `bson:"address"`
	UserType           UserType           `bson:"user_type"`
	TotalClaimedPoints int64              `bson:"total_claimed_points"`
	LastClaimTime      time.Time          `bson:"last_claim_time"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at"`
	Delegators         map[string]int64   `bson:"delegators,omitempty"` // map[delegatorAddress]delegatedAmount
	Operators          map[string]int64   `bson:"operators,omitempty"`  // map[operatorAddress]delegatedAmount
}

func (d *Database) GetOrCreateUser(address string, userType UserType) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := d.users.FindOne(ctx, bson.M{
		"address": address,
		"user_type": userType,
	}).Decode(&user)
	
	if err == nil {
		now := time.Now()
		update := bson.M{
			"$set": bson.M{
				"updated_at": now,
			},
		}
		
		_, updateErr := d.users.UpdateOne(ctx, bson.M{
			"address": address,
			"user_type": userType,
		}, update)
		
		if updateErr != nil {
			d.logger.Printf("Warning: Failed to update user %s (%s): %v", address, userType, updateErr)
		}
		
		user.UpdatedAt = now
		
		if user.Delegators == nil && userType == UserTypeOperator {
			user.Delegators = make(map[string]int64)
			
			_, updateErr = d.users.UpdateOne(ctx, bson.M{
				"address": address,
				"user_type": userType,
			}, bson.M{
				"$set": bson.M{
					"delegators": user.Delegators,
				},
			})
			
			if updateErr != nil {
				d.logger.Printf("Warning: Failed to initialize delegators map for user %s: %v", address, updateErr)
			}
		}
		
		if user.Operators == nil && userType == UserTypeDelegate {
			user.Operators = make(map[string]int64)
			
			_, updateErr = d.users.UpdateOne(ctx, bson.M{
				"address": address,
				"user_type": userType,
			}, bson.M{
				"$set": bson.M{
					"operators": user.Operators,
				},
			})
			
			if updateErr != nil {
				d.logger.Printf("Warning: Failed to initialize operators map for user %s: %v", address, updateErr)
			}
		}
		
		return &user, nil
	}

	if err != mongo.ErrNoDocuments {
		// Unexpected error
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	now := time.Now()
	newUser := User{
		Address:            address,
		UserType:           userType,
		TotalClaimedPoints: 0,
		LastClaimTime:      time.Time{}, // Zero time
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	
	if userType == UserTypeOperator {
		newUser.Delegators = make(map[string]int64)
	}

	if userType == UserTypeDelegate {
		newUser.Operators = make(map[string]int64)
	}

	result, err := d.users.InsertOne(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	
	d.logger.Printf("Created new user record: %s (type: %s) with ID: %v", address, userType, result.InsertedID)

	return &newUser, nil
}

func (d *Database) GetAllDelegations() ([]DelegationRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := d.delegations.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delegations: %v", err)
	}
	defer cursor.Close(ctx)

	var delegations []DelegationRecord
	if err := cursor.All(ctx, &delegations); err != nil {
		return nil, fmt.Errorf("failed to decode delegations: %v", err)
	}

	return delegations, nil
}

func (d *Database) GetOperatorDelegators(operatorAddress string) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var user User
	err := d.users.FindOne(ctx, bson.M{
		"address": operatorAddress,
		"user_type": UserTypeOperator,
	}).Decode(&user)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return make(map[string]int64), nil // Return empty map if operator not found
		}
		return nil, fmt.Errorf("failed to find operator user: %v", err)
	}
	
	if user.Delegators == nil {
		user.Delegators = make(map[string]int64)
		
		_, updateErr := d.users.UpdateOne(ctx, bson.M{
			"address": operatorAddress,
			"user_type": UserTypeOperator,
		}, bson.M{
			"$set": bson.M{
				"delegators": user.Delegators,
			},
		})
		
		if updateErr != nil {
			d.logger.Printf("Warning: Failed to initialize delegators map for user %s: %v", operatorAddress, updateErr)
		}
	}
	
	return user.Delegators, nil
}


func (d *Database) SyncOperatorDelegators(operatorAddress string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	operatorUser, err := d.GetOrCreateUser(operatorAddress, UserTypeOperator)
	if err != nil {
		return fmt.Errorf("failed to get operator user: %v", err)
	}
	
	if operatorUser.Delegators == nil {
		operatorUser.Delegators = make(map[string]int64)
	}
	
	cursor, err := d.delegations.Find(ctx, bson.M{"to_address": operatorAddress})
	if err != nil {
		return fmt.Errorf("failed to fetch delegations: %v", err)
	}
	defer cursor.Close(ctx)
	
	operatorUser.Delegators = make(map[string]int64)
	
	var delegations []DelegationRecord
	if err = cursor.All(ctx, &delegations); err != nil {
		return fmt.Errorf("failed to decode delegations: %v", err)
	}
	
	// Update the delegators map with current delegations
	for _, delegation := range delegations {
		if delegation.FromAddress != "" && delegation.FromAddress != "0x" && 
		   delegation.FromAddress != "0x0" && len(delegation.FromAddress) > 3 {
			operatorUser.Delegators[delegation.FromAddress] = delegation.Amount
		}
	}
	
	_, err = d.users.UpdateOne(
		ctx,
		bson.M{
			"address": operatorAddress,
			"user_type": UserTypeOperator,
		},
		bson.M{
			"$set": bson.M{
				"delegators": operatorUser.Delegators,
				"updated_at": time.Now(),
			},
		},
	)
	
	if err != nil {
		return fmt.Errorf("failed to update operator's delegators map: %v", err)
	}
	
	d.logger.Printf("Synced delegators for operator %s: %d delegators found", 
		operatorAddress, len(operatorUser.Delegators))
	
	return nil
}

func (d *Database) SyncDelegatorOperators(delegatorAddress string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	delegatorUser, err := d.GetOrCreateUser(delegatorAddress, UserTypeDelegate)
	if err != nil {
		return fmt.Errorf("failed to get delegator user: %v", err)
	}
	
	if delegatorUser.Operators == nil {
		delegatorUser.Operators = make(map[string]int64)
	}
	
	cursor, err := d.delegations.Find(ctx, bson.M{"from_address": delegatorAddress})
	if err != nil {
		return fmt.Errorf("failed to fetch delegations: %v", err)
	}
	defer cursor.Close(ctx)
	
	delegatorUser.Operators = make(map[string]int64)
	
	var delegations []DelegationRecord
	if err = cursor.All(ctx, &delegations); err != nil {
		return fmt.Errorf("failed to decode delegations: %v", err)
	}
	
	for _, delegation := range delegations {
		if delegation.ToAddress != "" && delegation.ToAddress != "0x" && 
		   delegation.ToAddress != "0x0" && len(delegation.ToAddress) > 3 {
			
			exists, err := d.ClientExists(delegation.ToAddress)
			if err != nil {
				d.logger.Printf("Warning: Error checking if client exists: %v", err)
				continue
			}
			
			if !exists {
				d.logger.Printf("Skipping invalid operator %s for delegator %s", 
					delegation.ToAddress, delegatorAddress)
				continue
			}
			
			delegatorUser.Operators[delegation.ToAddress] = delegation.Amount
		}
	}
	
	_, err = d.users.UpdateOne(
		ctx,
		bson.M{
			"address": delegatorAddress,
			"user_type": UserTypeDelegate,
		},
		bson.M{
			"$set": bson.M{
				"operators": delegatorUser.Operators,
				"updated_at": time.Now(),
			},
		},
	)
	
	if err != nil {
		return fmt.Errorf("failed to update delegator's operators map: %v", err)
	}
	
	d.logger.Printf("Synced operators for delegator %s: %d valid operators found", 
		delegatorAddress, len(delegatorUser.Operators))
	
	return nil
}

func (d *Database) GetDelegatorOperators(delegatorAddress string) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var user User
	err := d.users.FindOne(ctx, bson.M{
		"address": delegatorAddress,
		"user_type": UserTypeDelegate,
	}).Decode(&user)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return make(map[string]int64), nil // Return empty map if delegator not found
		}
		return nil, fmt.Errorf("failed to find delegator user: %v", err)
	}
	
	if user.Operators == nil {
		user.Operators = make(map[string]int64)
		
		_, updateErr := d.users.UpdateOne(ctx, bson.M{
			"address": delegatorAddress,
			"user_type": UserTypeDelegate,
		}, bson.M{
			"$set": bson.M{
				"operators": user.Operators,
			},
		})
		
		if updateErr != nil {
			d.logger.Printf("Warning: Failed to initialize operators map for user %s: %v", delegatorAddress, updateErr)
		}
	}
	
	return user.Operators, nil
} 