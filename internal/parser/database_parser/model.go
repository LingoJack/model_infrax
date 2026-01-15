package database_parser

type mysqlTable struct {
	Name          string  `json:"ColumnName"`      // 表名
	Engine        string  `json:"Engine"`          // 存储引擎
	Version       int     `json:"Version"`         // 版本号
	RowFormat     string  `json:"Row_format"`      // 行格式
	Rows          int64   `json:"Rows"`            // 行数（估算值）
	AvgRowLength  int64   `json:"Avg_row_length"`  // 平均行长度
	DataLength    int64   `json:"Data_length"`     // 数据长度
	MaxDataLength int64   `json:"Max_data_length"` // 最大数据长度
	IndexLength   int64   `json:"Index_length"`    // 索引长度
	DataFree      int64   `json:"Data_free"`       // 空闲空间
	AutoIncrement *int64  `json:"Auto_increment"`  // 自增值（可能为null）
	CreateTime    *string `json:"Create_time"`     // 创建时间（可能为null）
	UpdateTime    *string `json:"Update_time"`     // 更新时间（可能为null）
	CheckTime     *string `json:"Check_time"`      // 检查时间（可能为null）
	Collation     string  `json:"Collation"`       // 字符集校对规则
	Checksum      *int64  `json:"Checksum"`        // 校验和（可能为null）
	CreateOptions string  `json:"Create_options"`  // 创建选项
	Comment       string  `json:"Comment"`         // 表注释
}

type mysqlField struct {
	Field      string  `json:"Field"`      // 字段名
	Type       string  `json:"Type"`       // 字段类型
	Collation  *string `json:"Collation"`  // 字符集校对规则（可能为null）
	Null       string  `json:"Null"`       // 是否允许为NULL（YES/NO）
	Key        string  `json:"Key"`        // 键类型（PRI/UNI/MUL等）
	Default    *string `json:"Default"`    // 默认值（可能为null）
	Extra      string  `json:"Extra"`      // 额外信息（如auto_increment）
	Privileges string  `json:"Privileges"` // 权限信息
	Comment    string  `json:"Comment"`    // 字段注释
}

type mysqlIndex struct {
	Table        string  `json:"Table"`         // 表名
	NonUnique    int     `json:"Non_unique"`    // 是否非唯一索引（0=唯一索引，1=非唯一索引）
	KeyName      string  `json:"Key_name"`      // 索引名称
	SeqInIndex   int     `json:"Seq_in_index"`  // 字段在索引中的序号（从1开始）
	ColumnName   string  `json:"Column_name"`   // 列名
	Collation    *string `json:"Collation"`     // 排序方式（A=升序，D=降序，NULL=未排序）
	Cardinality  *int64  `json:"Cardinality"`   // 索引中唯一值的数量估算（可能为null）
	SubPart      *int    `json:"Sub_part"`      // 索引前缀长度（可能为null）
	Packed       *string `json:"Packed"`        // 关键字如何被压缩（可能为null）
	Null         string  `json:"Null"`          // 列是否可以包含NULL值
	IndexType    string  `json:"Index_type"`    // 索引类型（BTREE/HASH/FULLTEXT/SPATIAL）
	Comment      string  `json:"Comment"`       // 索引注释
	IndexComment string  `json:"Index_comment"` // 索引注释（创建索引时的COMMENT）
	Visible      string  `json:"Visible"`       // 索引是否可见（YES/NO）
	Expression   *string `json:"Expression"`    // 表达式索引的表达式（可能为null）
}
