package flashsign

import (
	"fmt"
	"time"

	"github.com/alackfeng/bytezero/cores/utils"
)

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

// RevenueDayTotalAmount - 当日合同签署总数(operate_type=4) ：状态为成功的记录.
func (f *FlashSignApp) BusinessDaySignSuccessTotalCount() error {
	sqlQuery := "SELECT count(1) as signSuccessTotalCount from t_contract_operate_record where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and operate_type=4; "
	err := f.db.QueryRow(sqlQuery, f.business.currentDate).Scan(&f.business.signSuccessTotalCount)
	if err != nil {
		fmt.Println("FlashSignApp.BusinessDaySignSuccessTotalCount - error.", err.Error())
		return err
	}
	fmt.Println("FlashSignApp.BusinessDaySignSuccessTotalCount - signSuccessTotalCount.", f.business.signSuccessTotalCount)
	return nil
}

// BusinessDaySignTotalCount - 当日合同签署次数(operate_type=1) : 一份合同存在多人签署，每个操作都算.
func (f *FlashSignApp) BusinessDaySignTotalCount() error {
	sqlQuery := "SELECT count(1) as signTotalCount from t_contract_operate_record where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? and operate_type=1; "
	err := f.db.QueryRow(sqlQuery, f.business.currentDate).Scan(&f.business.signTotalCount)
	if err != nil {
		fmt.Println("FlashSignApp.BusinessDaySignTotalCount - error.", err.Error())
		return err
	}
	fmt.Println("FlashSignApp.BusinessDaySignTotalCount - signTotalCount.", f.business.signTotalCount)
	return nil
}

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

// BusinessDayContractSignCount - 自定义类合同当日签署数 自定义类合同占比 模板类合同当日签署份数 模板类合同当日占比 - template_type= 0自定义类合同 | 1模板类合同.
func (f *FlashSignApp) BusinessDayContractSignCount() error {
	sqlQuery := "select template_type as kind, SUM(template_count) as count from  (SELECT template_id, count(1) as template_count, CASE WHEN template_id>0 THEN 1 ELSE 0 END template_type from t_contract where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? GROUP BY template_id) t GROUP BY template_type; "
	rows, err := f.db.Query(sqlQuery, f.business.currentDate)
	if err != nil {
		fmt.Println("FlashSignApp.BusinessDayContractSignCount - error.", err.Error())
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
			f.business.customContractSignCount = res.count
		} else if res.kind == ContractKindTemplate {
			f.business.templateContractSignCount = res.count
		}
		total += res.count
	}
	if total != 0 {
		f.business.customContractSignPercent = utils.CalcPercent(f.business.customContractSignCount*100, total)
		f.business.templateContractSignPercent = utils.CalcPercent(f.business.templateContractSignCount*100, total)
	}
	fmt.Println("FlashSignApp.BusinessDayContractSignCount - ", f.business)
	return nil
}

// BusinessDayTemplateContractSignCount - 模板类借贷类合同签署份数	模板类借贷类合同当日占比 模板类非借贷类合同签署份数	模板类非借贷类合同当日占比.
func (f *FlashSignApp) BusinessDayTemplateContractSignCount() error {
	sqlQuery := "select t.template_type as kind, SUM(t.template_count) as count from (SELECT template_id, count(1) as template_count, case when template_id in (SELECT id from t_template where template_class_id in (SELECT id from t_template_class where name like '%借%')) then 2 else 3 end template_type from t_contract where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = ? GROUP BY template_id HAVING template_id>0) t GROUP BY template_type; "
	rows, err := f.db.Query(sqlQuery, f.business.currentDate)
	if err != nil {
		fmt.Println("FlashSignApp.BusinessDayTemplateContractSignCount - error.", err.Error())
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
			f.business.templateLoanContractSignCount = res.count
		} else if res.kind == ContractKindTemplateNoLoan {
			f.business.templateNoLoanContractSignCount = res.count
		}
		total += res.count
	}
	if total != 0 {
		f.business.templateLoanContractSignPercent = utils.CalcPercent(f.business.templateLoanContractSignCount*100, total)
		f.business.templateNoLoanContractSignPercent = utils.CalcPercent(f.business.templateNoLoanContractSignCount*100, total)
	}
	fmt.Println("FlashSignApp.BusinessDayTemplateContractSignCount - ", f.business)
	return nil
}

// dbSourceUrlBaasSeal -
func (f *FlashSignApp) dbSourceUrlBaasReport() string {
	return f.dbSourceUrl + f.dbNameBaasReport
}

// BusinessInsert -
func (f *FlashSignApp) BusinessInsert() error {
	db, err := DBConnect(f.driverName, f.dbSourceUrlBaasReport(), 3)
	if err != nil {
		return err
	}
	defer db.Close()
	sqlQuery := "insert into t_report_business(currentDate, signSuccessTotalCount, signTotalCount, customContractSignCount, customContractSignPercent, templateContractSignCount, templateContractSignPercent, templateLoanContractSignCount, templateLoanContractSignPercent, templateNoLoanContractSignCount, templateNoLoanContractSignPercent, createTime) values (?,?,?,?,?,?,?,?,?,?,?,?)"
	b := f.business
	primaryKeyDateName := utils.FormatDate(b.currentDate)
	res, err := db.Exec(sqlQuery, primaryKeyDateName, b.signSuccessTotalCount, b.signTotalCount,
		b.customContractSignCount, b.customContractSignPercent,
		b.templateContractSignCount, b.templateContractSignPercent,
		b.templateLoanContractSignCount, b.templateLoanContractSignPercent,
		b.templateNoLoanContractSignCount, b.templateNoLoanContractSignPercent,
		time.Now(),
	)
	if err != nil {
		fmt.Println("FlashSignApp.BusinessInsert - error", err.Error())
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("FlashSignApp.BusinessInsert - LastInsertId error", err.Error())
		return err
	}
	fmt.Println("FlashSignApp.BusinessInsert - insert id ", id)
	return nil
}

// Business -
func (f *FlashSignApp) Business() error {
	f.BusinessDaySignSuccessTotalCount()
	f.BusinessDaySignTotalCount()
	f.BusinessDayContractSignCount()
	f.BusinessDayTemplateContractSignCount()
	f.BusinessInsert()
	return nil
}
