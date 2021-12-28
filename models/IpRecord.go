package models

import (
	"database/sql"
	"log"
	"math"
)

type IpRecord struct {
	ip       string
	hitCount uint8
}

const maxHitCount = uint8(3)

func (ipRecord *IpRecord) Create(db *sql.DB) (sql.Result, error) {

	tx, errBegin := db.Begin()
	if errBegin != nil {
		return nil, errBegin
	}
	defer tx.Rollback()
	stmt, errPrepare := tx.Prepare("INSERT INTO hits(ip, hit_count) VALUES($1,$2)")
	if errPrepare != nil {
		return nil, errPrepare
	}
	defer stmt.Close()
	result, errExec := stmt.Exec(ipRecord.ip, ipRecord.hitCount)
	if errExec != nil {
		log.Printf("Failed to create IpRecord")
		return nil, errExec
	}
	errCommit := tx.Commit()
	if errCommit != nil {
		return nil, errCommit
	}
	return result, nil
}

func (ipRecord *IpRecord) Read(db *sql.DB) (*IpRecord, error) {
	foundRecord := IpRecord{}
	row := db.QueryRow("SELECT * FROM hits WHERE ip = $1", ipRecord.ip)
	err := row.Scan(&foundRecord.ip, &foundRecord.hitCount)
	if err != nil {
		return nil, err
	}
	return &foundRecord, nil
}

func (ipRecord *IpRecord) Update(db *sql.DB) (sql.Result, error) {
	tx, errBegin := db.Begin()
	if errBegin != nil {
		return nil, errBegin
	}
	defer tx.Rollback()
	stmt, errPrepare := tx.Prepare("UPDATE hits SET hit_count=$1 WHERE ip =$2")
	if errPrepare != nil {
		return nil, errPrepare
	}
	defer stmt.Close()
	result, errExec := stmt.Exec(ipRecord.hitCount, ipRecord.ip)
	if errExec != nil {
		log.Printf("Failed to update IpRecord")
		return nil, errExec
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return nil, errCommit
	}
	return result, nil
}

func (ipRecord *IpRecord) Delete(db *sql.DB) (*IpRecord, error) {

	foundIpRecord, err := ipRecord.Read(db)
	if err != nil {
		return nil, err
	}
	tx, errBegin := db.Begin()
	if errBegin != nil {
		return nil, errBegin
	}
	defer tx.Rollback()
	stmt, errPrepare := tx.Prepare("DELETE FROM hits WHERE ip =$1")
	if errPrepare != nil {
		return nil, errPrepare
	}
	defer stmt.Close()
	_, errExec := stmt.Exec(ipRecord.ip)
	if errExec != nil {
		log.Printf("Failed to update IpRecord")
		return nil, errExec
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return nil, errCommit
	}
	return foundIpRecord, nil
}

func NewIpRecord(ip string) *IpRecord {
	return &IpRecord{ip, 1}
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
