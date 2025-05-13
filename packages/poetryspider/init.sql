DROP TABLE IF EXISTS `poetry`;
CREATE TABLE `poetry` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增唯一ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `poetry_id` varchar(255) DEFAULT NULL COMMENT '诗文ID',
    `title` varchar(255) DEFAULT NULL COMMENT '名称',
    `author` varchar(100) DEFAULT NULL COMMENT '作者',
    `dynasty` varchar(100) DEFAULT NULL COMMENT '朝代',
    `poetry_type` varchar(255) DEFAULT NULL COMMENT '类型',
    `paragraphs` text DEFAULT NULL COMMENT '主题',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_poetry_id` (`poetry_id`),
    KEY `idx_author` (`author`),
    KEY `idx_title` (`title`)
) AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '诗词表';

DROP TABLE IF EXISTS `interpret`;
CREATE TABLE `interpret` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增唯一ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `poetry_id` varchar(255) DEFAULT NULL COMMENT '诗文ID',
    `translation` text DEFAULT NULL COMMENT '译文',
    `annotation` text DEFAULT NULL COMMENT '注释',
    `intro` text DEFAULT NULL COMMENT '评价',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_poetry_id` (`poetry_id`),
    KEY `idx_poetry_id` (`poetry_id`)
) AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '诗词表';
