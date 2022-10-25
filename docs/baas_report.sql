/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.90.146_3306
 Source Server Type    : MySQL
 Source Server Version : 80031
 Source Host           : 192.168.90.146:3306
 Source Schema         : baas_report

 Target Server Type    : MySQL
 Target Server Version : 80031
 File Encoding         : 65001

 Date: 25/10/2022 18:21:31
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_basedata_statistics
-- ----------------------------
DROP TABLE IF EXISTS `t_basedata_statistics`;
CREATE TABLE `t_basedata_statistics`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `date_time` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '日期',
  `order_value_total` int NULL DEFAULT NULL COMMENT '订单总金额',
  `order_total` int NULL DEFAULT NULL COMMENT '订单总数',
  `sign_flow_record` int NULL DEFAULT NULL COMMENT '当日签署总数',
  `user_org_auth` int NULL DEFAULT NULL COMMENT '认证企业用户总数',
  `user_org_cancel` int NULL DEFAULT NULL COMMENT '注销企业用户总数',
  `user_person_auth` int NULL DEFAULT NULL COMMENT '认证个人用户总数',
  `user_person_cancel` int NULL DEFAULT NULL COMMENT '注销个人用户总数',
  PRIMARY KEY (`id`, `date_time`) USING BTREE,
  UNIQUE INDEX `date_index`(`date_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 826 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for t_report_business
-- ----------------------------
DROP TABLE IF EXISTS `t_report_business`;
CREATE TABLE `t_report_business`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `currentDate` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '日期',
  `signSuccessTotalCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '当日合同签署总数',
  `signTotalCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '当日合同签署次数',
  `signSuccessPresentTotalCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '当日体验合同签署总数',
  `signPresentTotalCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '当日体验合同签署次数',
  `customContractSignCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '自定义类合同当日签署数',
  `customContractSignPercent` float UNSIGNED NULL DEFAULT NULL COMMENT '自定义类合同占比',
  `templateContractSignCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '模板类合同当日签署份数',
  `templateContractSignPercent` float UNSIGNED NULL DEFAULT NULL COMMENT '模板类合同当日占比',
  `templateLoanContractSignCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '模板类借贷类合同签署份数',
  `templateLoanContractSignPercent` float UNSIGNED NULL DEFAULT NULL COMMENT '模板类借贷类合同当日占比',
  `templateNoLoanContractSignCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '模板类非借贷类合同签署份数',
  `templateNoLoanContractSignPercent` float UNSIGNED NULL DEFAULT NULL COMMENT '模板类非借贷类合同当日占比',
  `createTime` timestamp(0) NULL DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`, `currentDate`) USING BTREE,
  INDEX `currentDate`(`currentDate`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2217 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = '业务维度分析' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for t_report_dict
-- ----------------------------
DROP TABLE IF EXISTS `t_report_dict`;
CREATE TABLE `t_report_dict`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `item` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT 'key 主键',
  `value` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT 'key => value',
  `create_time` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT 'operator log time',
  `update_time` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT 'operator log update time',
  PRIMARY KEY (`id`, `item`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = '定义一些报表需要的字典，如lastReportDate etc' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for t_report_revenue
-- ----------------------------
DROP TABLE IF EXISTS `t_report_revenue`;
CREATE TABLE `t_report_revenue`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'auto seq',
  `currentDate` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '日期',
  `totalAmount` float UNSIGNED NULL DEFAULT NULL COMMENT '当日总收入',
  `stockCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '当天统计过往待签署数量(购买的+赠送的)',
  `expendCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '当日消耗份数',
  `expiredPurchaseCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '当日过期份数(已购)',
  `expiredPresentCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '当日过期份数(赠送)',
  `wechatTransAmount` float UNSIGNED NULL DEFAULT NULL COMMENT '微信当日交易金额',
  `wechatTransAccess` bigint UNSIGNED NULL DEFAULT NULL COMMENT '微信当日交易用户',
  `wechatTransCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '微信当日交易笔数',
  `wechatRepurchaseAccess` bigint UNSIGNED NULL DEFAULT NULL COMMENT '微信日复购用户',
  `wechatFirstPurchaseAccess` bigint UNSIGNED NULL DEFAULT NULL COMMENT '微信日新购用户',
  `alipayTransAmount` float UNSIGNED NULL DEFAULT NULL COMMENT '支付宝当日交易金额',
  `alipayTransAccess` bigint UNSIGNED NULL DEFAULT NULL COMMENT '支付宝当日交易用户',
  `alipayTransCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '支付宝当日交易笔数',
  `alipayRepurchaseAccess` bigint UNSIGNED NULL DEFAULT NULL COMMENT '支付宝日复购用户',
  `alipayFirstPurchaseAccess` bigint UNSIGNED NULL DEFAULT NULL COMMENT '支付宝日新购用户',
  `presentCount` bigint UNSIGNED NULL DEFAULT NULL COMMENT '赠送份数',
  `purchasePackageAmount1` bigint UNSIGNED NULL DEFAULT NULL COMMENT '购买单份套餐次数',
  `purchasePackageAmount5` bigint UNSIGNED NULL DEFAULT NULL COMMENT '购买5份套餐次数',
  `purchasePackageAmount10` bigint UNSIGNED NULL DEFAULT NULL COMMENT '购买10份套餐次数',
  `purchasePackageAmount50` bigint UNSIGNED NULL DEFAULT NULL COMMENT '购买50份套餐次数',
  `purchasePackageAmount100` bigint UNSIGNED NULL DEFAULT NULL COMMENT '购买100份套餐次数',
  `createTime` timestamp(0) NULL DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`, `currentDate`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2217 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = '营收维度分析' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for t_user_statistics
-- ----------------------------
DROP TABLE IF EXISTS `t_user_statistics`;
CREATE TABLE `t_user_statistics`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `date_time` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '日期',
  `user_total` int NULL DEFAULT NULL COMMENT '累计用户总数',
  `user_person_add` int NULL DEFAULT NULL COMMENT '新增个人用户数',
  `user_active_total` int NULL DEFAULT NULL COMMENT '日活用户总数',
  `user_visit_total` int NULL DEFAULT NULL COMMENT '日访问用户总数',
  `org_active_day` int NULL DEFAULT NULL COMMENT '企业日活跃用户',
  `org_active_week` int NULL DEFAULT NULL COMMENT '企业7天活跃用户',
  `org_active_15` int NULL DEFAULT NULL COMMENT '企业15天活跃用户',
  `org_active_30` int NULL DEFAULT NULL COMMENT '企业30天活跃用户',
  `org_active_90` int NULL DEFAULT NULL COMMENT '企业90天活跃用户',
  `org_silent_90` int NULL DEFAULT NULL COMMENT '企业90天沉默用户',
  `person_active_day` int NULL DEFAULT NULL COMMENT '个人日活跃用户',
  `person_active_week` int NULL DEFAULT NULL COMMENT '个人7天活跃用户',
  `person_active_15` int NULL DEFAULT NULL COMMENT '个人15天活跃用户',
  `person_active_30` int NULL DEFAULT NULL COMMENT '个人30天活跃用户',
  `person_active_90` int NULL DEFAULT NULL COMMENT '个人90天活跃用户',
  `person_silen_90` int NULL DEFAULT NULL COMMENT '个人90天沉默用户',
  `person_resign_7` int NULL DEFAULT NULL COMMENT '个人7天复签用户',
  `person_resign_15` int NULL DEFAULT NULL COMMENT '个人15天复签用户',
  `person_resign_30` int NULL DEFAULT NULL COMMENT '个人30天复签用户',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `date_index`(`date_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 827 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
