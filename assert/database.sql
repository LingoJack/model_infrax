CREATE TABLE IF NOT EXISTS `t_user`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `userId`     varchar(128)        NOT NULL DEFAULT '' COMMENT '用户ID',
    `userName`   varchar(128)        NOT NULL DEFAULT '' COMMENT '用户名称',
    `createTime` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_userId_userName` (`userId`, `userName`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT = '用户表';