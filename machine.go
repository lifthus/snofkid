package snofkid

import (
	"fmt"
	"sync"
	"time"
)

// NewMachine returns a new SnowflakeMachine based on the given epoch and machine ID.
// It returns an error if the given epoch or machine ID is invalid.
func NewMachine(epoch int64, machineID int64) (*SnowflakeMachine, error) {
	if !IsEpochValid(epoch) || !IsMachineIDValid(machineID) {
		return nil, fmt.Errorf("generating new machine failed by invalid epoch or machine ID")
	}
	return &SnowflakeMachine{
		epoch:     epoch,
		machineID: machineID,
	}, nil
}

// SnowflakeMachine represents a specific machine that generates SnowflakeID.
type SnowflakeMachine struct {
	epoch     int64
	machineID int64

	mu     sync.Mutex
	curMs  int64
	curSeq int64
}

// New generates a new SnowflakeID based on the machine's epoch and ID.
// it returns an error if more than 4096 Snowflakes are generated in the same millisecond.
func (m *SnowflakeMachine) New() (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	curMs := time.Now().UnixMilli()
	timestamp := curMs - m.epoch
	if curMs > m.curMs {
		m.curMs = curMs
		m.curSeq = 0
	}
	if m.curSeq > MaxSequence {
		return 0, fmt.Errorf("snowflakes in a millisecond got exhausted")
	}
	curSeq := m.curSeq
	m.curSeq++
	return from(timestamp, m.machineID, curSeq), nil
}

// Validate validates the sign bit and the machine ID of the given SnowflakeID.
func (m *SnowflakeMachine) Validate(sfid int64) bool {
	if sfid < 0 {
		return false
	}
	if m.machineID != MachineID(sfid) {
		return false
	}
	return true
}

// Time returns the timestamp of the given Snowflake with reflection of the machine's epoch, but without validation.
// the phrase "without validation" means that it doesn't check the sign bit and the machine ID.
func (m *SnowflakeMachine) ParseTime(sfid int64) time.Time {
	ts := RawTimestamp(sfid) + m.epoch
	return time.UnixMilli(ts)
}

func (m SnowflakeMachine) Epoch() int64 {
	return m.epoch
}

func (m SnowflakeMachine) MachineID() int64 {
	return m.machineID
}
