package snofkid

import (
	"testing"
)

var (
	testMachine *SnowflakeMachine

	snowflakes4096   []int64
	actualMillis4096 []int64

	snowflakes500ms []int64
)

func TestMain(m *testing.M) {
	var err error
	testMachine, err = NewMachine(TestEpoch, TestMachineID)
	if err != nil {
		panic(err)
	}

	snowflakes4096, actualMillis4096 = generateSnowflakesWithMilliSecsGenerated(testMachine)
	snowflakes500ms = generateSnowflakesFor500ms(testMachine)
	m.Run()
}

func TestSnowflakeFrom(t *testing.T) {
	sfid := from(123456789, 123, 10)
	switch {
	case RawTimestamp(sfid) != 123456789:
		fallthrough
	case MachineID(sfid) != 123:
		fallthrough
	case Sequence(sfid) != 10:
		t.Errorf("SnowflakeIDFrom generated invalid Snowflake ID")
	}
}

func TestRawTimestampGetter(t *testing.T) {
	target := from(123456789, 1, 10)
	if res := RawTimestamp(target); res != 123456789 {
		t.Errorf("Timestamp is expected to be %d but got %d", 123456789, res)
	}
}

func TestMachineIDGetter(t *testing.T) {
	target := from(123456789, 123, 10)
	if res := MachineID(target); res != 123 {
		t.Errorf("MachineID is expected to be %d but got %d", 123, res)
	}
}

func TestSequenceGetter(t *testing.T) {
	target := from(123456789, 1, 123)
	if res := Sequence(target); res != 123 {
		t.Errorf("Sequence is expected to be %d but got %d", 123, res)
	}
}
func TestMaxTimestampConst(t *testing.T) {
	calcedMaxTs := int64(-1) ^ (int64(-1) << TimestampBits)
	if MaxTimestamp != calcedMaxTs {
		t.Errorf("MaxTimestamp is expected to be %d but got %d", calcedMaxTs, MaxTimestamp)
	}
}

func TestMaxMachineIDConst(t *testing.T) {
	calcedMaxMID := -1 ^ (-1 << MachineIDBits)
	if MaxMachineID != calcedMaxMID {
		t.Errorf("MaxMachineID is expected to be %d but got %d", calcedMaxMID, MaxMachineID)
	}
}

func TestMaxSequenceConst(t *testing.T) {
	calcedMaxSeq := -1 ^ (-1 << SequenceBits)
	if MaxSequence != calcedMaxSeq {
		t.Errorf("MaxSequence is expected to be %d but got %d", calcedMaxSeq, MaxSequence)
	}
}
