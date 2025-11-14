package model

import (
	"encoding/json"
	"time"
)

// TFeedbackTask 反馈任务表
type TFeedbackTask struct {
	Id                  uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement;comment:主键ID;not null" json:"id"`
	FeedbackTaskId      string    `gorm:"column:feedbackTaskId;type:varchar(128);comment:任务ID;not null" json:"feedbackTaskId"`
	AppId               string    `gorm:"column:appId;type:varchar(128);comment:应用ID;not null" json:"appId"`
	AppName             string    `gorm:"column:appName;type:varchar(128);comment:应用名称;not null" json:"appName"`
	Status              string    `gorm:"column:status;type:varchar(50);comment:修复状态: success, running, fail, cancel, init, confirm;not null" json:"status"`
	ApiName             *string   `gorm:"column:apiName;type:varchar(256);comment:API名称;" json:"apiName"`
	IsNeedFix           bool      `gorm:"column:isNeedFix;type:tinyint(1);comment:是否需要修复: 0-否 1-是;not null" json:"isNeedFix"`
	FixPart             *string   `gorm:"column:fixPart;type:varchar(10);comment:需要修复的部分: backend, frontend, all;" json:"fixPart"`
	Reason              *string   `gorm:"column:reason;type:text;comment:错误原因及修复说明;" json:"reason"`
	ConciseReason       *string   `gorm:"column:conciseReason;type:varchar(128);comment:错误原因及修复说明(简洁);" json:"conciseReason"`
	BeforeCode          *string   `gorm:"column:beforeCode;type:text;comment:修复前代码;" json:"beforeCode"`
	AfterCode           *string   `gorm:"column:afterCode;type:text;comment:修复后代码;" json:"afterCode"`
	BeforeMod           *string   `gorm:"column:beforeMod;type:text;comment:修复前模块;" json:"beforeMod"`
	AfterMod            *string   `gorm:"column:afterMod;type:text;comment:修复后模块;" json:"afterMod"`
	ErrContext          *string   `gorm:"column:errContext;type:text;comment:错误上下文;" json:"errContext"`
	HashCode            *string   `gorm:"column:hashCode;type:varchar(128);comment:hash值;" json:"hashCode"`
	Route               *string   `gorm:"column:route;type:varchar(256);comment:FIP协议描述的前端页面路由;" json:"route"`
	FeedbackSource      *string   `gorm:"column:feedbackSource;type:varchar(20);comment:反馈来源: system(兼容旧逻辑，为空默认system), user;" json:"feedbackSource"`
	UserFeedbackContent *string   `gorm:"column:userFeedbackContent;type:text;comment:用户反馈内容，只有feedbackSource为user时，该字段才有效;" json:"userFeedbackContent"`
	TaskGroupId         *string   `gorm:"column:taskGroupId;type:varchar(128);comment:任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同;" json:"taskGroupId"`
	CreateTime          time.Time `gorm:"column:createTime;type:datetime;comment:创建时间;not null" json:"createTime"`
	UpdateTime          time.Time `gorm:"column:updateTime;type:datetime;comment:更新时间;not null" json:"updateTime"`
}

// TableName 返回表名
func (t *TFeedbackTask) TableName() string {
	return "t_feedback_task"
}

// Jsonify 将结构体序列化为 JSON 字符串（紧凑格式）
// 返回:
//   - string: JSON 字符串，如果序列化失败则返回错误信息的 JSON
func (t *TFeedbackTask) Jsonify() string {
	byts, err := json.Marshal(t)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(byts)
}

// JsonifyIndent 将结构体序列化为格式化的 JSON 字符串（带缩进）
// 返回:
//   - string: 格式化的 JSON 字符串，如果序列化失败则返回错误信息的 JSON
func (t *TFeedbackTask) JsonifyIndent() string {
	byts, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(byts)
}

// TFeedbackTaskBuilder 用于构建 TFeedbackTask 实例的 Builder
type TFeedbackTaskBuilder struct {
	instance *TFeedbackTask
}

// NewTFeedbackTaskBuilder 创建一个新的 TFeedbackTaskBuilder 实例
// 返回:
//   - *TFeedbackTaskBuilder: Builder 实例，用于链式调用
func NewTFeedbackTaskBuilder() *TFeedbackTaskBuilder {
	return &TFeedbackTaskBuilder{
		instance: &TFeedbackTask{},
	}
}

// WithFeedbackTaskId 设置 feedbackTaskId 字段
// 参数:
//   - feedbackTaskId: 任务ID
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithFeedbackTaskId(feedbackTaskId string) *TFeedbackTaskBuilder {
	b.instance.FeedbackTaskId = feedbackTaskId
	return b
}

// WithAppId 设置 appId 字段
// 参数:
//   - appId: 应用ID
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithAppId(appId string) *TFeedbackTaskBuilder {
	b.instance.AppId = appId
	return b
}

// WithAppName 设置 appName 字段
// 参数:
//   - appName: 应用名称
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithAppName(appName string) *TFeedbackTaskBuilder {
	b.instance.AppName = appName
	return b
}

// WithStatus 设置 status 字段
// 参数:
//   - status: 修复状态: success, running, fail, cancel, init, confirm
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithStatus(status string) *TFeedbackTaskBuilder {
	b.instance.Status = status
	return b
}

// WithApiName 设置 apiName 字段
// 参数:
//   - apiName: API名称
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithApiName(apiName *string) *TFeedbackTaskBuilder {
	b.instance.ApiName = apiName
	return b
}

// WithApiNameValue 设置 apiName 字段（便捷方法，自动转换为指针）
// 参数:
//   - apiName: API名称
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithApiNameValue(apiName string) *TFeedbackTaskBuilder {
	b.instance.ApiName = &apiName
	return b
}

// WithIsNeedFix 设置 isNeedFix 字段
// 参数:
//   - isNeedFix: 是否需要修复: 0-否 1-是
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithIsNeedFix(isNeedFix bool) *TFeedbackTaskBuilder {
	b.instance.IsNeedFix = isNeedFix
	return b
}

// WithFixPart 设置 fixPart 字段
// 参数:
//   - fixPart: 需要修复的部分: backend, frontend, all
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithFixPart(fixPart *string) *TFeedbackTaskBuilder {
	b.instance.FixPart = fixPart
	return b
}

// WithFixPartValue 设置 fixPart 字段（便捷方法，自动转换为指针）
// 参数:
//   - fixPart: 需要修复的部分: backend, frontend, all
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithFixPartValue(fixPart string) *TFeedbackTaskBuilder {
	b.instance.FixPart = &fixPart
	return b
}

// WithReason 设置 reason 字段
// 参数:
//   - reason: 错误原因及修复说明
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithReason(reason *string) *TFeedbackTaskBuilder {
	b.instance.Reason = reason
	return b
}

// WithReasonValue 设置 reason 字段（便捷方法，自动转换为指针）
// 参数:
//   - reason: 错误原因及修复说明
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithReasonValue(reason string) *TFeedbackTaskBuilder {
	b.instance.Reason = &reason
	return b
}

// WithConciseReason 设置 conciseReason 字段
// 参数:
//   - conciseReason: 错误原因及修复说明(简洁)
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithConciseReason(conciseReason *string) *TFeedbackTaskBuilder {
	b.instance.ConciseReason = conciseReason
	return b
}

// WithConciseReasonValue 设置 conciseReason 字段（便捷方法，自动转换为指针）
// 参数:
//   - conciseReason: 错误原因及修复说明(简洁)
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithConciseReasonValue(conciseReason string) *TFeedbackTaskBuilder {
	b.instance.ConciseReason = &conciseReason
	return b
}

// WithBeforeCode 设置 beforeCode 字段
// 参数:
//   - beforeCode: 修复前代码
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithBeforeCode(beforeCode *string) *TFeedbackTaskBuilder {
	b.instance.BeforeCode = beforeCode
	return b
}

// WithBeforeCodeValue 设置 beforeCode 字段（便捷方法，自动转换为指针）
// 参数:
//   - beforeCode: 修复前代码
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithBeforeCodeValue(beforeCode string) *TFeedbackTaskBuilder {
	b.instance.BeforeCode = &beforeCode
	return b
}

// WithAfterCode 设置 afterCode 字段
// 参数:
//   - afterCode: 修复后代码
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithAfterCode(afterCode *string) *TFeedbackTaskBuilder {
	b.instance.AfterCode = afterCode
	return b
}

// WithAfterCodeValue 设置 afterCode 字段（便捷方法，自动转换为指针）
// 参数:
//   - afterCode: 修复后代码
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithAfterCodeValue(afterCode string) *TFeedbackTaskBuilder {
	b.instance.AfterCode = &afterCode
	return b
}

// WithBeforeMod 设置 beforeMod 字段
// 参数:
//   - beforeMod: 修复前模块
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithBeforeMod(beforeMod *string) *TFeedbackTaskBuilder {
	b.instance.BeforeMod = beforeMod
	return b
}

// WithBeforeModValue 设置 beforeMod 字段（便捷方法，自动转换为指针）
// 参数:
//   - beforeMod: 修复前模块
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithBeforeModValue(beforeMod string) *TFeedbackTaskBuilder {
	b.instance.BeforeMod = &beforeMod
	return b
}

// WithAfterMod 设置 afterMod 字段
// 参数:
//   - afterMod: 修复后模块
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithAfterMod(afterMod *string) *TFeedbackTaskBuilder {
	b.instance.AfterMod = afterMod
	return b
}

// WithAfterModValue 设置 afterMod 字段（便捷方法，自动转换为指针）
// 参数:
//   - afterMod: 修复后模块
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithAfterModValue(afterMod string) *TFeedbackTaskBuilder {
	b.instance.AfterMod = &afterMod
	return b
}

// WithErrContext 设置 errContext 字段
// 参数:
//   - errContext: 错误上下文
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithErrContext(errContext *string) *TFeedbackTaskBuilder {
	b.instance.ErrContext = errContext
	return b
}

// WithErrContextValue 设置 errContext 字段（便捷方法，自动转换为指针）
// 参数:
//   - errContext: 错误上下文
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithErrContextValue(errContext string) *TFeedbackTaskBuilder {
	b.instance.ErrContext = &errContext
	return b
}

// WithHashCode 设置 hashCode 字段
// 参数:
//   - hashCode: hash值
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithHashCode(hashCode *string) *TFeedbackTaskBuilder {
	b.instance.HashCode = hashCode
	return b
}

// WithHashCodeValue 设置 hashCode 字段（便捷方法，自动转换为指针）
// 参数:
//   - hashCode: hash值
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithHashCodeValue(hashCode string) *TFeedbackTaskBuilder {
	b.instance.HashCode = &hashCode
	return b
}

// WithRoute 设置 route 字段
// 参数:
//   - route: FIP协议描述的前端页面路由
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithRoute(route *string) *TFeedbackTaskBuilder {
	b.instance.Route = route
	return b
}

// WithRouteValue 设置 route 字段（便捷方法，自动转换为指针）
// 参数:
//   - route: FIP协议描述的前端页面路由
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithRouteValue(route string) *TFeedbackTaskBuilder {
	b.instance.Route = &route
	return b
}

// WithFeedbackSource 设置 feedbackSource 字段
// 参数:
//   - feedbackSource: 反馈来源: system(兼容旧逻辑，为空默认system), user
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithFeedbackSource(feedbackSource *string) *TFeedbackTaskBuilder {
	b.instance.FeedbackSource = feedbackSource
	return b
}

// WithFeedbackSourceValue 设置 feedbackSource 字段（便捷方法，自动转换为指针）
// 参数:
//   - feedbackSource: 反馈来源: system(兼容旧逻辑，为空默认system), user
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithFeedbackSourceValue(feedbackSource string) *TFeedbackTaskBuilder {
	b.instance.FeedbackSource = &feedbackSource
	return b
}

// WithUserFeedbackContent 设置 userFeedbackContent 字段
// 参数:
//   - userFeedbackContent: 用户反馈内容，只有feedbackSource为user时，该字段才有效
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithUserFeedbackContent(userFeedbackContent *string) *TFeedbackTaskBuilder {
	b.instance.UserFeedbackContent = userFeedbackContent
	return b
}

// WithUserFeedbackContentValue 设置 userFeedbackContent 字段（便捷方法，自动转换为指针）
// 参数:
//   - userFeedbackContent: 用户反馈内容，只有feedbackSource为user时，该字段才有效
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithUserFeedbackContentValue(userFeedbackContent string) *TFeedbackTaskBuilder {
	b.instance.UserFeedbackContent = &userFeedbackContent
	return b
}

// WithTaskGroupId 设置 taskGroupId 字段
// 参数:
//   - taskGroupId: 任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithTaskGroupId(taskGroupId *string) *TFeedbackTaskBuilder {
	b.instance.TaskGroupId = taskGroupId
	return b
}

// WithTaskGroupIdValue 设置 taskGroupId 字段（便捷方法，自动转换为指针）
// 参数:
//   - taskGroupId: 任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithTaskGroupIdValue(taskGroupId string) *TFeedbackTaskBuilder {
	b.instance.TaskGroupId = &taskGroupId
	return b
}

// WithCreateTime 设置 createTime 字段
// 参数:
//   - createTime: 创建时间
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithCreateTime(createTime time.Time) *TFeedbackTaskBuilder {
	b.instance.CreateTime = createTime
	return b
}

// WithUpdateTime 设置 updateTime 字段
// 参数:
//   - updateTime: 更新时间
//
// 返回:
//   - *TFeedbackTaskBuilder: 返回 Builder 实例，支持链式调用
func (b *TFeedbackTaskBuilder) WithUpdateTime(updateTime time.Time) *TFeedbackTaskBuilder {
	b.instance.UpdateTime = updateTime
	return b
}

// Build 构建并返回 TFeedbackTask 实例
// 返回:
//   - *TFeedbackTask: 构建完成的实例
func (b *TFeedbackTaskBuilder) Build() *TFeedbackTask {
	return b.instance
}
