package datastructures

import "math"

type IpRecord struct {
	ip       string
	hitCount uint8
}

const maxHitCount = uint8(3)

func NewIpRecord(ip string) *IpRecord {
	return &IpRecord{ip,1}
}

func (ipRecord *IpRecord) IncreaseHitCount() {
	if ipRecord.hitCount < math.MaxUint8 {
		ipRecord.hitCount++
	}
}

func (ipRecord *IpRecord) IsBlocked(ip string) bool {
	return ipRecord.ip == ip && ipRecord.hitCount > maxHitCount
}

func (ipRecord *IpRecord) GetIp() string {
	return ipRecord.ip
}

func (ipRecord *IpRecord) GetHitCount() uint8 {
	return ipRecord.hitCount
}
