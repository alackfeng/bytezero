package flashsign

import (
	"fmt"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// FlashSignApp - flashsign data analysis.
type FlashSignApp struct {
	driverName       string
	dbSourceUrl      string
	dbNameBaasSeal   string
	dbNameBaasReport string
	currentDate      string
	db               *sql.DB

	revenue  *RevenueDay
	business *BusinessDay
}

func NewFlashSignApp() *FlashSignApp {
	currentDate := "2021-11-19 00:00:00"
	return &FlashSignApp{
		driverName:       "mysql",
		dbSourceUrl:      "root:123456@tcp(192.168.90.146:3306)/",
		dbNameBaasSeal:   "baas_seal",
		dbNameBaasReport: "baas_report",
		currentDate:      currentDate,
		revenue:          &RevenueDay{},
		business:         &BusinessDay{currentDate: currentDate},
	}
}

// DBConnect -
func DBConnect(driverName string, dbSourceUrl string, conns int) (*sql.DB, error) {
	db, err := sql.Open(driverName, dbSourceUrl)
	if err != nil {
		fmt.Printf("DBConnect - Open <%s:%s> Error:%v.\n", driverName, dbSourceUrl, err.Error())
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(conns)
	db.SetMaxIdleConns(conns)
	err = db.Ping()
	if err != nil {
		fmt.Printf("DBConnect - Ping <%s:%s> Error:%v.\n", driverName, dbSourceUrl, err.Error())
		return nil, err
	}
	return db, err
}

// dbSourceUrlBaasSeal -
func (f *FlashSignApp) dbSourceUrlBaasSeal() string {
	return f.dbSourceUrl + f.dbNameBaasSeal
}

func (f *FlashSignApp) init() error {
	fmt.Println("FlashSignApp.init -")
	db, err := DBConnect(f.driverName, f.dbSourceUrlBaasSeal(), 10)
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
	f.Business()
	fmt.Println("FlashSignApp.Main - over")
}
