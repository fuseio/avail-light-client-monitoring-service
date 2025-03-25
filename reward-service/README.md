# Reward Distribution Service

This service is responsible for calculating and distributing rewards to node operators and delegators based on their performance and delegations.

## Overview

The Reward Distribution Service is part of a microservice architecture that works alongside the Monitoring Service. While the Monitoring Service tracks node activity, NFT ownership, and delegations, the Reward Service focuses on calculating and distributing rewards based on this data.

## Reward Calculation

### Node Operator Points

- Base Points: 1000 points if the node was active for at least 50% of the past 24 hours.
- Commission Points: Additional points from delegators based on the commission rate.

Formula: `OperatorPoints = 1000 + (1000 × DelegatedNFTs × CommissionRate / 100)`

Example with 100 NFTs (1 owned and 99 delegated) and 10% commission:
`OperatorPoints = 1000 + (1000 * 99 * 0.10) = 1000 + 9900 = 10900`

### Delegator Points

- Base Points: 1000 points per NFT if the delegated node operator is eligible for points.
- Commission Deduction: Points are reduced by the operator's commission rate.

Formula: `DelegatorPoints = (1000 * DelegatedNFTs) - (1000 * DelegatedNFTs * CommissionRate / 100)`

Example with 99 NFTs delegated and 10% commission:
`DelegatorPoints = (1000 * 99) - (1000 * 99 * 0.10) = 99000 - 9900 = 89100`

## API Endpoints

### GET /clients

Retrieves all clients with their points.

### GET /clients/{address}/rewards

Retrieves the reward history for a specific client.

### GET /clients/{address}/points

Retrieves the global points for a specific client.

### GET /rewards/summary

Retrieves a summary of all rewards, including:

- Total number of rewards distributed
- Total number of users
- Total points distributed
- Latest reward timestamp
- Average points per user

### GET /health

Health check endpoint.

## Logging

The service includes comprehensive logging for all operations:

- **Reward Calculation Logs**: Detailed logs of the reward calculation process, including uptime checks, commission calculations, and point distributions.
- **API Request Logs**: Logs for all API requests, including request details, response times, and status codes.
- **Database Operation Logs**: Logs for database operations, including reward storage and user point updates.
- **Service Lifecycle Logs**: Logs for service startup, shutdown, and periodic reward calculations.

## Configuration

The service can be configured using environment variables:

- `MONGO_URI`: MongoDB connection URI
- `MONGO_DB`: MongoDB database name
- `PORT`: Server port
- `RPC_URL`: Blockchain RPC URL
- `NFT_CONTRACT_ADDRESS`: NFT contract address
- `DELEGATE_CONTRACT_ADDRESS`: Delegation contract address
- `RIGHTS`: Rights for the NFT
- `CHECK_NFT_INTERVAL`: Interval in minutes for reward calculation

## Running the Service

### Using Docker

```bash
docker-compose up
```

### Manually

```bash
cd reward-service
go run cmd/main.go
```
