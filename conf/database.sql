/*
Navicat MySQL Data Transfer

Source Server         : 192.168.10.160-kangyi
Source Server Version : 50720
Source Host           : 192.168.10.160:3306
Source Database       : xsy_finance

Target Server Type    : MYSQL
Target Server Version : 50720
File Encoding         : 65001

Date: 2022-10-22 08:54:49
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for address
-- ----------------------------
DROP TABLE IF EXISTS `address`;
CREATE TABLE `address` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source_address` varchar(255) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `trx20_total` int(11) DEFAULT NULL,
  `lable` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `address` (`address`(191)) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4985 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for exchange
-- ----------------------------
DROP TABLE IF EXISTS `exchange`;
CREATE TABLE `exchange` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1359 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for lable
-- ----------------------------
DROP TABLE IF EXISTS `lable`;
CREATE TABLE `lable` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `address` varchar(255) DEFAULT NULL,
  `lable` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `lable` (`lable`(191)) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4985 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tx
-- ----------------------------
DROP TABLE IF EXISTS `tx`;
CREATE TABLE `tx` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `address` varchar(255) DEFAULT NULL,
  `from_address` varchar(255) DEFAULT NULL,
  `to_address` varchar(255) DEFAULT NULL,
  `direction` varchar(3) DEFAULT NULL,
  `token_abbr` varchar(255) DEFAULT NULL,
  `number` varchar(255) DEFAULT NULL,
  `transaction_id` varchar(255) DEFAULT NULL,
  `event_type` varchar(255) DEFAULT NULL,
  `contract_type` varchar(255) DEFAULT NULL,
  `FromAddressIsContract` int(11) DEFAULT NULL,
  `ToAddressIsContract` int(11) DEFAULT NULL,
  `time` int(20) DEFAULT NULL,
  `block` int(20) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `to_address` (`to_address`(191)) USING BTREE,
  KEY `from_address` (`from_address`(191)) USING BTREE,
  KEY `number` (`token_abbr`(191),`number`(191)) USING BTREE,
  KEY `time` (`time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4394148 DEFAULT CHARSET=utf8mb4;
