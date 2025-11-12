CREATE TABLE IF NOT EXISTS `t_artifact`
(
    `id`           bigint(20) unsigned                     NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `artifactId`   varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产物ID',
    `artifactName` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产物名称',
    `sessionId`    varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属的会话',
    `step`         int(11)                                NOT NULL COMMENT '大的步骤点',
    `subStep`      varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '小的步骤点',
    `content`      text COLLATE utf8mb4_unicode_ci COMMENT '内容',
    `version`      varchar(128) COLLATE utf8mb4_unicode_ci          DEFAULT NULL COMMENT '版本',
    `createTime`   datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime`   datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_artifactId` (`artifactId`),
    KEY `idx_sessionId` (`sessionId`),
    KEY `idx_step` (`step`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='任务执行流程中生成的中间产物表';