/*
 Navicat Premium Data Transfer

 Source Server         : 本机
 Source Server Type    : MySQL
 Source Server Version : 50738
 Source Host           : localhost:3306
 Source Schema         : xll-core

 Target Server Type    : MySQL
 Target Server Version : 50738
 File Encoding         : 65001

 Date: 13/09/2023 15:35:43
*/

SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

SET
FOREIGN_KEY_CHECKS = 1;

CREATE TABLE `tb_job_management`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `app_name`   varchar(64) NOT NULL,
    `name`       varchar(64) NOT NULL,
    PRIMARY KEY (`id`),
    KEY          `idx_tb_job_management_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;