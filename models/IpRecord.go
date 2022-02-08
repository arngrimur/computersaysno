package models

import (
	"database/sql"
	"log"
	"math"
	"reflect"
)

type IpRecord struct {
	ip       string
	hitCount uint8
}

const MaxHitCount = uint8(3)

func (ipRecord *IpRecord) Create(db *sql.DB) (sql.Result, error) {

	tx, errBegin := db.Begin()
	if errBegin != nil {
		return nil, errBegin
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if errBegin != nil {
			log.Printf("Failed to to rollback when creating new %s, %s", reflect.TypeOf(ipRecord), err)
		}
	}(tx)
	stmt, errPrepare := tx.Prepare("INSERT INTO hits(ip, hit_count) VALUES($1,$2)")
	if errPrepare != nil {
		return nil, errPrepare
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("Failed preparing statement, %s.", err)
		}
	}(stmt)
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
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("Failed to to rollback when creating new %s, %s", reflect.TypeOf(ipRecord), rollbackErr)
			}
		}
	}()

	stmt, errPrepare := tx.Prepare("UPDATE hits SET hit_count=$1 WHERE ip =$2")
	if errPrepare != nil {
		return nil, errPrepare
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("Failed preparing statement, %s.", err)
		}
	}(stmt)
	result, err := stmt.Exec(ipRecord.hitCount, ipRecord.ip)
	if err != nil {
		log.Printf("Failed to update IpRecord")
		return nil, err
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
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("Failed to to rollback when creating new %s, %s", reflect.TypeOf(ipRecord), rollbackErr)
			}
		}
	}()
	stmt, errPrepare := tx.Prepare("DELETE FROM hits WHERE ip =$1")
	if errPrepare != nil {
		return nil, errPrepare
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("Failed to close statemet %s, %s", reflect.TypeOf(ipRecord), err)
		}
	}(stmt)
	_, err = stmt.Exec(ipRecord.ip)
	if err != nil {
		log.Printf("Failed to update IpRecord")
		return nil, err
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
	return ipRecord.ip == ip && ipRecord.hitCount > MaxHitCount
}

func (ipRecord *IpRecord) GetIp() string {
	return ipRecord.ip
}

func (ipRecord *IpRecord) GetHitCount() uint8 {
	return ipRecord.hitCount
}

func (ipRecord *IpRecord) SetMaxHitCount() {
	ipRecord.hitCount = MaxHitCount
}
