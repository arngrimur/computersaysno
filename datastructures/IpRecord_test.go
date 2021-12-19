package datastructures

import (
	"github.com/stretchr/testify/assert"
	"testing"
)
const ip = "1.1.1.1"
func TestIpRecord_New(t *testing.T) {
	record := NewIpRecord(ip)
	assert.Equal(t, ip, record.ip)
	assert.Equal(t, uint8(1), record.hitCount)
}

func TestIpRecord_CountIncrease(t *testing.T) {
	record := NewIpRecord(ip)
	record.IncreaseHitCount()
	assert.Equal(t, uint8(2), record.hitCount)
}

func TestIpRecord_CountIncreaseStops(t *testing.T) {
	record := NewIpRecord(ip)
	record.hitCount = 0xff
	record.IncreaseHitCount()
	assert.Equal(t, uint8(255), record.hitCount)
}

func TestIpRecord_CountIsNotBlocked(t *testing.T) {
	record := NewIpRecord(ip)
	record.hitCount = uint8(3)
	assert.Equal(t, false, record.IsBlocked(ip))
}

func TestIpRecord_CountIsBlocked(t *testing.T) {
	record := NewIpRecord(ip)
	record.hitCount = uint8(4)
	assert.Equal(t, true, record.IsBlocked(ip))
}

func TestIpRecord_GetIp(t *testing.T) {
	record := NewIpRecord(ip)
	assert.Equal(t, ip,record.GetIp())
}

func TestIpRecord_GetHitCount(t *testing.T) {
	record := NewIpRecord(ip)
	record.IncreaseHitCount()
	assert.Equal(t, uint8(2), record.GetHitCount())
}

