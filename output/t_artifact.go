package model

import (
	"time"
)

// TArtifact 任务执行流程中生成的中间产物表
type TArtifact struct {
	Id uint64 `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement;comment:主键ID;not null" json:"id"`
	ArtifactId uint64 `gorm:"column:artifactId;type:bigint(20);comment:产物ID;not null" json:"artifactId"`
	ArtifactName string `gorm:"column:artifactName;type:varchar(128);comment:产物名称;not null" json:"artifactName"`
	SessionId uint64 `gorm:"column:sessionId;type:bigint(20);comment:所属的会话;not null" json:"sessionId"`
	Step int `gorm:"column:step;type:int(11);comment:大的步骤点;not null" json:"step"`
	SubStep int `gorm:"column:subStep;type:int(11);comment:小的步骤点;not null" json:"subStep"`
	Content *string `gorm:"column:content;type:text;comment:内容;" json:"content"`
	Version *string `gorm:"column:version;type:varchar(128);comment:版本;" json:"version"`
	CreateTime time.Time `gorm:"column:createTime;type:datetime;comment:创建时间;not null" json:"createTime"`
	UpdateTime time.Time `gorm:"column:updateTime;type:datetime;comment:更新时间;not null" json:"updateTime"`
}

// TableName 返回表名
func (t *TArtifact) TableName() string {
	return "t_artifact"
}