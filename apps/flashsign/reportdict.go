package flashsign

import (
	"database/sql"
	"fmt"
)

// BaasReportDict - 报告字典表.
type BaasReportDict struct {
	lastReportDate   string
	dbSourceBaasSeal string // 报表数据来源数据库.
}

// NewBaasReportDict -
func NewBaasReportDict() *BaasReportDict {
	return &BaasReportDict{lastReportDate: "2020-10-19 00:00:00"}
}

// String -
func (b BaasReportDict) String() string {
	return fmt.Sprintf("BaasReportDict[lastReportDate:%s,dbSourceBaasSeal:%s]", b.lastReportDate, b.dbSourceBaasSeal)
}

// valid -
func (b *BaasReportDict) valid() error {
	if b.dbSourceBaasSeal == "" {
		return fmt.Errorf("t_report_dict.%s is null", b.dbSourceBaasSeal)
	} else if b.lastReportDate == "" {
		return fmt.Errorf("t_report_dict.%s is null", b.lastReportDate)
	}
	return nil
}

// load -
func (b *BaasReportDict) load(db *sql.DB) error {
	sqlQuery := "select item, value from t_report_dict; "
	rows, err := db.Query(sqlQuery)
	if err != nil {
		fmt.Println("BaasReportDict.load - error.", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var item, value string
		if err := rows.Scan(&item, &value); err != nil {
			return err
		}
		if item == "lastReportDate" {
			b.lastReportDate = value
		} else if item == "dbSourceBaasSeal" {
			b.dbSourceBaasSeal = value
		}
	}
	return b.valid()
}

// get -
func (b *BaasReportDict) get(db *sql.DB, key string) (value string) {
	sqlQuery := "SELECT value from t_report_dict where item = ?; "
	err := db.QueryRow(sqlQuery, key).Scan(&value)
	if err != nil {
		fmt.Println("FlashSignApp.RevenueDayExpendCount - error.", err.Error())
		return ""
	}
	return value
}

// updateLastReportDate -
func (b *BaasReportDict) updateLastReportDate(db *sql.DB) error {
	// now := time.Now()
	sqlQuery := "update t_report_dict set value = ? where item = 'lastReportDate'; "
	_, err := db.Exec(sqlQuery, b.lastReportDate)
	if err != nil {
		return err
	}
	// count, _ := res.RowsAffected()
	// fmt.Println("BaasReportDict.updateLastReportDate - update t_report_dict.lastReportDate value ", b.lastReportDate, now, count)
	return nil
}
