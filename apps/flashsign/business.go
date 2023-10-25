package flashsign

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alackfeng/bytezero/cores/utils"
)

// ContractKind -
const (
	ContractKindCustom         = 0 // 自定义类合同.
	ContractKindTemplate       = 1 // 模板类合同.
	ContractKindTemplateLoan   = 2 // 模板类借贷类合同.
	ContractKindTemplateNoLoan = 3 // 模板类非借贷类合同.
)

// ContractType -
const (
	ContractTypeNormal  = 0 // 常规合同.
	ContractTypePresent = 1 // 体验合同.
	ContractTypeOther   = 2 // 其他合同.
)

// Contract -
type Contract struct {
	kind  int
	count int
}

// BusinessDay - 业务维度分析.
// # 日期 当日合同签署总数 自定义类合同当日签署数 自定义类合同占比 模板类合同当日签署份数	模板类合同当日占比	模板类借贷类合同签署份数 模板类借贷类合同当日占比 模板类非借贷类合同签署份数 模板类非借贷类合同当日占比	法律增值业务.
type BusinessDay struct {
	currentDate                  string  // 日期.
	signSuccessTotalCount        int     // 当日常规合同签署总数.
	signInvalidCount             int     // 当日解除合同签署份数.
        signSettlementCount          int     // 当日结清合同签署份数.
	signTotalCount               int     // 当日常规合同签署次数.
	signSuccessPresentTotalCount int     // 当日体验合同签署总数.
	signPresentTotalCount        int     // 当日体验合同签署次数.
	customContractSignCount      int     // 自定义类合同当日签署数.
	customContractSignPercent    float64 // 自定义类合同占比
	templateContractSignCount    int     // 模板类合同当日签署份数.
	templateContractSignPercent  float64 // 模板类合同当日占比.

	templateLoanContractSignCount     int     // 模板类借贷类合同签署份数.
	templateLoanContractSignPercent   float64 // 模板类借贷类合同当日占比.
	templateNoLoanContractSignCount   int     // 模板类非借贷类合同签署份数.
	templateNoLoanContractSignPercent float64 // 模板类非借贷类合同当日占比.

	// 2023-10-25 - .
	notarizationBuyCountDay		int // 当日购买公证书数量.
	notarizationBuyCountSum		int // 累记购买公证书数量.
	notarizationReqSuccessCountDay	int // 当日公证书出证成功数量.
	notarizationReqSuccessCountSum	int // 累记公证书出证成功数量.
}

// NewBusinessDay -
func NewBusinessDay(name string) *BusinessDay {
	return &BusinessDay{
		currentDate: name,
	}
}

// String -
func (b BusinessDay) String() string {
	return fmt.Sprintf("%s: [Success:%d, Count:%d, Custom:%d-%.2f, Template:%d-%.2f, Loan:%d-%.2f, NoLoan:%d-%.2f]", b.currentDate, b.signSuccessTotalCount, b.signTotalCount,
		b.customContractSignCount, b.customContractSignPercent,
		b.templateContractSignCount, b.templateContractSignPercent,
		b.templateLoanContractSignCount, b.templateLoanContractSignPercent,
		b.templateNoLoanContractSignCount, b.templateNoLoanContractSignPercent)
}

// BusinessDaySignSuccessTotalCount - 当日合同签署总数(operate_type=4) - 已完成4, contract_kind=0正式.
func (f *BusinessDay) BusinessDaySignSuccessTotalCount(db *sql.DB) error {
	sqlQuery := "select max(contract_kind) as type, count(contract_id) as count from t_contract_operate_record t1 INNER JOIN t_contract t2 ON t1.contract_id=t2.id where FROM_UNIXTIME(t1.create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and t1.operate_type=4 group by contract_kind; "
	rows, err := db.Query(sqlQuery, f.currentDate)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDaySignSuccessTotalCount - error.", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var res Contract
		if err := rows.Scan(&res.kind, &res.count); err != nil {
			return err
		}
		if res.kind == ContractTypeNormal {
			f.signSuccessTotalCount = res.count
		} else if res.kind == ContractTypePresent {
			f.signSuccessPresentTotalCount = res.count
		}
	}
	// fmt.Println("BusinessDay.BusinessDaySignSuccessTotalCount - signSuccessTotalCount.", f.business.signSuccessTotalCount)
	return nil
}

// BusinessDaySignInvalidCount - 当日解除合同签署份数.
func (f *BusinessDay) BusinessDaySignInvalidCount(db *sql.DB) error {
	sqlQuery := "SELECT COUNT(id) AS count FROM t_contract  WHERE FROM_UNIXTIME(finish_time DIV 1000, '%Y-%m-%d 00:00:00') = ? AND contract_type = 2 AND status = 3; "
	rows, err := db.Query(sqlQuery, f.currentDate)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDaySignInvalidCount - error.", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			return err
		}
		f.signInvalidCount = count
	 	fmt.Println("BusinessDay.BusinessDaySignInvalidCount - signInvalidCount.", f.signInvalidCount)
	}
	return nil
}

// BusinessDaySignSettlementCount - 当日结清合同签署份数.
func (f *BusinessDay) BusinessDaySignSettlementCount(db *sql.DB) error {
	sqlQuery := "SELECT COUNT(id) AS count FROM t_contract  WHERE FROM_UNIXTIME(finish_time DIV 1000, '%Y-%m-%d 00:00:00') = ? AND contract_type = 4 AND status = 3; "
	rows, err := db.Query(sqlQuery, f.currentDate)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDaySignInvalidCount - error.", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			return err
		}
		f.signSettlementCount = count
	 	fmt.Println("BusinessDay.BusinessDaySignSettlementCount - signSettlementCount.", f.signSettlementCount)
	}
	return nil
}

// BusinessDaySignTotalCount - 当日合同签署次数(operate_type=1) : 一份合同存在多人签署，每个操作都算, contract_kind=0正式合同.
func (f *BusinessDay) BusinessDaySignTotalCount(db *sql.DB) error {
	sqlQuery := "select max(contract_kind) as type, count(contract_id) as count from t_contract_operate_record t1 INNER JOIN t_contract t2 ON t1.contract_id=t2.id where FROM_UNIXTIME(t1.create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and t1.operate_type=1 group by contract_kind; "
	rows, err := db.Query(sqlQuery, f.currentDate)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDaySignTotalCount - error.", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var res Contract
		if err := rows.Scan(&res.kind, &res.count); err != nil {
			return err
		}
		if res.kind == ContractTypeNormal {
			f.signTotalCount = res.count
		} else if res.kind == ContractTypePresent {
			f.signPresentTotalCount = res.count
		}
	}
	// fmt.Println("BusinessDay.BusinessDaySignTotalCount - signTotalCount.", f.business.signTotalCount)
	return nil
}

// BusinessDayContractSignCount - 自定义类合同当日签署数 自定义类合同占比 模板类合同当日签署份数 模板类合同当日占比 - template_type= 0自定义类合同 | 1模板类合同（contract_kind=0为正常合同）.
func (f *BusinessDay) BusinessDayContractSignCount(db *sql.DB) error {
	sqlQuery := "select template_type as kind, SUM(template_count) as count from  (SELECT template_id, count(1) as template_count, CASE WHEN template_id>0 THEN 1 ELSE 0 END template_type from t_contract where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and contract_kind=0 GROUP BY template_id) t GROUP BY template_type; "
	rows, err := db.Query(sqlQuery, f.currentDate)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDayContractSignCount - error.", err.Error())
		return err
	}
	defer rows.Close()
	total := 0
	for rows.Next() {
		var res Contract
		if err := rows.Scan(&res.kind, &res.count); err != nil {
			return err
		}
		if res.kind == ContractKindCustom {
			f.customContractSignCount = res.count
		} else if res.kind == ContractKindTemplate {
			f.templateContractSignCount = res.count
		}
		total += res.count
	}
	if total != 0 {
		f.customContractSignPercent = utils.CalcPercent(f.customContractSignCount*100, total)
		f.templateContractSignPercent = utils.CalcPercent(f.templateContractSignCount*100, total)
	}
	// fmt.Println("BusinessDay.BusinessDayContractSignCount - ", f.business)
	return nil
}

// BusinessDayTemplateContractSignCount - 模板类借贷类合同签署份数	模板类借贷类合同当日占比 模板类非借贷类合同签署份数	模板类非借贷类合同当日占比.
func (f *BusinessDay) BusinessDayTemplateContractSignCount(db *sql.DB) error {
	sqlQuery := "select t.template_type as kind, SUM(t.template_count) as count from (SELECT template_id, count(1) as template_count, case when template_id in (SELECT id from t_template where template_class_id in (SELECT id from t_template_class where name like '%借%')) then 2 else 3 end template_type from t_contract where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and contract_kind=0 GROUP BY template_id HAVING template_id>0) t GROUP BY template_type; "
	rows, err := db.Query(sqlQuery, f.currentDate)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDayTemplateContractSignCount - error.", err.Error())
		return err
	}
	defer rows.Close()
	total := 0
	for rows.Next() {
		var res Contract
		if err := rows.Scan(&res.kind, &res.count); err != nil {
			return err
		}
		if res.kind == ContractKindTemplateLoan {
			f.templateLoanContractSignCount = res.count
		} else if res.kind == ContractKindTemplateNoLoan {
			f.templateNoLoanContractSignCount = res.count
		}
		total += res.count
	}
	if total != 0 {
		f.templateLoanContractSignPercent = utils.CalcPercent(f.templateLoanContractSignCount*100, total)
		f.templateNoLoanContractSignPercent = utils.CalcPercent(f.templateNoLoanContractSignCount*100, total)
	}
	// fmt.Println("BusinessDay.BusinessDayTemplateContractSignCount - ", f.business)
	return nil
}

// BusinessDayNotarizationBuyCountDay - 当日购买公证书数量.
func (f *BusinessDay) BusinessDayNotarizationBuyCountDay(db *sql.DB) error {
	sqlQuery := "SELECT COUNT(id) AS count FROM notarization_buy_record  WHERE FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ?; "
	rows, err := db.Query(sqlQuery, f.currentDate)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDayNotarizationBuyCountDay - error.", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			return err
		}
		f.notarizationBuyCountDay = count
	 	fmt.Println("BusinessDay.BusinessDayNotarizationBuyCountDay - notarizationBuyCountDay.", f.notarizationBuyCountDay)
	}
	return nil
}

// BusinessDayNotarizationBuyCountSum - 累记购买公证书数量.
func (f *BusinessDay) BusinessDayNotarizationBuyCountSum(db *sql.DB) error {
	sqlQuery := "SELECT COUNT(id) AS count FROM notarization_buy_record; "
	rows, err := db.Query(sqlQuery)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDayNotarizationBuyCountSum - error.", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			return err
		}
		f.notarizationBuyCountSum = count
	 	fmt.Println("BusinessDay.BusinessDayNotarizationBuyCountSum - notarizationBuyCountSum.", f.notarizationBuyCountSum)
	}
	return nil
}

// BusinessDayNotarizationReqSuccessCountDay - 当日公证书出证成功数量.
func (f *BusinessDay) BusinessDayNotarizationReqSuccessCountDay(db *sql.DB) error {
	sqlQuery := "SELECT COUNT(id) AS count FROM notarization_request_record  WHERE FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? AND status = 1 ; "
	rows, err := db.Query(sqlQuery, f.currentDate)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDayNotarizationReqSuccessCountDay - error.", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			return err
		}
		f.notarizationReqSuccessCountDay = count
	 	fmt.Println("BusinessDay.BusinessDayNotarizationReqSuccessCountDay - notarizationReqSuccessCountDay.", f.notarizationReqSuccessCountDay)
	}
	return nil
}

// BusinessDayNotarizationReqSuccessCountSum - 累记公证书出证成功数量.
func (f *BusinessDay) BusinessDayNotarizationReqSuccessCountSum(db *sql.DB) error {
	sqlQuery := "SELECT COUNT(id) AS count FROM notarization_request_record  WHERE status = 1 ; "
	rows, err := db.Query(sqlQuery)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDayNotarizationReqSuccessCountSum - error.", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			return err
		}
		f.notarizationReqSuccessCountSum = count
	 	fmt.Println("BusinessDay.BusinessDayNotarizationReqSuccessCountSum - notarizationReqSuccessCountSum.", f.notarizationReqSuccessCountSum)
	}
	return nil
}

// BusinessUpdate - 
func (f *BusinessDay) BusinessUpdate(reportDb *sql.DB, field string) error {
	primaryKeyDateName := utils.FormatDate(f.currentDate)
	if field == "notarizationBuyCountDay" {
 		sqlQuery := "update t_report_business set notarizationBuyCountDay = ? where currentDate = ?; "
		_, err := reportDb.Exec(sqlQuery, f.notarizationBuyCountDay, primaryKeyDateName)
		if err != nil {
			fmt.Println("BusinessDay.BusinessUpdate - error.", field, err.Error())
			return err
		}
	} else if field == "notarizationReqSuccessCountDay" {
 		sqlQuery := "update t_report_business set notarizationReqSuccessCountDay = ? where currentDate = ?; "
		_, err := reportDb.Exec(sqlQuery, f.notarizationReqSuccessCountDay, primaryKeyDateName)
		if err != nil {
			fmt.Println("BusinessDay.BusinessUpdate - error.", field, err.Error())
			return err
		}
	} else {
		fmt.Println("BusinessDay.BusinessUpdate - no field.", field)
	}
	return nil
}


// RevenueRemove -
func (f *BusinessDay) BusinessRemove(reportDb *sql.DB) error {
	primaryKeyDateName := utils.FormatDate(f.currentDate)
	sqlQuery := "delete from t_report_business where currentDate = ?; "
	_, err := reportDb.Exec(sqlQuery, primaryKeyDateName)
	if err != nil {
		fmt.Println("BusinessDay.BusinessRemove - error", err.Error())
		return err
	}
	// count, _ := res.RowsAffected()
	// fmt.Println("BusinessDay.BusinessRemove - delete rows: ", count)
	return nil
}

// BusinessInsert -
func (b *BusinessDay) BusinessInsert(reportDb *sql.DB) error {
	sqlQuery := "insert into t_report_business(currentDate, signSuccessTotalCount, signInvalidCount, signSettlementCount, signTotalCount, signSuccessPresentTotalCount, signPresentTotalCount, customContractSignCount, customContractSignPercent, templateContractSignCount, templateContractSignPercent, templateLoanContractSignCount, templateLoanContractSignPercent, templateNoLoanContractSignCount, templateNoLoanContractSignPercent, notarizationBuyCountDay, notarizationBuyCountSum, notarizationReqSuccessCountDay, notarizationReqSuccessCountSum, createTime) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	primaryKeyDateName := utils.FormatDate(b.currentDate)
	res, err := reportDb.Exec(sqlQuery, primaryKeyDateName, b.signSuccessTotalCount, b.signInvalidCount, b.signSettlementCount, b.signTotalCount, b.signSuccessPresentTotalCount, b.signPresentTotalCount,
		b.customContractSignCount, b.customContractSignPercent,
		b.templateContractSignCount, b.templateContractSignPercent,
		b.templateLoanContractSignCount, b.templateLoanContractSignPercent,
		b.templateNoLoanContractSignCount, b.templateNoLoanContractSignPercent,
		b.notarizationBuyCountDay, b.notarizationBuyCountSum, b.notarizationReqSuccessCountDay, b.notarizationReqSuccessCountSum,
		time.Now(),
	)
	if err != nil {
		fmt.Println("BusinessDay.BusinessInsert - error", err.Error())
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		fmt.Println("BusinessDay.BusinessInsert - LastInsertId error", err.Error())
		return err
	}
	// fmt.Println("BusinessDay.BusinessInsert - insert id ", id)
	return nil
}

// Business -
func (f *FlashSignApp) Business(lastDate string) error {
	o := &BusinessDay{currentDate: lastDate}
	o.BusinessDaySignSuccessTotalCount(f.db)
	o.BusinessDaySignTotalCount(f.db)
	o.BusinessDayContractSignCount(f.db)
	o.BusinessDayTemplateContractSignCount(f.db)
	// fmt.Println("FlashSignApp.Business - ", o)
	o.BusinessDaySignInvalidCount(f.db)
	o.BusinessDaySignSettlementCount(f.db)
	o.BusinessDayNotarizationBuyCountDay(f.db)
	o.BusinessDayNotarizationBuyCountSum(f.db)
	o.BusinessDayNotarizationReqSuccessCountDay(f.db)
	o.BusinessDayNotarizationReqSuccessCountSum(f.db)
	o.BusinessRemove(f.reportDb)
	o.BusinessInsert(f.reportDb)
	return nil
}
