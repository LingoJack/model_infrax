package db_infra

import (
	"context"
	"fmt"

	"github.com/LingoJack/model_infrax/internal/conf"
	"github.com/LingoJack/model_infrax/internal/constant"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	if conf.ValueStr("generate_config.generate_mode") == constant.GenerateModeDatabase {
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

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			// TODO
			panic(err)
		}
	}
}

func ExecuteSql(ctx context.Context, sql string, recvPtr interface{}) error {
	return db.WithContext(ctx).Raw(sql).Scan(recvPtr).Error
}
