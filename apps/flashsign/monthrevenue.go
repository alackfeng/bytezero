package flashsign

import (
	"database/sql"
	"fmt"
	"time"
)


// RevenueMonth - 月度营收维度分析.
type RevenueMonth struct {
	currentMonth         string  // 月份.
	totalAmount          float64 // 月营收.
	growthRate           float64 // 增长率.
	forecastGrowth       float64 // 预测增长率.
	forecastAmount       float64 // 预测收益.

	lastMonthAmount           float64 // 上月收入.
}

// NewRevenueMonth -
func NewRevenueMonth() *RevenueMonth {
	return &RevenueMonth{}
}

// String -
func (b RevenueMonth) String() string {
	return fmt.Sprintf("%s: [totalAmount:%f, growthRate:%f, forecastGrowth:%f, forecastAmount:%f]",
		b.currentMonth, b.totalAmount, b.growthRate, b.forecastGrowth, b.forecastAmount)
}

// CurrentMonth - 当月.
func CurrentMonth(d string) string {
	now, _ := time.Parse("2006-01-02 00:00:00", d)
        return fmt.Sprintf("%04d-%02d", now.Year(), now.Month())
}

// LastMonth - 上个月.
func LastMonth(d string, month int) (string, string) {
	now, _ := time.Parse("2006-01-02 00:00:00", d)
        lastMonth := now.AddDate(0, month, 0)
        return fmt.Sprintf("%04d-%02d", now.Year(), now.Month()), fmt.Sprintf("%04d-%02d", lastMonth.Year(), lastMonth.Month())
}

// LastThreeMonth - 前三个月.
func LastThreeMonth(d string, month int) (string, string, string) {
	now, _ := time.Parse("2006-01-02 00:00:00", d)
        lastMonth := now.AddDate(0, month, 0)
        lastMonth2 := now.AddDate(0, month-1, 0)
        lastMonth3 := now.AddDate(0, month-2, 0)
        return fmt.Sprintf("%04d-%02d", lastMonth.Year(), lastMonth.Month()), fmt.Sprintf("%04d-%02d", lastMonth2.Year(), lastMonth2.Month()), fmt.Sprintf("%04d-%02d", lastMonth3.Year(), lastMonth3.Month())
}

// RevenueMonthTotalAmount - 月营收.
func (f *RevenueMonth) RevenueMonthTotalAmount(db *sql.DB) error {
	currentMonth1 := CurrentMonth(f.currentMonth)
	sqlQuery := "select IFNULL(sum(totalAmount), 0) from t_report_revenue where currentDate like ?; "
	err := db.QueryRow(sqlQuery, currentMonth1+"%").Scan(&f.totalAmount)
	if err != nil {
		fmt.Println("RevenueMonth.RevenueMonthTotalAmount - error.", err.Error())
		return err
	}
	fmt.Println("RevenueMonth.RevenueMonthTotalAmount - totalAmount.", f.totalAmount, currentMonth1)
	return nil
}

// RevenueMonthGrowthRate - 增长率.
func (f *RevenueMonth) RevenueMonthGrowthRate(db *sql.DB) error {
        _, lastMonth := LastMonth(f.currentMonth, -1)
        sqlQuery := "select totalAmount from t_report_month_revenue where currentMonth = ?; "
	err := db.QueryRow(sqlQuery, lastMonth).Scan(&f.lastMonthAmount)
        if err != nil {
                fmt.Println("RevenueMonth.RevenueMonthGrowthRate - error.", err.Error())
                return err
        }
        fmt.Println("RevenueMonth.RevenueMonthGrowthRate - growthRate ", f.growthRate, f.lastMonthAmount)
	if f.lastMonthAmount != 0 {
		f.growthRate = (f.totalAmount - f.lastMonthAmount) / f.lastMonthAmount
	}
        return nil
}

// RevenueMonthForecastGrowth - 预测增长率.
func (f *RevenueMonth) RevenueMonthForecastGrowth(db *sql.DB) error {
	lm1, lm2, lm3 := LastThreeMonth(f.currentMonth, -1)
	sqlQuery := "SELECT IFNULL(SUM(growthRate),0)/3 from t_report_month_revenue where currentMonth in (?,?,?); "
	err := db.QueryRow(sqlQuery, lm1, lm2, lm3).Scan(&f.forecastGrowth)
	if err != nil {
		fmt.Println("RevenueMonth.RevenueMonthForecastGrowth - error.", err.Error())
		return err
	}
        fmt.Println("RevenueMonth.RevenueMonthForecastGrowth - forecastGrowth ", f.forecastGrowth)
	return nil
}

// RevenueMonthForecastAmount - 预测收益.
func (f *RevenueMonth) RevenueMonthForecastAmount(db *sql.DB) error {
	f.forecastAmount = (1+f.forecastGrowth) * f.lastMonthAmount
        fmt.Println("RevenueMonth.RevenueMonthForecastAmount - forecastAmount ", f.forecastAmount)
	return nil
}

// RevenueRemove -
func (f *RevenueMonth) RevenueMonthRemove(reportDb *sql.DB) error {
	currentMonth, _ := LastMonth(f.currentMonth, -1)
	sqlQuery := "delete from t_report_month_revenue where currentMonth = ?; "
	_, err := reportDb.Exec(sqlQuery, currentMonth)
	if err != nil {
		fmt.Println("RevenueMonth.RevenueRemove - error", err.Error())
		return err
	}
	// count, _ := res.RowsAffected()
	// fmt.Println("RevenueMonth.RevenueRemove - delete rows: ", count)
	return nil
}

// RevenueInsert -
func (f *RevenueMonth) RevenueMonthInsert(reportDb *sql.DB) error {
	sqlQuery := "insert into t_report_month_revenue(currentMonth, totalAmount, growthRate, forecastGrowth, forecastAmount, createTime) " +
		"values (?,?,?,?,?,?)"
	currentMonth, _ := LastMonth(f.currentMonth, -1)
	res, err := reportDb.Exec(sqlQuery, currentMonth, f.totalAmount, f.growthRate, f.forecastGrowth, f.forecastAmount, time.Now() )
	if err != nil {
		fmt.Println("RevenueMonth.RevenueInsert - error", err.Error())
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		fmt.Println("RevenueMonth.RevenueInsert - LastInsertId error", err.Error())
		return err
	}
	// fmt.Println("RevenueMonth.RevenueInsert - insert id ", id)
	return nil
}

// Revenue - 营收维度分析.
func (f *FlashSignApp) RevenueMonth(lastDate string) {
	o := &RevenueMonth{currentMonth: lastDate}
	o.RevenueMonthTotalAmount(f.reportDb)
	o.RevenueMonthGrowthRate(f.reportDb)
	o.RevenueMonthForecastGrowth(f.reportDb)
	o.RevenueMonthForecastAmount(f.reportDb)
	fmt.Println("FlashSignApp.RevenueMonth - ", o)
	o.RevenueMonthRemove(f.reportDb)
	o.RevenueMonthInsert(f.reportDb)
}
