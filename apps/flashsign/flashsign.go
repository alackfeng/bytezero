package flashsign

import (
	"fmt"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// FlashSignApp - flashsign data analysis.
type FlashSignApp struct {
	driverName  string
	dbSourceUrl string
	currentDate string
	db          *sql.DB

	revenue *RevenueDay
}

func NewFlashSignApp() *FlashSignApp {
	return &FlashSignApp{
		driverName:  "mysql",
		dbSourceUrl: "root:123456@tcp(192.168.90.146:3306)/baas_seal",
		currentDate: "2021-11-19 00:00:00",
		revenue:     &RevenueDay{},
	}
}

func (f *FlashSignApp) init() error {
	fmt.Println("FlashSignApp.init -")
	db, err := sql.Open(f.driverName, f.dbSourceUrl)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		return err
	}
	f.db = db
	return nil
}

// Main -
func (f *FlashSignApp) Main() {
	if err := f.init(); err != nil {
		fmt.Printf("FlashSignApp.Main - error %v\n", err.Error())
		return
	}
	f.Revenue()
	fmt.Println("FlashSignApp.Main - over")
}
