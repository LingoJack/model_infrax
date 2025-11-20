参考[规范]，根据[表设计手稿]生成[建表sql]

[sql规范]
```sql
CREATE TABLE IF NOT EXISTS `t_user`
(
    `id`         bigint(20) unsigned                     NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `userId`     varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户ID',
    `userName`   varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名称',
    `createTime` datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime` datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_userId` (`userId`),
    KEY `idx_userId_userName` (`userId`, `userName`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户表';

CREATE TABLE IF NOT EXISTS `t_memory`
(
    `id`               bigint(20) unsigned                     NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `memoryId`         varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '记忆ID',
    `userId`           varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户ID',
    `summary`          text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '摘要',
    `raw_dialog_array` text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '对话的JSON格式，包括了节点的输入输出等',
    `createTime`       datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime`       datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_memoryId` (`memoryId`),
    KEY `idx_userId` (`userId`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='记忆表';

CREATE TABLE IF NOT EXISTS `t_llm_history`
(
    `id`         bigint(20) unsigned                     NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `model`      varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模型名称',
    `input`      text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '输入内容',
    `output`     text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '输出内容',
    `createTime` datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime` datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_model_createTime` (`model`, `createTime`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='LLM历史记录表';

CREATE TABLE IF NOT EXISTS `t_work_node`
(
    `id`           bigint(20) unsigned                     NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `nodeId`       varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '节点ID',
    `nodeName`     varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '节点名称',
    `executor`     varchar(50) COLLATE utf8mb4_unicode_ci  NOT NULL COMMENT '执行器: robot / brain',
    `type`         varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '节点类型',
    `description`  text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '描述',
    `inputSchema`  text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '输入模式',
    `outputSchema` text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '输出模式',
    `timeout`      int(11)                                NOT NULL DEFAULT '-1' COMMENT '超时时间(秒)，-1表示无超时',
    `extra`        text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '附加信息，是map的json表达',
    `createTime`   datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updateTime`   datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_nodeId` (`nodeId`),
    KEY `idx_nodeId_nodeName` (`nodeId`, `nodeName`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='工作节点表';
```

[表设计手稿]
```

```