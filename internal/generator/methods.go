package generator

import (
	"fmt"
	"strings"

	"github.com/LingoJack/model_infrax/internal/conf"
)

// 配置键
const (
	cfgGormTagStyle = "generate_option.gorm_tag_style"
	cfgCommentStyle = "generate_option.comment_style"
	cfgDaoMethods   = "generate_option.dao_methods"
)

// 风格合法值
var (
	ValidGormTagStyles = []string{"full", "standard", "minimal"}
	ValidCommentStyles = []string{"full", "brief", "none"}
)

// MethodMeta 描述一个可配置的 DAO 方法/方法模式
type MethodMeta struct {
	ID    string `json:"id"`    // 配置用方法名（如 SelectByPk）
	Desc  string `json:"desc"`  // 中文说明
	Group string `json:"group"` // 所属分区
}

// AllDaoMethods 全部可配置方法清单（dao_methods 白名单的合法值）
// Pk 系列按主键实例化（SelectByPk → SelectById），Index 系列按唯一/普通索引列实例化
var AllDaoMethods = []MethodMeta{
	{"WithTx", "事务副本（开启新事务连接）", "事务"},
	{"Transaction", "嵌套事务回调", "事务"},
	{"SelectList", "分页查询列表", "列表查询"},
	{"SelectCount", "查询总数", "列表查询"},
	{"SelectListWithAppendConditionFunc", "追加自定义条件的分页查询", "列表查询"},
	{"SelectCountWithAppendConditionFunc", "追加自定义条件的计数", "列表查询"},
	{"Insert", "插入单条", "写入"},
	{"InsertBatch", "批量插入", "写入"},
	{"InsertOrUpdateNullable", "插入或更新（零值覆盖）", "写入"},
	{"InsertOrUpdateBatchNullable", "批量插入或更新（零值覆盖）", "写入"},
	{"ExecSql", "原生 SQL 查询", "原生 SQL"},
	{"SelectByPk", "按主键查单条（如 SelectById）", "主键方法"},
	{"SelectByPkList", "按主键列表批量查询", "主键方法"},
	{"UpdateByPk", "按主键更新（非零值）", "主键方法"},
	{"UpdateByPkWithMap", "按主键 Map 更新（可零值覆盖）", "主键方法"},
	{"UpdateByPkWithCondition", "按主键+条件更新", "主键方法"},
	{"UpdateByPkWithMapAndCondition", "按主键 Map+条件更新", "主键方法"},
	{"DeleteByPk", "按主键删除", "主键方法"},
	{"SelectByIndex", "按唯一/普通索引查单条（如 SelectByUserId）", "索引方法"},
	{"SelectByIndexList", "按索引批量查询（仅单列索引）", "索引方法"},
	{"UpdateByIndex", "按索引更新（非零值）", "索引方法"},
	{"UpdateByIndexWithMap", "按索引 Map 更新", "索引方法"},
	{"UpdateByIndexWithCondition", "按索引+条件更新", "索引方法"},
	{"UpdateByIndexWithMapAndCondition", "按索引 Map+条件更新", "索引方法"},
	{"DeleteByIndex", "按索引删除", "索引方法"},
	{"PoBuilder", "的 Builder 链式 setter（WithXxx 系列）", "PO 附加"},
	{"PoJsonify", "的 Jsonify/JsonifyIndent 序列化方法", "PO 附加"},
}

// AllDaoMethodIDs 返回全部可配置方法名（供 flag 校验等使用）
func AllDaoMethodIDs() []string {
	ids := make([]string, 0, len(AllDaoMethods))
	for _, m := range AllDaoMethods {
		ids = append(ids, m.ID)
	}
	return ids
}

// MethodSwitches DAO/PO 生成内容的方法开关
// 前半部分与 AllDaoMethods 一一对应；后两个为内部联动推导结果（不可直接配置）
type MethodSwitches struct {
	WithTx                             bool
	Transaction                        bool
	SelectList                         bool
	SelectCount                        bool
	SelectListWithAppendConditionFunc  bool
	SelectCountWithAppendConditionFunc bool
	Insert                             bool
	InsertBatch                        bool
	InsertOrUpdateNullable             bool
	InsertOrUpdateBatchNullable        bool
	ExecSql                            bool
	SelectByPk                         bool
	SelectByPkList                     bool
	UpdateByPk                         bool
	UpdateByPkWithMap                  bool
	UpdateByPkWithCondition            bool
	UpdateByPkWithMapAndCondition      bool
	DeleteByPk                         bool
	SelectByIndex                      bool
	SelectByIndexList                  bool
	UpdateByIndex                      bool
	UpdateByIndexWithMap               bool
	UpdateByIndexWithCondition         bool
	UpdateByIndexWithMapAndCondition   bool
	DeleteByIndex                      bool
	PoBuilder                          bool
	PoJsonify                          bool

	// 内部联动推导（不可配置）
	BuildQueryCondition bool // 被 4 个 Select 系列依赖
	ValidOrderBy        bool // isValidOrderBy/getValidOrderByFields，被 SelectList/SelectCount 依赖
}

// NewMethodSwitches 从配置解析方法开关
// dao_methods 未配置（或为空）时全部生成；配置后按白名单开启
// 必须在 conf.InitWithPath 之后调用
func NewMethodSwitches() (MethodSwitches, error) {
	sw := MethodSwitches{}

	names, err := conf.ValueStrSlice(cfgDaoMethods)
	if err != nil || len(names) == 0 {
		// 未配置或空列表 = 全部生成
		sw = allOn()
		sw.derive()
		return sw, nil
	}

	valid := map[string]bool{}
	for _, m := range AllDaoMethods {
		valid[m.ID] = true
	}
	set := map[string]bool{}
	for _, name := range names {
		name = strings.TrimSpace(name)
		if !valid[name] {
			return sw, fmt.Errorf("配置 %s 含非法方法名: %s，可用值: %s", cfgDaoMethods, name, strings.Join(AllDaoMethodIDs(), ", "))
		}
		set[name] = true
	}

	sw.WithTx = set["WithTx"]
	sw.Transaction = set["Transaction"]
	sw.SelectList = set["SelectList"]
	sw.SelectCount = set["SelectCount"]
	sw.SelectListWithAppendConditionFunc = set["SelectListWithAppendConditionFunc"]
	sw.SelectCountWithAppendConditionFunc = set["SelectCountWithAppendConditionFunc"]
	sw.Insert = set["Insert"]
	sw.InsertBatch = set["InsertBatch"]
	sw.InsertOrUpdateNullable = set["InsertOrUpdateNullable"]
	sw.InsertOrUpdateBatchNullable = set["InsertOrUpdateBatchNullable"]
	sw.ExecSql = set["ExecSql"]
	sw.SelectByPk = set["SelectByPk"]
	sw.SelectByPkList = set["SelectByPkList"]
	sw.UpdateByPk = set["UpdateByPk"]
	sw.UpdateByPkWithMap = set["UpdateByPkWithMap"]
	sw.UpdateByPkWithCondition = set["UpdateByPkWithCondition"]
	sw.UpdateByPkWithMapAndCondition = set["UpdateByPkWithMapAndCondition"]
	sw.DeleteByPk = set["DeleteByPk"]
	sw.SelectByIndex = set["SelectByIndex"]
	sw.SelectByIndexList = set["SelectByIndexList"]
	sw.UpdateByIndex = set["UpdateByIndex"]
	sw.UpdateByIndexWithMap = set["UpdateByIndexWithMap"]
	sw.UpdateByIndexWithCondition = set["UpdateByIndexWithCondition"]
	sw.UpdateByIndexWithMapAndCondition = set["UpdateByIndexWithMapAndCondition"]
	sw.DeleteByIndex = set["DeleteByIndex"]
	sw.PoBuilder = set["PoBuilder"]
	sw.PoJsonify = set["PoJsonify"]

	sw.derive()
	return sw, nil
}

// allOn 返回全开开关
func allOn() MethodSwitches {
	return MethodSwitches{
		WithTx:                             true,
		Transaction:                        true,
		SelectList:                         true,
		SelectCount:                        true,
		SelectListWithAppendConditionFunc:  true,
		SelectCountWithAppendConditionFunc: true,
		Insert:                             true,
		InsertBatch:                        true,
		InsertOrUpdateNullable:             true,
		InsertOrUpdateBatchNullable:        true,
		ExecSql:                            true,
		SelectByPk:                         true,
		SelectByPkList:                     true,
		UpdateByPk:                         true,
		UpdateByPkWithMap:                  true,
		UpdateByPkWithCondition:            true,
		UpdateByPkWithMapAndCondition:      true,
		DeleteByPk:                         true,
		SelectByIndex:                      true,
		SelectByIndexList:                  true,
		UpdateByIndex:                      true,
		UpdateByIndexWithMap:               true,
		UpdateByIndexWithCondition:         true,
		UpdateByIndexWithMapAndCondition:   true,
		DeleteByIndex:                      true,
		PoBuilder:                          true,
		PoJsonify:                          true,
	}
}

// derive 推导内部联动开关
func (s *MethodSwitches) derive() {
	s.BuildQueryCondition = s.SelectList || s.SelectCount ||
		s.SelectListWithAppendConditionFunc || s.SelectCountWithAppendConditionFunc
	s.ValidOrderBy = s.SelectList || s.SelectCount
}

// GormTagStyle 获取 gorm tag 风格配置（缺省 full）
func GormTagStyle() string {
	style := conf.ValueStr(cfgGormTagStyle)
	if style == "" {
		return "full"
	}
	return style
}

// CommentStyle 获取注释风格配置（缺省 full）
func CommentStyle() string {
	style := conf.ValueStr(cfgCommentStyle)
	if style == "" {
		return "full"
	}
	return style
}
