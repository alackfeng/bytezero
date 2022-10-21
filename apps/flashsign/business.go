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

// Contract -
type Contract struct {
	kind  int
	count int
}

// BusinessDay - 业务维度分析.
// # 日期 当日合同签署总数 自定义类合同当日签署数 自定义类合同占比 模板类合同当日签署份数	模板类合同当日占比	模板类借贷类合同签署份数 模板类借贷类合同当日占比 模板类非借贷类合同签署份数 模板类非借贷类合同当日占比	法律增值业务.
type BusinessDay struct {
	currentDate                 string  // 日期.
	signSuccessTotalCount       int     // 当日合同签署总数.
	signTotalCount              int     // 当日合同签署次数.
	customContractSignCount     int     // 自定义类合同当日签署数.
	customContractSignPercent   float64 // 自定义类合同占比
	templateContractSignCount   int     // 模板类合同当日签署份数.
	templateContractSignPercent float64 // 模板类合同当日占比.

	templateLoanContractSignCount     int     // 模板类借贷类合同签署份数.
	templateLoanContractSignPercent   float64 // 模板类借贷类合同当日占比.
	templateNoLoanContractSignCount   int     // 模板类非借贷类合同签署份数.
	templateNoLoanContractSignPercent float64 // 模板类非借贷类合同当日占比.
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

// BusinessDaySignSuccessTotalCount - 当日合同签署总数(operate_type=4) ：状态为成功的记录.
func (f *BusinessDay) BusinessDaySignSuccessTotalCount(db *sql.DB) error {
	sqlQuery := "SELECT IFNULL(count(1),0) as signSuccessTotalCount from t_contract_operate_record where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and operate_type=4; "
	err := db.QueryRow(sqlQuery, f.currentDate).Scan(&f.signSuccessTotalCount)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDaySignSuccessTotalCount - error.", err.Error())
		return err
	}
	// fmt.Println("BusinessDay.BusinessDaySignSuccessTotalCount - signSuccessTotalCount.", f.business.signSuccessTotalCount)
	return nil
}

// BusinessDaySignTotalCount - 当日合同签署次数(operate_type=1) : 一份合同存在多人签署，每个操作都算.
func (f *BusinessDay) BusinessDaySignTotalCount(db *sql.DB) error {
	sqlQuery := "SELECT IFNULL(count(1),0) as signTotalCount from t_contract_operate_record where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and operate_type=1; "
	err := db.QueryRow(sqlQuery, f.currentDate).Scan(&f.signTotalCount)
	if err != nil {
		fmt.Println("BusinessDay.BusinessDaySignTotalCount - error.", err.Error())
		return err
	}
	// fmt.Println("BusinessDay.BusinessDaySignTotalCount - signTotalCount.", f.business.signTotalCount)
	return nil
}

// BusinessDayContractSignCount - 自定义类合同当日签署数 自定义类合同占比 模板类合同当日签署份数 模板类合同当日占比 - template_type= 0自定义类合同 | 1模板类合同.
func (f *BusinessDay) BusinessDayContractSignCount(db *sql.DB) error {
	sqlQuery := "select template_type as kind, SUM(template_count) as count from  (SELECT template_id, count(1) as template_count, CASE WHEN template_id>0 THEN 1 ELSE 0 END template_type from t_contract where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? GROUP BY template_id) t GROUP BY template_type; "
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
	sqlQuery := "select t.template_type as kind, SUM(t.template_count) as count from (SELECT template_id, count(1) as template_count, case when template_id in (SELECT id from t_template where template_class_id in (SELECT id from t_template_class where name like '%借%')) then 2 else 3 end template_type from t_contract where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? GROUP BY template_id HAVING template_id>0) t GROUP BY template_type; "
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
	sqlQuery := "insert into t_report_business(currentDate, signSuccessTotalCount, signTotalCount, customContractSignCount, customContractSignPercent, templateContractSignCount, templateContractSignPercent, templateLoanContractSignCount, templateLoanContractSignPercent, templateNoLoanContractSignCount, templateNoLoanContractSignPercent, createTime) values (?,?,?,?,?,?,?,?,?,?,?,?)"
	primaryKeyDateName := utils.FormatDate(b.currentDate)
	res, err := reportDb.Exec(sqlQuery, primaryKeyDateName, b.signSuccessTotalCount, b.signTotalCount,
		b.customContractSignCount, b.customContractSignPercent,
		b.templateContractSignCount, b.templateContractSignPercent,
		b.templateLoanContractSignCount, b.templateLoanContractSignPercent,
		b.templateNoLoanContractSignCount, b.templateNoLoanContractSignPercent,
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
	o.BusinessRemove(f.reportDb)
	o.BusinessInsert(f.reportDb)
	return nil
}
