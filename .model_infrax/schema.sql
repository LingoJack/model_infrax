-- 在此编写你的 CREATE TABLE 语句
-- 示例:
CREATE TABLE IF NOT EXISTS `t_example`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name`       varchar(128)        NOT NULL COMMENT '名称',
    `createTime` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '示例表';
