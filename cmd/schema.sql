CREATE TABLE IF NOT EXISTS `t_session_step` (
                                                `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                                `sessionId` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '会话ID',
                                                `query` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '查询内容',
                                                `think` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '思考内容',
                                                `answer` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '回答内容',
                                                `type` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类型',
                                                `step` int(11) NOT NULL COMMENT '步骤',
                                                `subStep` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '子步骤',
                                                `modelName` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模型名称',
                                                `extra` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '附加信息',
                                                `createTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                                `updateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                                PRIMARY KEY (`id`),
                                                UNIQUE KEY `uk_sessionId` (`sessionId`),
                                                KEY `idx_step_subStep` (`step`, `subStep`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '会话表';