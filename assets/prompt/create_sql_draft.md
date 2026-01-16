参考[规范] ，根据[表设计手稿] 生成[建表sql] 

[sql规范] 
```sql
CREATE TABLE IF NOT EXISTS `t_robot_instance`
(
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `robot_id`            varchar(128)        NOT NULL DEFAULT '' COMMENT '机器人ID',
    `robot_name`          varchar(256)        NOT NULL DEFAULT '' COMMENT '机器人名称',
    `robot_key`           varchar(256)        NOT NULL DEFAULT '' COMMENT '机器人密钥',
    `status`              varchar(20)         NOT NULL DEFAULT 'enable' COMMENT '状态: enable-启用, disable-停用',
    `owner_uin`           varchar(128)        NOT NULL DEFAULT '' COMMENT '主账号',
    `uin`                 varchar(128)        NOT NULL DEFAULT '' COMMENT 'UIN',
    `robot_model_id`      varchar(128)        NOT NULL DEFAULT '' COMMENT '机器人型号ID',
    `robot_model_name`    varchar(128)        NOT NULL DEFAULT '' COMMENT '机器人型号名称',
    `robot_model_version` varchar(128)        NOT NULL DEFAULT '' COMMENT '机器人型号版本',
    `create_time`         datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`         datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_robot_id` (`robot_id`),
    KEY `idx_robot_model_id` (`robot_model_id`),
    KEY `idx_owner_uin_uin` (`owner_uin`, `uin`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='机器人实例表';
```

[表设计手稿] 
```

```