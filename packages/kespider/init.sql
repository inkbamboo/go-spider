DROP TABLE IF EXISTS `area`;
CREATE TABLE `area` (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `district_id` varchar(255) DEFAULT NULL COMMENT '板块ID',
  `district_name` datetime(0) NULL DEFAULT NULL COMMENT '板块名称',
  `area_id` varchar(255) DEFAULT NULL COMMENT '区域ID',
  `area_name` datetime(0) NULL DEFAULT NULL COMMENT '区域名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_district_id` (`district_id`)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '区域表';

DROP TABLE IF EXISTS `house`;
CREATE TABLE `house` (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `housedel_id` varchar(255) DEFAULT NULL COMMENT '房屋ID',
  `district_id` varchar(255) DEFAULT NULL COMMENT '板块ID',
  `house_area` varchar(255) DEFAULT NULL COMMENT '房屋面积',
  `house_orientation` varchar(255) DEFAULT NULL COMMENT '房屋朝向',
  `house_type` varchar(255) DEFAULT NULL COMMENT '房屋类型',
  `house_year` varchar(255) DEFAULT NULL COMMENT '房屋年限',
  `xiaoqu_name` varchar(255) DEFAULT NULL COMMENT '小区名称',
  `house_floor` varchar(255) DEFAULT NULL COMMENT '楼层总高度',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_housedel_id` (`housedel_id`)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '房屋信息表';

DROP TABLE IF EXISTS `house_price`;
CREATE TABLE `house_price` (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `housedel_id` varchar(255) DEFAULT NULL COMMENT '房屋ID',
  `version` varchar(255) DEFAULT NULL COMMENT '版本(更新日期)',
  `district_id` varchar(255) DEFAULT NULL COMMENT '板块ID',
  `total_price` int(0) NULL DEFAULT NULL COMMENT '总价',
  `unit_price` int(0) NULL DEFAULT NULL COMMENT '单价',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_housedel_id_version` (`housedel_id`,`version`)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '价格信息';
