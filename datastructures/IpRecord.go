package datastructures

import (
	"db"
	"math"
)

type IpRecord struct {
	ip       string
	hitCount uint8
}

func (ipRecord *IpRecord) Create() {
	db, err := db.Init()
	if err != nil {
		return
	}
	defer db.Close()
	ctx, _ := db.Begin()
	stmt, _ := ctx.Prepare("INSERT INTO ip_record(ip, hit_count) VALUES(?,?)")
	stmt.Close()
	stmt.Exec(ipRecord.ip,ipRecord.hitCount)
	ctx.Commit()
}

func (ipRecord *IpRecord) Read() {
	panic("implement me")
}

func (ipRecord *IpRecord) Update() {
	panic("implement me")
}

func (ipRecord *IpRecord) Delete() {
	panic("implement me")
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
