package flashsign

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alackfeng/bytezero/cores/utils"
)

// PaymodeWechat -
const (
	PaymodeWechat = 0
	PaymodeAliPay = 1
)

// PurchasePackageAmount -
const (
	PurchasePackageAmount1   = 1
	PurchasePackageAmount5   = 5
	PurchasePackageAmount10  = 10
	PurchasePackageAmount50  = 50
	PurchasePackageAmount100 = 100
)

const (
	ActivityTypePurchase = 0
	ActivityTypePresent  = 1
)

// RevenueDay - 营收维度分析.
type RevenueDay struct {
	currentDate          string  // 日期.
	totalAmount          float64 // 当日总收入.
	stockAll             int     // 当前所有份数.
	stockCount           int     // 当日库存份数.
	stockPurchaseCount   int     // 当日库存已购份数.
	stockPresentCount    int     // 当日库存赠送份数.
	expiredPurchaseCount int     // 当日过期份数(已购).
	expiredPresentCount  int     // 当日过期份数(赠送).
	expendCount          int     // 当日消耗份数.

	wechatTransAmount         float64 // 微信当日交易金额.
	wechatTransAccess         int     // 微信当日交易用户.
	wechatTransCount          int     // 微信当日交易笔数.
	wechatRepurchaseAccess    int     // 微信日复购用户.
	wechatFirstPurchaseAccess int     // 微信日新购用户.

	alipayTransAmount         float64 // 支付宝当日交易金额.
	alipayTransAccess         int     // 支付宝当日交易用户.
	alipayTransCount          int     // 支付宝当日交易笔数.
	alipayRepurchaseAccess    int     // 支付宝日复购用户.
	alipayFirstPurchaseAccess int     // 支付宝日新购用户.

	presentCount             int // 赠送份数.
	purchasePackageAmount1   int // 购买单份套餐次数.
	purchasePackageAmount5   int // 购买5份套餐次数.
	purchasePackageAmount10  int // 购买10份套餐次数.
	purchasePackageAmount50  int // 购买50份套餐次数.
	purchasePackageAmount100 int // 购买100份套餐次数.
	averageAmount30day	 float64 // 30日收入均值.
}

// NewRevenueDay -
func NewRevenueDay() *RevenueDay {
	return &RevenueDay{}
}

// String -
func (b RevenueDay) String() string {
	return fmt.Sprintf("%s: [totalAmount:%f, stockCount:%d, expendCount:%d, wechat<amount:%f,access:%d,count:%d,repurchase:%d,first:%d>, alipay<amount:%f,access:%d,count:%d,repurchase:%d,first:%d>, package<present:%d, p1:%d, p5:%d, p10:%d, p50:%d, p100:%d>]",
		b.currentDate, b.totalAmount, b.stockCount, b.expendCount,
		b.wechatTransAmount, b.wechatTransAccess, b.wechatTransCount, b.wechatRepurchaseAccess, b.wechatFirstPurchaseAccess,
		b.alipayTransAmount, b.alipayTransAccess, b.alipayTransCount, b.alipayRepurchaseAccess, b.alipayFirstPurchaseAccess,
		b.presentCount, b.purchasePackageAmount1, b.purchasePackageAmount5, b.purchasePackageAmount10, b.purchasePackageAmount50, b.purchasePackageAmount100)
}

// RevenueDayTotalAmount - 当日总收入.
func (f *RevenueDay) RevenueDayTotalAmount(db *sql.DB) error {
	sqlQuery := "SELECT IFNULL(SUM(price-dis_amount),0) as totalAmount FROM t_order where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and status = 2 "
	err := db.QueryRow(sqlQuery, f.currentDate).Scan(&f.totalAmount)
	if err != nil {
		fmt.Println("RevenueDay.RevenueDayTotalAmount - error.", err.Error())
		return err
	}
	// fmt.Println("RevenueDay.RevenueDayTotalAmount - totalAmount.", f.revenue.totalAmount)
	return nil
}

// RevenueDayStockAll - 当前所有份数 - 当前系统中全部已产生总的合同份数.
func (f *RevenueDay) RevenueDayStockAll(db *sql.DB) error {
        currentDateEnd := utils.FormatNextDateMs(f.currentDate)
        sqlQuery := "SELECT IFNULL(SUM(amount),0) as stockAll from t_bought_package where status != 4 and create_time < ?; "
        err := db.QueryRow(sqlQuery, currentDateEnd).Scan(&f.stockAll)
        if err != nil {
                fmt.Println("RevenueDay.RevenueDayStockAll - error.", err.Error())
                return err
        }
        // fmt.Println("RevenueDay.RevenueDayStockAll - stockAll.", f.revenue.stockAll)
        return nil
}

// RevenueDayStockCount - 当日库存份数 - 当天统计过往待签署数量(购买的+赠送的), 去掉单份体验合同.
func (f *RevenueDay) RevenueDayStockCount(db *sql.DB) error {
	currentDateEnd := utils.FormatNextDateMs(f.currentDate)
	sqlQuery := "SELECT IFNULL(SUM(count),0) as stockCount, activity_type from t_bought_package where status = 0 and expired_time > ? and create_time < ? group by activity_type; "
	rows, err := db.Query(sqlQuery, currentDateEnd, currentDateEnd)
	if err != nil {
		fmt.Println("RevenueDay.RevenueDayStockCount - error.", err.Error())
		return err
	}
	defer rows.Close()
	
	var total int
	for rows.Next() {
		var count int
		var activity int
		if err := rows.Scan(&count, &activity); err != nil {
			return err
		}
		if activity == ActivityTypePurchase {
			f.stockPurchaseCount = count
		} else if activity == ActivityTypePresent {
			f.stockPresentCount = count
		} else {
			total = count
		}
	}
	f.stockCount = f.stockPurchaseCount + f.stockPresentCount + total
	// fmt.Println("RevenueDay.RevenueDayStockCount - stockQuantity.", f.revenue.stockCount)
	return nil
}

// RevenueDayExpiredPurchaseCount - 当日过期份数(已购).
func (f *RevenueDay) RevenueDayExpiredPurchaseCount(db *sql.DB) error {
	sqlQuery := "select IFNULL(SUM(count),0) as expiredPurchaseCount from t_bought_package where status = 2 and activity_type = 0 and FROM_UNIXTIME(expired_time DIV 1000, '%Y-%m-%d 00:00:00') = ?; "
	err := db.QueryRow(sqlQuery, f.currentDate).Scan(&f.expiredPurchaseCount)
	if err != nil {
		fmt.Println("RevenueDay.RevenueDayExpiredPurchaseCount - error.", err.Error())
		return err
	}
	// fmt.Println("RevenueDay.RevenueDayExpiredPurchaseCount - expiredPurchaseCount.", f.revenue.expiredPurchaseCount)
	return nil
}

// RevenueDayExpiredPresentCount - 当日过期份数(赠送).
func (f *RevenueDay) RevenueDayExpiredPresentCount(db *sql.DB) error {
	sqlQuery := "select IFNULL(SUM(count),0) as expiredPresentCount from t_bought_package where status = 2 and activity_type = 1 and FROM_UNIXTIME(expired_time DIV 1000, '%Y-%m-%d 00:00:00') = ?; "
	err := db.QueryRow(sqlQuery, f.currentDate).Scan(&f.expiredPresentCount)
	if err != nil {
		fmt.Println("RevenueDay.RevenueDayExpiredPresentCount - error.", err.Error())
		return err
	}
	// fmt.Println("RevenueDay.RevenueDayExpiredPresentCount - expiredPresentCount.", f.revenue.expiredPresentCount)
	return nil
}

type PackageDeductionKindCount struct {
	count int // 份数.
	kind  int // group by.
}

// RevenueDayExpendCount - 当日消耗份数: 当日套餐抵扣0（划扣） -  撤回3（撤销划扣）.
func (f *RevenueDay) RevenueDayExpendCount(db *sql.DB) error {
	sqlQuery := "SELECT IFNULL(SUM(amount),0) as count, package_deduction_kind as kind  from t_deduction_record WHERE FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? GROUP BY package_deduction_kind; "
	rows, err := db.Query(sqlQuery, f.currentDate)
	if err != nil {
		fmt.Println("FlashSignApp.RevenueDayExpendCount - error.", err.Error())
		return err
	}
	defer rows.Close()
	var packageDeduction map[int]int = make(map[int]int)
	for rows.Next() {
		var res PackageDeductionKindCount
		if err := rows.Scan(&res.count, &res.kind); err != nil {
			return err
		}
		packageDeduction[res.kind] = res.count
	}
	if count, ok := packageDeduction[0]; ok {
		if reduct, ok := packageDeduction[3]; ok {
			f.expendCount = count - reduct
		} else {
			f.expendCount = count
		}
	}
	// fmt.Println("RevenueDay.RevenueDayExpendCount - expendCount.", f.revenue.expendCount)
	return nil
}

// RevenueDayWechatTrans - 微信当日交易金额 微信当日交易用户 微信当日交易笔数.
func (f *RevenueDay) RevenueDayWechatTrans(db *sql.DB, paymode int) error {
	sqlQuery := "SELECT IFNULL(SUM(price-dis_amount),0) as wechatTransAmount, IFNULL(COUNT(DISTINCT access_id),'') as wechatTransAccess, IFNULL(COUNT(1),0) as wechatTransCount FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and `status` = 2 and pay_method=?; "
	var err error
	if paymode == PaymodeAliPay {
		err = db.QueryRow(sqlQuery, f.currentDate, paymode).Scan(&f.alipayTransAmount, &f.alipayTransAccess, &f.alipayTransCount)
	} else {
		err = db.QueryRow(sqlQuery, f.currentDate, paymode).Scan(&f.wechatTransAmount, &f.wechatTransAccess, &f.wechatTransCount)
	}
	if err != nil {
		fmt.Println("RevenueDay.RevenueDayWechatTrans - error.", err.Error())
		return err
	}
	// if paymode == PaymodeAliPay {
	// 	fmt.Println("RevenueDay.RevenueDayWechatTrans - alipayTrans.", f.revenue.alipayTransAmount, f.revenue.alipayTransAccess, f.revenue.alipayTransCount)
	// } else {
	// 	fmt.Println("RevenueDay.RevenueDayWechatTrans - wechatTrans.", f.revenue.wechatTransAmount, f.revenue.wechatTransAccess, f.revenue.wechatTransCount)
	// }
	return nil
}

// RevenueDayWechatRepurchaseAccess - 微信日复购用户.
func (f *RevenueDay) RevenueDayWechatRepurchaseAccess(db *sql.DB, paymode int) error {
	sqlQuery := "SELECT IFNULL(COUNT(1),0) as repurchaseAccess from (SELECT COUNT(access_id) as access_ids FROM t_order where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and status = 2 and pay_method=? GROUP BY access_id) t where t.access_ids>1; "
	var access int
	err := db.QueryRow(sqlQuery, f.currentDate, paymode).Scan(&access)
	if err != nil {
		fmt.Println("RevenueDay.RevenueDayWechatRepurchaseAccess - error.", err.Error())
		return err
	}
	if paymode == PaymodeAliPay {
		f.alipayRepurchaseAccess = access
		// fmt.Println("RevenueDay.RevenueDayWechatRepurchaseAccess - alipayRepurchaseAccess.", f.revenue.alipayRepurchaseAccess)
	} else {
		f.wechatRepurchaseAccess = access
		// fmt.Println("RevenueDay.RevenueDayWechatRepurchaseAccess - repurchaseAccess.", f.revenue.wechatRepurchaseAccess)
	}
	return nil
}

// FirstPurchaseAccess -
type FirstPurchaseAccess struct {
	accessId string
	exist    string
}

// RevenueDayWechatFirstPurchaseAccess - 微信日新购用户 - 微信平台首次购买的记录（程序实现：先建立首次购买表按平台，在查找是否存在购买）.
func (f *RevenueDay) RevenueDayWechatFirstPurchaseAccess(db *sql.DB, paymode int) error {
	sqlQuery := "SELECT DISTINCT access_id as accessId FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and `status` = 2 and pay_method=? GROUP BY access_id;  "
	rows, err := db.Query(sqlQuery, f.currentDate, paymode)
	if err != nil {
		fmt.Println("FlashSignApp.RevenueDayWechatFirstPurchaseAccess - error.", err.Error())
		return err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var res FirstPurchaseAccess
		if err := rows.Scan(&res.accessId); err != nil {
			return err
		}
		sqlQueryo := "SELECT access_id from t_order where access_id = ? and FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') <> ? and pay_method=? LIMIT 0,1;"
		err := db.QueryRow(sqlQueryo, res.accessId, f.currentDate, paymode).Scan(&res.exist)
		if err != nil {
			if err == sql.ErrNoRows {
				count++
				continue
			}
			fmt.Println("RevenueDay.RevenueDayWechatFirstPurchaseAccess - error.", err.Error())
			return err
		}
		if res.exist == "" {
			count++
		}
	}
	if paymode == PaymodeAliPay {
		f.alipayFirstPurchaseAccess = count
		// fmt.Println("RevenueDay.RevenueDayWechatFirstPurchaseAccess - alipayFirstPurchaseAccess.", f.revenue.wechatFirstPurchaseAccess)
	} else {
		f.wechatFirstPurchaseAccess = count
		// fmt.Println("RevenueDay.RevenueDayWechatFirstPurchaseAccess - wechatFirstPurchaseAccess.", f.revenue.wechatFirstPurchaseAccess)
	}

	return nil
}

// RevenueDayPresentCount - 赠送份数.
func (f *RevenueDay) RevenueDayPresentCount(db *sql.DB) error {
	sqlQuery := "SELECT IFNULL(SUM(amount), 0) as presentCount from t_bought_package where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and activity_type=1; "
	err := db.QueryRow(sqlQuery, f.currentDate).Scan(&f.presentCount)
	if err != nil {
		fmt.Println("RevenueDay.RevenueDayPresentCount - error.", err.Error())
		return err
	}
	// fmt.Println("RevenueDay.RevenueDayPresentCount - presentCount.", f.presentCount)
	return nil
}

// RevenueDayPurchasePackage - 购买单份套餐次数	购买5份套餐次数	购买10份套餐次数	购买50份套餐次数	购买100份套餐次数.
func (f *RevenueDay) RevenueDayPurchasePackage(db *sql.DB, amount int) (int, error) {
	sqlQuery := "SELECT IFNULL(count(1),0) from t_bought_package where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and activity_type=0 and amount=?; "
	var count int
	err := db.QueryRow(sqlQuery, f.currentDate, amount).Scan(&count)
	if err != nil {
		fmt.Println("RevenueDay.RevenueDayPurchasePackage - error.", err.Error())
		return 0, err
	}
	// fmt.Println("RevenueDay.RevenueDayPurchasePackage - purchasePackage.", amount, "=>", count)
	return count, nil
}

// AverageAmount30Day - 30日均值.
func (f *RevenueDay) AverageAmount30Day(db *sql.DB, update bool) error {
	sqlQuery := "select sum(totalAmount)/30 as averageAmount30day from t_report_revenue where currentDate > ? and currentDate <= ?;"
	current, n30day := NDayDate(f.currentDate, -30)
	fmt.Println("RevenueDay.AverageAmount30Day - ", current, n30day)
	err := db.QueryRow(sqlQuery, n30day, current).Scan(&f.averageAmount30day)
	if err != nil {
		fmt.Println("RevenueDay.AverageAmount30Day - error.", err.Error())
		return err
	}
	fmt.Println("RevenueDay.AverageAmount30Day - averageAmount30day.", f.averageAmount30day)
	if update {
		sqlQuery = "update t_report_revenue set averageAmount30day = ? where currentDate = ?;"
		_, err := db.Exec(sqlQuery, f.averageAmount30day, current)
		if err != nil {
			return err
		}	
	}	
	return nil
}

// RevenueRemove -
func (f *RevenueDay) RevenueRemove(reportDb *sql.DB) error {
	primaryKeyDateName := utils.FormatDate(f.currentDate)
	sqlQuery := "delete from t_report_revenue where currentDate = ?; "
	_, err := reportDb.Exec(sqlQuery, primaryKeyDateName)
	if err != nil {
		fmt.Println("RevenueDay.RevenueRemove - error", err.Error())
		return err
	}
	// count, _ := res.RowsAffected()
	// fmt.Println("RevenueDay.RevenueRemove - delete rows: ", count)
	return nil
}

// RevenueInsert -
func (b *RevenueDay) RevenueInsert(reportDb *sql.DB) error {
	sqlQuery := "insert into t_report_revenue(currentDate, totalAmount, stockAll, stockCount, stockPurchaseCount, stockPresentCount, expiredPurchaseCount, expiredPresentCount, expendCount, " +
		"wechatTransAmount, wechatTransAccess, wechatTransCount, wechatRepurchaseAccess, wechatFirstPurchaseAccess, " +
		"alipayTransAmount, alipayTransAccess, alipayTransCount, alipayRepurchaseAccess, alipayFirstPurchaseAccess, " +
		"presentCount, purchasePackageAmount1, purchasePackageAmount5, purchasePackageAmount10, purchasePackageAmount50, purchasePackageAmount100, averageAmount30day, createTime) " +
		"values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	primaryKeyDateName := utils.FormatDate(b.currentDate)
	res, err := reportDb.Exec(sqlQuery, primaryKeyDateName, b.totalAmount, b.stockAll, b.stockCount, b.stockPurchaseCount, b.stockPresentCount, b.expiredPurchaseCount, b.expiredPresentCount, b.expendCount,
		b.wechatTransAmount, b.wechatTransAccess, b.wechatTransCount, b.wechatRepurchaseAccess, b.wechatFirstPurchaseAccess,
		b.alipayTransAmount, b.alipayTransAccess, b.alipayTransCount, b.alipayRepurchaseAccess, b.alipayFirstPurchaseAccess,
		b.presentCount, b.purchasePackageAmount1, b.purchasePackageAmount5, b.purchasePackageAmount10, b.purchasePackageAmount50, b.purchasePackageAmount100, b.averageAmount30day,
		time.Now(),
	)
	if err != nil {
		fmt.Println("RevenueDay.RevenueInsert - error", err.Error())
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		fmt.Println("RevenueDay.RevenueInsert - LastInsertId error", err.Error())
		return err
	}
	// fmt.Println("RevenueDay.RevenueInsert - insert id ", id)
	return nil
}

// Revenue - 营收维度分析.
func (f *FlashSignApp) Revenue(lastDate string) {
	o := &RevenueDay{currentDate: lastDate}
	o.RevenueDayTotalAmount(f.db)
	o.RevenueDayStockAll(f.db)
	o.RevenueDayStockCount(f.db)
	o.RevenueDayExpiredPurchaseCount(f.db)
	o.RevenueDayExpiredPresentCount(f.db)
	o.RevenueDayExpendCount(f.db)
	o.RevenueDayWechatTrans(f.db, PaymodeWechat)
	o.RevenueDayWechatRepurchaseAccess(f.db, PaymodeWechat)
	o.RevenueDayWechatFirstPurchaseAccess(f.db, PaymodeWechat)
	o.RevenueDayWechatTrans(f.db, PaymodeAliPay)
	o.RevenueDayWechatRepurchaseAccess(f.db, PaymodeAliPay)
	o.RevenueDayWechatFirstPurchaseAccess(f.db, PaymodeAliPay)
	o.RevenueDayPresentCount(f.db)
	o.purchasePackageAmount1, _ = o.RevenueDayPurchasePackage(f.db, PurchasePackageAmount1)
	o.purchasePackageAmount5, _ = o.RevenueDayPurchasePackage(f.db, PurchasePackageAmount5)
	o.purchasePackageAmount10, _ = o.RevenueDayPurchasePackage(f.db, PurchasePackageAmount10)
	o.purchasePackageAmount50, _ = o.RevenueDayPurchasePackage(f.db, PurchasePackageAmount50)
	o.purchasePackageAmount100, _ = o.RevenueDayPurchasePackage(f.db, PurchasePackageAmount100)
	// fmt.Println("FlashSignApp.Revenue - ", o)
	o.RevenueRemove(f.reportDb)
	o.RevenueInsert(f.reportDb)
	o.AverageAmount30Day(f.reportDb, true)	
}
