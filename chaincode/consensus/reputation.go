package consensus

import (
    "math"
    "time"
)

type ReputationMetrics struct {
    TransactionSuccess float64
    ResponseTime      float64
    UptimePercentage  float64
    DataQuality       float64
}

// CalculateReputation computes the reputation score for a node
func CalculateReputation(metrics ReputationMetrics) float64 {
    // Weights for different metrics
    const (
        wTransaction = 0.4
        wResponse    = 0.2
        wUptime      = 0.2
        wQuality     = 0.2
    )

    // Calculate weighted score
    score := (metrics.TransactionSuccess * wTransaction) +
             (metrics.ResponseTime * wResponse) +
             (metrics.UptimePercentage * wUptime) +
             (metrics.DataQuality * wQuality)

    // Normalize between 0 and 1
    return math.Max(0, math.Min(1, score))
}

// UpdateReputationBasedOnPerformance updates reputation based on node performance
func UpdateReputationBasedOnPerformance(currentRep float64, success bool, responseTime time.Duration) float64 {
    const (
        maxResponseTime = time.Second * 5
        successWeight   = 0.3
        timeWeight     = 0.7
    )

    // Calculate success impact
    successImpact := 0.0
    if success {
        successImpact = 1.0
    }

    // Calculate response time impact
    timeImpact := 1.0 - math.Min(1.0, float64(responseTime)/float64(maxResponseTime))

    // Calculate new reputation
    newRep := currentRep + (((successImpact * successWeight) + (timeImpact * timeWeight) - currentRep) * 0.1)

    // Ensure reputation stays between 0 and 1
    return math.Max(0, math.Min(1, newRep))
}