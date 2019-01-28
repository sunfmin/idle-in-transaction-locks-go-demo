package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/theplant/gormutils"
)

func main() {
	db, err := gorm.Open("postgres", "user=ec password=123 dbname=ec_test sslmode=disable host=localhost port=5001")
	if err != nil {
		panic(fmt.Sprintf("can not connect to database: %v", err))
	}

	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(70)
	db.Debug()

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Printf("Stats: %#v\n", db.DB().Stats())
		}
	}()

	var runtimeCount = 1000
	var w sync.WaitGroup
	w.Add(runtimeCount)
	for i := 0; i < runtimeCount; i++ {
		go func() {
			defer w.Done()

			deadLockByRequireConnection(db)
		}()
	}
	w.Wait()
}

func deadLockByRequireConnection(db *gorm.DB) {
	var err error
	err = gormutils.Transact(db, func(tx *gorm.DB) (e error) {
		e = tx.Exec(`SELECT pg_sleep(2)`).Error
		if e != nil {
			return e
		}

		// NOTE: this db is not tx, So it will need to execute in different transaction
		return db.Exec(`SELECT pg_sleep(1)`).Error
	})

	if err != nil {
		fmt.Println("err:", err)
	}
}
