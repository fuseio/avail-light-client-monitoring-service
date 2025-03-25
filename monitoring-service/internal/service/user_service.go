package service

import (
    "log"
    "time"
    
    "monitoring-service/internal/database"
)

type UserService struct {
    db     *database.Database
    logger *log.Logger
}

func NewUserService(db *database.Database, logger *log.Logger) *UserService {
    return &UserService{
        db:     db,
        logger: logger,
    }
}

func (s *UserService) Start() {
    s.logger.Println("Starting user service")
    
    s.scanForMissingUsers()
    
    ticker := time.NewTicker(24 * time.Hour)
    go func() {
        for range ticker.C {
            s.scanForMissingUsers()
        }
    }()
}

func (s *UserService) scanForMissingUsers() {
    s.logger.Println("Scanning for missing users")
    
    clients, err := s.db.GetAllClients()
    if err != nil {
        s.logger.Printf("Error getting clients: %v", err)
        return
    }
    s.logger.Printf("Found %d clients/operators", len(clients))
    
    validOperators := make(map[string]bool)
    for _, client := range clients {
        if client.Address == "" || client.Address == "0x" || client.Address == "0x0" || len(client.Address) <= 3 {
            continue
        }
        validOperators[client.Address] = true
    }
    
    var operatorsCreated, operatorsExisting int
    for address := range validOperators {
        user, err := s.db.GetOrCreateUser(address, database.UserTypeOperator)
        if err != nil {
            s.logger.Printf("Error creating operator user for %s: %v", address, err)
        } else if user.CreatedAt.Equal(user.UpdatedAt) {
            operatorsCreated++
        } else {
            operatorsExisting++
        }
        
        if err := s.db.SyncOperatorDelegators(address); err != nil {
            s.logger.Printf("Error syncing delegators for operator %s: %v", address, err)
        }
    }
    
    delegations, err := s.db.GetAllDelegations()
    if err != nil {
        s.logger.Printf("Error getting delegations: %v", err)
        return
    }
    s.logger.Printf("Found %d total delegations", len(delegations))
    
    delegatorMap := make(map[string]bool)
    for _, delegation := range delegations {
        if delegation.FromAddress != "" && delegation.FromAddress != "0x" && 
           delegation.FromAddress != "0x0" && len(delegation.FromAddress) > 3 {
            if validOperators[delegation.ToAddress] {
                delegatorMap[delegation.FromAddress] = true
            } else {
                s.logger.Printf("Skipping delegator %s with invalid operator %s", 
                    delegation.FromAddress, delegation.ToAddress)
            }
        }
    }
    
    s.logger.Printf("Found %d unique delegator addresses with valid operators", len(delegatorMap))
    
    var delegatorsCreated, delegatorsExisting int
    for address := range delegatorMap {
        user, err := s.db.GetOrCreateUser(address, database.UserTypeDelegate)
        if err != nil {
            s.logger.Printf("Error creating delegator user for %s: %v", address, err)
        } else if user.CreatedAt.Equal(user.UpdatedAt) {
            delegatorsCreated++
        } else {
            delegatorsExisting++
        }
        
        if err := s.db.SyncDelegatorOperators(address); err != nil {
            s.logger.Printf("Error syncing operators for delegator %s: %v", address, err)
        }
    }
    
    s.logger.Printf("User scan complete: found %d unique operators (%d new, %d existing), %d unique delegators (%d new, %d existing)",
        len(validOperators), operatorsCreated, operatorsExisting, 
        len(delegatorMap), delegatorsCreated, delegatorsExisting)
} 