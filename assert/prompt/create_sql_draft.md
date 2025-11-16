参考[规范]，根据[表设计手稿]生成[建表sql]

[sql规范]
```sql
CREATE TABLE IF NOT EXISTS `t_feedback_task`
(
    `id`                  bigint(20) unsigned                     NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `feedbackTaskId`      varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务ID',
    `appId`               varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '应用ID',
    `appName`             varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '应用名称',
    `status`              varchar(50) COLLATE utf8mb4_unicode_ci  NOT NULL DEFAULT 'init' COMMENT '修复状态: success, running, fail, cancel, init, confirm',
    `apiName`             varchar(256) COLLATE utf8mb4_unicode_ci          DEFAULT NULL COMMENT 'API名称',
    `isNeedFix`           tinyint(1)                              NOT NULL DEFAULT '0' COMMENT '是否需要修复: 0-否 1-是',
    `fixPart`             varchar(10) COLLATE utf8mb4_unicode_ci           DEFAULT NULL COMMENT '需要修复的部分: backend, frontend, all',
    `reason`              text COLLATE utf8mb4_unicode_ci COMMENT '错误原因及修复说明',
    `conciseReason`       varchar(128) COLLATE utf8mb4_unicode_ci          DEFAULT NULL COMMENT '错误原因及修复说明(简洁)',
    `beforeCode`          text COLLATE utf8mb4_unicode_ci COMMENT '修复前代码',
    `afterCode`           text COLLATE utf8mb4_unicode_ci COMMENT '修复后代码',
    `beforeMod`           text COLLATE utf8mb4_unicode_ci COMMENT '修复前模块',
    `afterMod`            text COLLATE utf8mb4_unicode_ci COMMENT '修复后模块',
    `errContext`          text COLLATE utf8mb4_unicode_ci COMMENT '错误上下文',
    `hashCode`            varchar(128) COLLATE utf8mb4_unicode_ci          DEFAULT NULL COMMENT 'hash值',
    `route`               varchar(256) COLLATE utf8mb4_unicode_ci          DEFAULT '' COMMENT 'FIP协议描述的前端页面路由',
    `feedbackSource`      varchar(20) COLLATE utf8mb4_unicode_ci           DEFAULT 'system' COMMENT '反馈来源: system(兼容旧逻辑，为空默认system), user',
    `userFeedbackContent` text COLLATE utf8mb4_unicode_ci COMMENT '用户反馈内容，只有feedbackSource为user时，该字段才有效',
    `taskGroupId`         varchar(128) COLLATE utf8mb4_unicode_ci          DEFAULT NULL COMMENT '任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同',
    `createTime`          datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime`          datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_feedbackTaskId` (`feedbackTaskId`),
    KEY `idx_appId` (`appId`),
    KEY `idx_taskGroupId` (`taskGroupId`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='反馈任务表';
```

[表设计手稿]
```


```