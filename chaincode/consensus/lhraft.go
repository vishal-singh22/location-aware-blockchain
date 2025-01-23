package consensus

import (
    "time"
    "math"
    "sync"
)

// ConsensusNode represents a node in the LH-Raft consensus
type ConsensusNode struct {
    ID            string
    Location      string
    Reputation    float64
    IsLeader      bool
    GroupMembers  []string
    LastHeartbeat time.Time
    State         NodeState
    mu            sync.Mutex
}

// NodeState represents the state of a node in the consensus
type NodeState int

const (
    Follower NodeState = iota
    Candidate
    Leader
)

// LHRaftConsensus is the main consensus structure
type LHRaftConsensus struct {
    Nodes         map[string]*ConsensusNode
    Threshold     float64
    ZoneLeaders   map[string]string  // ZoneID -> LeaderID
    mu            sync.RWMutex
}

// NewLHRaftConsensus creates a new instance of the consensus
func NewLHRaftConsensus(threshold float64) *LHRaftConsensus {
    return &LHRaftConsensus{
        Nodes:       make(map[string]*ConsensusNode),
        Threshold:   threshold,
        ZoneLeaders: make(map[string]string),
    }
}

// RegisterNode adds a new node to the consensus
func (l *LHRaftConsensus) RegisterNode(id, location string, reputation float64) error {
    l.mu.Lock()
    defer l.mu.Unlock()

    node := &ConsensusNode{
        ID:           id,
        Location:     location,
        Reputation:   reputation,
        IsLeader:     false,
        GroupMembers: make([]string, 0),
        State:       Follower,
    }

    l.Nodes[id] = node
    return nil
}

// FormCandidateGroups creates location-based consensus groups
func (l *LHRaftConsensus) FormCandidateGroups(location string) []string {
    l.mu.RLock()
    defer l.mu.RUnlock()

    candidates := make([]string, 0)
    for id, node := range l.Nodes {
        if node.Location == location && node.Reputation >= l.Threshold {
            candidates = append(candidates, id)
        }
    }
    return candidates
}

// ElectZoneLeader implements the leader election for a specific zone
func (l *LHRaftConsensus) ElectZoneLeader(zoneID string) (string, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    var bestCandidate string
    var highestRep float64 = -1

    // Find the node with highest reputation in the zone
    for id, node := range l.Nodes {
        if node.Location == zoneID && node.Reputation > highestRep {
            highestRep = node.Reputation
            bestCandidate = id
        }
    }

    if bestCandidate != "" {
        l.ZoneLeaders[zoneID] = bestCandidate
        l.Nodes[bestCandidate].IsLeader = true
        l.Nodes[bestCandidate].State = Leader
    }

    return bestCandidate, nil
}

// UpdateNodeReputation updates a node's reputation and triggers re-election if needed
func (l *LHRaftConsensus) UpdateNodeReputation(nodeID string, newReputation float64) error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if node, exists := l.Nodes[nodeID]; exists {
        oldRep := node.Reputation
        node.Reputation = newReputation

        // If this is a leader and reputation dropped below threshold
        if node.IsLeader && newReputation < l.Threshold {
            node.IsLeader = false
            node.State = Follower
            // Trigger re-election for the zone
            go l.ElectZoneLeader(node.Location)
        }
        return nil
    }
    return fmt.Errorf("node not found: %s", nodeID)
}

// PropagateTransaction handles transaction propagation in the hierarchy
func (l *LHRaftConsensus) PropagateTransaction(transaction []byte, zoneID string) error {
    l.mu.RLock()
    leaderID, exists := l.ZoneLeaders[zoneID]
    l.mu.RUnlock()

    if !exists {
        return fmt.Errorf("no leader found for zone: %s", zoneID)
    }

    // Simulate local consensus
    success := l.achieveLocalConsensus(leaderID, transaction)
    if !success {
        return fmt.Errorf("failed to achieve local consensus in zone: %s", zoneID)
    }

    // Simulate global consensus propagation
    return l.propagateToGlobalConsensus(transaction)
}

func (l *LHRaftConsensus) achieveLocalConsensus(leaderID string, transaction []byte) bool {
    l.mu.RLock()
    leader := l.Nodes[leaderID]
    l.mu.RUnlock()

    if leader == nil {
        return false
    }

    // Simulate consensus achievement (In real implementation, this would involve
    // actual communication between nodes)
    return true
}

func (l *LHRaftConsensus) propagateToGlobalConsensus(transaction []byte) error {
    // In real implementation, this would coordinate with other zone leaders
    return nil
}