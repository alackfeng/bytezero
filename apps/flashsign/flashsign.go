package flashsign

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/alackfeng/bytezero/cores/utils"
	_ "github.com/go-sql-driver/mysql"
)

// FlashSignApp - flashsign data analysis.
type FlashSignApp struct {
	driverName         string
	dbSourceBaasReport string
	config             *BaasReportDict
	reportDb           *sql.DB
	db                 *sql.DB
}

func NewFlashSignApp() *FlashSignApp {
	return &FlashSignApp{
		driverName:         "mysql",
		dbSourceBaasReport: "root:123456@tcp(192.168.90.53:3306)/baas_report",
		config: &BaasReportDict{
			lastReportDate: "2020-10-19 00:00:00",
		},
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
	return f.config.dbSourceBaasSeal
}

// dbSourceUrlBaasReport -
func (f *FlashSignApp) dbSourceUrlBaasReport() string {
	return f.dbSourceBaasReport
}

// lastReportDate -
func (f *FlashSignApp) lastReportDate() string {
	return f.config.lastReportDate
}

// init -
func (f *FlashSignApp) init() error {
	reportDb, err := DBConnect(f.driverName, f.dbSourceUrlBaasReport(), 3)
	if err != nil {
		return err
	}
	f.reportDb = reportDb
	if err := f.config.load(f.reportDb); err != nil {
		return err
	}
	fmt.Println("FlashSignApp.init - config: ", f.config)

	db, err := DBConnect(f.driverName, f.dbSourceUrlBaasSeal(), 10)
	if err != nil {
		return err
	}
	f.db = db

	return nil
}

// close -
func (f *FlashSignApp) close() {
	fmt.Println("FlashSignApp.close -")
	if f.db != nil {
		f.db.Close()
	}
	if f.reportDb != nil {
		f.reportDb.Close()
	}
}

// FormatNextDate -
func FormatNextDate(d string, day int) (string, bool) {
	now, _ := time.Parse("2006-01-02 00:00:00", d)
	currMs := now.AddDate(0, 0, day)
	if time.Now().Before(currMs) {
		return "", false
	}
	currMs = now.AddDate(0, 0, 1)
	return fmt.Sprintf("%04d-%02d-%02d 00:00:00", currMs.Year(), currMs.Month(), currMs.Day()), true
}

// CheckReportDate -
func CheckReportDate(d string) (string, error) {
	currMs, err := time.Parse("2006-01-02 00:00:00", d)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%04d-%02d-%02d 00:00:00", currMs.Year(), currMs.Month(), currMs.Day()), nil
}

// Main -
// go run .\main.go flashsign .
// go run .\main.go flashsign -d '2022-10-19 00:00:00' .
func (f *FlashSignApp) Main(reportDate string) {
	if err := f.init(); err != nil {
		fmt.Printf("FlashSignApp.Main - error %v\n", err.Error())
		return
	}

	count := 0
	bQuit := false
	if reportDate != "" {
		d, err := CheckReportDate(reportDate)
		if err != nil {
			fmt.Printf("FlashSignApp.Main - report date param <%s> error:%v\n.", reportDate, err.Error())
		} else {
			f.config.lastReportDate = d
			dura := utils.NewDuration()
			f.Revenue(f.config.lastReportDate)
			f.Business(f.config.lastReportDate)
			fmt.Printf("FlashSignApp.Main - execute task: %s, dura:[%v] %dms\n", f.config.lastReportDate, dura.Begin(), dura.DuraMs())
		}

	} else {
		for {
			if bQuit {
				break
			}
			_, ok := FormatNextDate(f.config.lastReportDate, 1)
			if !ok {
				fmt.Println("FlashSignApp.Main - current time ", time.Now(), " overload ", f.config.lastReportDate)
				break
			}
			dura := utils.NewDuration()
			f.Revenue(f.config.lastReportDate)
			f.Business(f.config.lastReportDate)
			fmt.Printf("FlashSignApp.Main - execute task: %s, dura:[%v] %dms\n", f.config.lastReportDate, dura.Begin(), dura.DuraMs())
			d, ok := FormatNextDate(f.config.lastReportDate, 1)
			if !ok {
				fmt.Println("FlashSignApp.Main - current time ", time.Now(), " overload ", f.config.lastReportDate)
				break
			}
			f.config.lastReportDate = d
			f.config.updateLastReportDate(f.reportDb)
			if count%100 == 0 {
				time.Sleep(time.Millisecond * 50)
			}
			count++
		}
	}

	f.close()
	fmt.Println("FlashSignApp.Main - over")
}
