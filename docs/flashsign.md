-- 营收维度分析
# 日期	当日总收入	当日库存份数	当日消耗份数	支付宝当日交易金额	支付宝当日交易用户	支付宝当日交易笔数	支付宝日复购用户	支付宝日新购用户	微信当日交易金额	微信当日交易用户	微信当日交易笔数	微信日复购用户	微信日新购用户	赠送份数	购买单份套餐次数	购买5份套餐次数	购买10份套餐次数	购买50份套餐次数	购买100份套餐次数

# 日期
SELECT FROM_UNIXTIME(1603079211703 DIV 1000, '%Y-%m-%d');
SELECT CURRENT_DATE();
select round(UNIX_TIMESTAMP(@report_date)*1000, 0);
SELECT DATE_ADD(@report_date,INTERVAL 1 DAY);

set @report_date = '2021-11-19 00:00:00';
set @report_date_end = DATE_ADD(@report_date,INTERVAL 1 DAY);
set @report_timestamp = round(UNIX_TIMESTAMP(@report_date)*1000, 0);

set @report_timestamp_end = round(UNIX_TIMESTAMP(@report_date_end)*1000, 0);
SELECT @report_date, @report_date_end, @report_timestamp, @report_timestamp_end;

select * from t_order ORDER BY create_time ASC LIMIT 0,1; #2020-10-19
select * from t_contract ORDER BY create_time ASC LIMIT 0,1;


# 当日总收入 - 
# SELECT price, dis_amount, (price-dis_amount) as diff, `status` FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d') = @report_date and `status` = 2;
SELECT SUM(price-dis_amount) as `当日总收入` FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2;
SELECT * from t_order where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2;

# 当日库存份数 - 当天统计过往用户已购买的套餐的待签署数量
# SELECT * from t_bought_package where status = 0 and activity_type = 0 and expired_time >= @report_timestamp_end and create_time <= @report_timestamp_end;
SELECT SUM(count) as `当日库存份数` from t_bought_package where status = 0 and activity_type = 0 and expired_time > @report_timestamp_end and create_time < @report_timestamp_end;
# SELECT SUM(count) as `当日库存份数` from t_bought_package where status = 0 and activity_type = 0 and FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date;

# 当日消耗份数 - 当日套餐抵扣0（划扣） -  撤回3（撤销划扣）
SELECT MAX(name), SUM(amount), package_deduction_kind from t_deduction_record WHERE FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date GROUP BY package_deduction_kind;
SELECT * from t_deduction_record where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and package_deduction_kind=2;
select (t1.amount - t2.amount) as `当日消耗份数` from (SELECT SUM(amount) as amount from t_deduction_record where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and package_deduction_kind=0) t1,
(SELECT SUM(amount) as amount from t_deduction_record where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and package_deduction_kind=3) t2;


# 支付宝当日交易金额	支付宝当日交易用户	支付宝当日交易笔数	支付宝日复购用户	支付宝日新购用户 - 
SELECT SUM(price-dis_amount) as `支付宝当日交易金额` FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=1;
SELECT COUNT(DISTINCT access_id) as `支付宝当日交易用户` FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=1;
SELECT COUNT(1) as `支付宝当日交易笔数` FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=1;
SELECT COUNT(1) as `支付宝日复购用户` from (SELECT COUNT(access_id) as access_ids FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=1 GROUP BY access_id) t where t.access_ids>1;
SELECT COUNT(1) as `支付宝日新购用户` from (SELECT COUNT(access_id) as access_ids FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=1 GROUP BY access_id) t where t.access_ids=1;


# 微信当日交易金额	微信当日交易用户	微信当日交易笔数	微信日复购用户	微信日新购用户 - 
SELECT SUM(price-dis_amount) as `微信当日交易金额` FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=0;
SELECT COUNT(DISTINCT access_id) as `微信当日交易用户` FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=0;
SELECT COUNT(1) as `微信当日交易笔数` FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=0;
# 微信日复购用户 - 当日有两笔以上订单的用户
SELECT COUNT(1) as `微信日复购用户` from (SELECT COUNT(access_id) as access_ids FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=0 GROUP BY access_id) t where t.access_ids>1;
# 微信日新购用户 - 微信平台首次购买的记录（程序实现：先建立首次购买表按平台，在查找是否存在购买）
SELECT COUNT(1) as `微信日新购用户` from (SELECT COUNT(access_id) as access_ids FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=0 GROUP BY access_id) t where t.access_ids=1;

SELECT DISTINCT access_id FROM `t_order` where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and `status` = 2 and pay_method=0 GROUP BY access_id;
SELECT access_id from t_order where access_id = '1695448a34ff4a87aecb563eba13f5c8' and FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') <> @report_date LIMIT 0,1;


# 赠送份数	购买单份套餐次数	购买5份套餐次数	购买10份套餐次数	购买50份套餐次数	购买100份套餐次数
SELECT IFNULL(SUM(amount), 0) as `赠送份数` from t_bought_package where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and activity_type=1;
SELECT count(1) from t_bought_package where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and activity_type=0 and amount=1;
SELECT count(1) from t_bought_package where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and activity_type=0 and amount=5;
SELECT count(1) from t_bought_package where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and activity_type=0 and amount=10;
SELECT count(1) from t_bought_package where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and activity_type=0 and amount=50;
SELECT count(1) from t_bought_package where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and activity_type=0 and amount=100;


-- 业务维度分析
# 日期	当日合同签署总数 当日合同签署次数	自定义类合同当日签署数	自定义类合同占比	模板类合同当日签署份数	模板类合同当日占比	模板类借贷类合同签署份数	模板类借贷类合同当日占比	模板类非借贷类合同签署份数	模板类非借贷类合同当日占比	法律增值业务

# 当日合同签署总数 - 已完成4, 
SELECT count(1) from t_contract_operate_record where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and operate_type=4;
# 当日合同签署次数 - 签署次数1
SELECT count(1) from t_contract_operate_record where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date and operate_type=1;

# 自定义类合同当日签署数	自定义类合同占比 模板类合同当日签署份数	模板类合同当日占比
# template_type= 0自定义类合同 | 1模板类合同
select template_type as `合同类别`, CASE WHEN template_type=0 THEN "自定义类合同当日签署数" ELSE "模板类合同当日签署份数" END as `名称`, SUM(template_count) as `当日签署份数` from  (SELECT template_id, count(1) as template_count, CASE WHEN template_id>0 THEN 1 ELSE 0 END template_type from t_contract where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date GROUP BY template_id) t GROUP BY template_type;

# 模板类借贷类合同签署份数	模板类借贷类合同当日占比	模板类非借贷类合同签署份数	模板类非借贷类合同当日占比
# SELECT id from t_template_class where name like '%借%'; // 3 19
# SELECT id from t_template where template_class_id in (3, 19);
# select count(t.template_count) as `模板类借贷类合同签署份数` from (SELECT template_id, count(1) as template_count from t_contract GROUP BY template_id HAVING template_id>0 and template_id in (SELECT id from t_template where template_class_id in (3, 19))) t;
# select count(t.template_count) as `模板类非借贷类合同签署份数` from (SELECT template_id, count(1) as template_count from t_contract GROUP BY template_id HAVING template_id>0 and template_id not in (SELECT id from t_template where template_class_id in (3, 19))) t;

select t.template_type as `合同类别`, CASE WHEN template_type=2 THEN "模板类借贷类合同签署份数" ELSE "模板类非借贷类合同签署份数" END as `名称`, SUM(t.template_count) as `签署份数` from (SELECT template_id, count(1) as template_count, case when template_id in (SELECT id from t_template where template_class_id in (SELECT id from t_template_class where name like '%借%')) then 2 else 3 end template_type from t_contract where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date GROUP BY template_id HAVING template_id>0) t GROUP BY template_type;

# 法律增值业务
SELECT * from contract_apply_order;
SELECT sum(price-dis_amount) as `总额`, count(1) as `记录数` from contract_apply_order where FROM_UNIXTIME(create_time DIV 1000, '%Y-%m-%d 00:00:00') = @report_date;

