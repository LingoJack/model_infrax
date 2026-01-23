package db_infra

import (
	"context"
	"fmt"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/constant"
	"github.com/LingoJack/model_infrax/internal/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// Init 初始化数据库连接
// 必须在 conf.InitWithPath 之后调用
// 仅在 generate_mode 为 database 时才会初始化数据库连接
func Init() error {
	generateMode := conf.ValueStr("generate_config.generate_mode")
	if generateMode == constant.GenerateModeDatabase {
		username := conf.ValueStr("generate_config.username")
		password := conf.ValueStr("generate_config.password")
		host := conf.ValueStr("generate_config.host")
		port := conf.ValueInt("generate_config.port")
		databaseName := conf.ValueStr("generate_config.database_name")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			username,
			password,
			host,
			port,
			databaseName,
		)

		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("[Init] 数据库连接失败, dsn=%s:%d/%s, err=%w", host, port, databaseName, err)
		}
		logger.Infof("[Init] 数据库连接成功, host=%s, port=%d, database=%s", host, port, databaseName)
	}
	return nil
}

func ExecSql(ctx context.Context, recvPtr any, sql string, args ...any) error {
	return db.WithContext(ctx).Raw(sql, args).Scan(recvPtr).Error
}
