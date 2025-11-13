package parser

import (
	"fmt"
	"model_infrax/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Parser SQL解析器，用于解析数据库表结构
type Parser struct {
	db     *gorm.DB
	config *config.Configger
}

const (
	// 获取所有表的描述
	sqlShowTableStatus = "show table status"

	// 获取单个表的描述
	sqlShowTableStatusForSingleTable = "show table status like '%s'"

	// 查看所有字段
	sqlShowTableFields = "show full fields from %s"

	// 查看所有索引
	sqlShowTableIndexes = "show full indexes from %s"
)

func NewParser(cfg *config.Configger) (*Parser, error) {
	// 构建数据库连接DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.GenerateConfig.Username,
		cfg.GenerateConfig.Password,
		cfg.GenerateConfig.Host,
		cfg.GenerateConfig.Port,
		cfg.GenerateConfig.DatabaseName)

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	return &Parser{
		db:     db,
		config: cfg,
	}, nil
}

// AllTables 获取所有表名
func (p *Parser) AllTables() ([]string, error) {
	var tables []string
	// TODO: 实现获取所有表的逻辑
	return tables, nil
}
