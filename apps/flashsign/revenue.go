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

// RevenueDay - 营收维度分析.
type RevenueDay struct {
	currentDate string  // 日期.
	totalAmount float64 // 当日总收入.
	stockCount  int     // 当日库存份数.
	expendCount int     // 当日消耗份数.

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
}

// NewRevenueDay -
func NewRevenueDay() *RevenueDay {
	return &RevenueDay{}
}

// RevenueDayTotalAmount - 当日总收入.
func (f *FlashSignApp) RevenueDayTotalAmount() error {
	sqlQuery := "SELECT SUM(price-dis_amount) as totalAmount FROM t_order where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and status = 2 "
	err := f.db.QueryRow(sqlQuery, f.currentDate).Scan(&f.revenue.totalAmount)
	if err != nil {
		fmt.Println("FlashSignApp.RevenueDayTotalAmount - error.", err.Error())
		return err
	}
	fmt.Println("FlashSignApp.RevenueDayTotalAmount - totalAmount.", f.revenue.totalAmount)
	return nil
}

// RevenueDayStockCount - 当日库存份数.
func (f *FlashSignApp) RevenueDayStockCount() error {
	sqlQuery := "SELECT SUM(count) as stockCount from t_bought_package where status = 0 and activity_type = 0 and FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? "
	err := f.db.QueryRow(sqlQuery, f.currentDate).Scan(&f.revenue.stockCount)
	if err != nil {
		fmt.Println("FlashSignApp.RevenueDayStockCount - error.", err.Error())
		return err
	}
	fmt.Println("FlashSignApp.RevenueDayStockCount - stockQuantity.", f.revenue.stockCount)
	return nil
}

type PackageDeductionKindCount struct {
	count int // 份数.
	kind  int // group by.
}

// RevenueDayExpendCount - 当日消耗份数: 当日套餐抵扣0（划扣） -  撤回3（撤销划扣）.
func (f *FlashSignApp) RevenueDayExpendCount() error {
	sqlQuery := "SELECT SUM(amount) as count, package_deduction_kind as kind  from t_deduction_record WHERE FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? GROUP BY package_deduction_kind; "

	rows, err := f.db.Query(sqlQuery, f.currentDate)
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
			f.revenue.expendCount = count - reduct
		} else {
			f.revenue.expendCount = count
		}
	}
	fmt.Println("FlashSignApp.RevenueDayExpendCount - expendCount.", f.revenue.expendCount)
	return nil
}

// RevenueDayWechatTrans - 微信当日交易金额 微信当日交易用户 微信当日交易笔数.
func (f *FlashSignApp) RevenueDayWechatTrans(paymode int) error {
	sqlQuery := "SELECT IFNULL(SUM(price-dis_amount),0) as wechatTransAmount, IFNULL(COUNT(DISTINCT access_id),'') as wechatTransAccess, IFNULL(COUNT(1),0) as wechatTransCount FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and `status` = 2 and pay_method=?; "
	var err error
	if paymode == PaymodeAliPay {
		err = f.db.QueryRow(sqlQuery, f.currentDate, paymode).Scan(&f.revenue.alipayTransAmount, &f.revenue.alipayTransAccess, &f.revenue.alipayTransCount)
	} else {
		err = f.db.QueryRow(sqlQuery, f.currentDate, paymode).Scan(&f.revenue.wechatTransAmount, &f.revenue.wechatTransAccess, &f.revenue.wechatTransCount)
	}
	if err != nil {
		fmt.Println("FlashSignApp.RevenueDayWechatTrans - error.", err.Error())
		return err
	}
	if paymode == PaymodeAliPay {
		fmt.Println("FlashSignApp.RevenueDayWechatTrans - alipayTrans.", f.revenue.alipayTransAmount, f.revenue.alipayTransAccess, f.revenue.alipayTransCount)
	} else {
		fmt.Println("FlashSignApp.RevenueDayWechatTrans - wechatTrans.", f.revenue.wechatTransAmount, f.revenue.wechatTransAccess, f.revenue.wechatTransCount)
	}
	return nil
}

// RevenueDayWechatRepurchaseAccess - 微信日复购用户.
func (f *FlashSignApp) RevenueDayWechatRepurchaseAccess(paymode int) error {
	sqlQuery := "SELECT COUNT(1) as repurchaseAccess from (SELECT COUNT(access_id) as access_ids FROM t_order where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and status = 2 and pay_method=? GROUP BY access_id) t where t.access_ids>1; "
	var access int
	err := f.db.QueryRow(sqlQuery, f.currentDate, paymode).Scan(&access)
	if err != nil {
		fmt.Println("FlashSignApp.RevenueDayWechatRepurchaseAccess - error.", err.Error())
		return err
	}
	if paymode == PaymodeAliPay {
		f.revenue.alipayRepurchaseAccess = access
		fmt.Println("FlashSignApp.RevenueDayWechatRepurchaseAccess - alipayRepurchaseAccess.", f.revenue.alipayRepurchaseAccess)
	} else {
		f.revenue.wechatRepurchaseAccess = access
		fmt.Println("FlashSignApp.RevenueDayWechatRepurchaseAccess - repurchaseAccess.", f.revenue.wechatRepurchaseAccess)
	}
	return nil
}

// FirstPurchaseAccess -
type FirstPurchaseAccess struct {
	accessId string
	exist    string
}

// RevenueDayWechatFirstPurchaseAccess - 微信日新购用户 - 微信平台首次购买的记录（程序实现：先建立首次购买表按平台，在查找是否存在购买）.
func (f *FlashSignApp) RevenueDayWechatFirstPurchaseAccess(paymode int) error {
	sqlQuery := "SELECT DISTINCT access_id as accessId FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and `status` = 2 and pay_method=? GROUP BY access_id;  "
	rows, err := f.db.Query(sqlQuery, f.currentDate, paymode)
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
		err := f.db.QueryRow(sqlQueryo, res.accessId, f.currentDate, paymode).Scan(&res.exist)
		if err != nil {
			if err == sql.ErrNoRows {
				count++
				continue
			}
			fmt.Println("FlashSignApp.RevenueDayWechatFirstPurchaseAccess - error.", err.Error())
			return err
		}
		if res.exist == "" {
			count++
		}
	}
	if paymode == PaymodeAliPay {
		f.revenue.alipayFirstPurchaseAccess = count
		fmt.Println("FlashSignApp.RevenueDayWechatFirstPurchaseAccess - alipayFirstPurchaseAccess.", f.revenue.wechatFirstPurchaseAccess)
	} else {
		f.revenue.wechatFirstPurchaseAccess = count
		fmt.Println("FlashSignApp.RevenueDayWechatFirstPurchaseAccess - wechatFirstPurchaseAccess.", f.revenue.wechatFirstPurchaseAccess)
	}

	return nil
}

// RevenueDayPresentCount - 赠送份数.
func (f *FlashSignApp) RevenueDayPresentCount() error {
	sqlQuery := "SELECT IFNULL(SUM(amount), 0) as presentCount from t_bought_package where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and activity_type=1; "
	err := f.db.QueryRow(sqlQuery, f.currentDate).Scan(&f.revenue.presentCount)
	if err != nil {
		fmt.Println("FlashSignApp.RevenueDayPresentCount - error.", err.Error())
		return err
	}
	fmt.Println("FlashSignApp.RevenueDayPresentCount - presentCount.", f.revenue.presentCount)
	return nil
}

// RevenueDayPurchasePackage - 购买单份套餐次数	购买5份套餐次数	购买10份套餐次数	购买50份套餐次数	购买100份套餐次数.
func (f *FlashSignApp) RevenueDayPurchasePackage(amount int) (int, error) {
	sqlQuery := "SELECT IFNULL(count(1),0) from t_bought_package where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and activity_type=0 and amount=?; "
	var count int
	err := f.db.QueryRow(sqlQuery, f.currentDate, amount).Scan(&count)
	if err != nil {
		fmt.Println("FlashSignApp.RevenueDayPurchasePackage - error.", err.Error())
		return 0, err
	}
	fmt.Println("FlashSignApp.RevenueDayPurchasePackage - purchasePackage.", amount, "=>", count)
	return count, nil
}

// RevenueInsert -
func (f *FlashSignApp) RevenueInsert() error {
	db, err := DBConnect(f.driverName, f.dbSourceUrlBaasReport(), 3)
	if err != nil {
		return err
	}
	defer db.Close()
	sqlQuery := "insert into t_report_revenue(currentDate, totalAmount, stockCount, expendCount, " +
		"wechatTransAmount, wechatTransAccess, wechatTransCount, wechatRepurchaseAccess, wechatFirstPurchaseAccess, " +
		"alipayTransAmount, alipayTransAccess, alipayTransCount, alipayRepurchaseAccess, alipayFirstPurchaseAccess, " +
		"presentCount, purchasePackageAmount1, purchasePackageAmount5, purchasePackageAmount10, purchasePackageAmount50, purchasePackageAmount100, createTime) " +
		"values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	b := f.revenue
	primaryKeyDateName := utils.FormatDate(b.currentDate)
	res, err := db.Exec(sqlQuery, primaryKeyDateName, b.totalAmount, b.stockCount, b.expendCount,
		b.wechatTransAmount, b.wechatTransAccess, b.wechatTransCount, b.wechatRepurchaseAccess, b.wechatFirstPurchaseAccess,
		b.alipayTransAmount, b.alipayTransAccess, b.alipayTransCount, b.alipayRepurchaseAccess, b.alipayFirstPurchaseAccess,
		b.presentCount, b.purchasePackageAmount1, b.purchasePackageAmount5, b.purchasePackageAmount10, b.purchasePackageAmount50, b.purchasePackageAmount100, time.Now(),
	)
	if err != nil {
		fmt.Println("FlashSignApp.RevenueInsert - error", err.Error())
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("FlashSignApp.RevenueInsert - LastInsertId error", err.Error())
		return err
	}
	fmt.Println("FlashSignApp.RevenueInsert - insert id ", id)
	return nil
}

// Revenue - 营收维度分析.
func (f *FlashSignApp) Revenue() {
	f.RevenueDayTotalAmount()
	f.RevenueDayStockCount()
	f.RevenueDayExpendCount()
	f.RevenueDayWechatTrans(PaymodeWechat)
	f.RevenueDayWechatRepurchaseAccess(PaymodeWechat)
	f.RevenueDayWechatFirstPurchaseAccess(PaymodeWechat)
	f.RevenueDayWechatTrans(PaymodeAliPay)
	f.RevenueDayWechatRepurchaseAccess(PaymodeAliPay)
	f.RevenueDayWechatFirstPurchaseAccess(PaymodeAliPay)
	f.RevenueDayPresentCount()
	f.revenue.purchasePackageAmount1, _ = f.RevenueDayPurchasePackage(PurchasePackageAmount1)
	f.revenue.purchasePackageAmount5, _ = f.RevenueDayPurchasePackage(PurchasePackageAmount5)
	f.revenue.purchasePackageAmount10, _ = f.RevenueDayPurchasePackage(PurchasePackageAmount10)
	f.revenue.purchasePackageAmount50, _ = f.RevenueDayPurchasePackage(PurchasePackageAmount50)
	f.revenue.purchasePackageAmount100, _ = f.RevenueDayPurchasePackage(PurchasePackageAmount100)
	f.RevenueInsert()
}
