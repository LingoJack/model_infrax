package model

import (
	"encoding/json"
	"time"
)

// TArtifact 任务执行流程中生成的中间产物表
type TArtifact struct {
	Id           uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement;comment:主键ID;not null" json:"id"`
	ArtifactId   string    `gorm:"column:artifactId;type:varchar(128);comment:产物ID;not null" json:"artifactId"`
	ArtifactName string    `gorm:"column:artifactName;type:varchar(128);comment:产物名称;not null" json:"artifactName"`
	SessionId    string    `gorm:"column:sessionId;type:varchar(128);comment:所属的会话;not null" json:"sessionId"`
	Step         int       `gorm:"column:step;type:int;comment:大的步骤点;not null" json:"step"`
	SubStep      string    `gorm:"column:subStep;type:varchar(256);comment:小的步骤点;not null" json:"subStep"`
	Content      *string   `gorm:"column:content;type:text;comment:内容;" json:"content"`
	Version      *string   `gorm:"column:version;type:varchar(128);comment:版本;" json:"version"`
	CreateTime   time.Time `gorm:"column:createTime;type:datetime;comment:创建时间;not null" json:"createTime"`
	UpdateTime   time.Time `gorm:"column:updateTime;type:datetime;comment:更新时间;not null" json:"updateTime"`
}

// TableName 返回表名
func (t *TArtifact) TableName() string {
	return "t_artifact"
}

// Jsonify 将结构体序列化为 JSON 字符串（紧凑格式）
// 返回:
//   - string: JSON 字符串，如果序列化失败则返回错误信息的 JSON
func (t *TArtifact) Jsonify() string {
	byts, err := json.Marshal(t)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(byts)
}

// JsonifyIndent 将结构体序列化为格式化的 JSON 字符串（带缩进）
// 返回:
//   - string: 格式化的 JSON 字符串，如果序列化失败则返回错误信息的 JSON
func (t *TArtifact) JsonifyIndent() string {
	byts, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(byts)
}

// TArtifactBuilder 用于构建 TArtifact 实例的 Builder
type TArtifactBuilder struct {
	instance *TArtifact
}

// NewTArtifactBuilder 创建一个新的 TArtifactBuilder 实例
// 返回:
//   - *TArtifactBuilder: Builder 实例，用于链式调用
func NewTArtifactBuilder() *TArtifactBuilder {
	return &TArtifactBuilder{
		instance: &TArtifact{},
	}
}

// WithArtifactId 设置 artifactId 字段
// 参数:
//   - artifactId: 产物ID
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithArtifactId(artifactId string) *TArtifactBuilder {
	b.instance.ArtifactId = artifactId
	return b
}

// WithArtifactName 设置 artifactName 字段
// 参数:
//   - artifactName: 产物名称
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithArtifactName(artifactName string) *TArtifactBuilder {
	b.instance.ArtifactName = artifactName
	return b
}

// WithSessionId 设置 sessionId 字段
// 参数:
//   - sessionId: 所属的会话
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithSessionId(sessionId string) *TArtifactBuilder {
	b.instance.SessionId = sessionId
	return b
}

// WithStep 设置 step 字段
// 参数:
//   - step: 大的步骤点
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithStep(step int) *TArtifactBuilder {
	b.instance.Step = step
	return b
}

// WithSubStep 设置 subStep 字段
// 参数:
//   - subStep: 小的步骤点
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithSubStep(subStep string) *TArtifactBuilder {
	b.instance.SubStep = subStep
	return b
}

// WithContent 设置 content 字段
// 参数:
//   - content: 内容
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithContent(content *string) *TArtifactBuilder {
	b.instance.Content = content
	return b
}

// WithContentValue 设置 content 字段（便捷方法，自动转换为指针）
// 参数:
//   - content: 内容
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithContentValue(content string) *TArtifactBuilder {
	b.instance.Content = &content
	return b
}

// WithVersion 设置 version 字段
// 参数:
//   - version: 版本
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithVersion(version *string) *TArtifactBuilder {
	b.instance.Version = version
	return b
}

// WithVersionValue 设置 version 字段（便捷方法，自动转换为指针）
// 参数:
//   - version: 版本
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithVersionValue(version string) *TArtifactBuilder {
	b.instance.Version = &version
	return b
}

// WithCreateTime 设置 createTime 字段
// 参数:
//   - createTime: 创建时间
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithCreateTime(createTime time.Time) *TArtifactBuilder {
	b.instance.CreateTime = createTime
	return b
}

// WithUpdateTime 设置 updateTime 字段
// 参数:
//   - updateTime: 更新时间
//
// 返回:
//   - *TArtifactBuilder: 返回 Builder 实例，支持链式调用
func (b *TArtifactBuilder) WithUpdateTime(updateTime time.Time) *TArtifactBuilder {
	b.instance.UpdateTime = updateTime
	return b
}

// Build 构建并返回 TArtifact 实例
// 返回:
//   - *TArtifact: 构建完成的实例
func (b *TArtifactBuilder) Build() *TArtifact {
	return b.instance
}
