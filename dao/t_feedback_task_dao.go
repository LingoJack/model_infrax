package dao

import (
	"fmt"
	"model_infrax/output/model/entity"
	"model_infrax/output/model/query"

	"strings"

	"gorm.io/gorm"
)

// TFeedbackTaskDAO 反馈任务表的DAO实现
type TFeedbackTaskDAO struct {
	db *gorm.DB
}

// NewTFeedbackTaskDAO 创建TFeedbackTaskDAO实例
// 参数:
//   - db: GORM数据库连接实例
//
// 返回:
//   - *TFeedbackTaskDAO: DAO实例
func NewTFeedbackTaskDAO(db *gorm.DB) *TFeedbackTaskDAO {
	return &TFeedbackTaskDAO{db: db}
}

// buildTFeedbackTaskQueryCondition 构建查询条件
// 参数:
//   - db: GORM数据库连接实例
//   - queryDTO: 查询条件DTO对象
//
// 返回:
//   - *gorm.DB: 应用了查询条件的数据库连接
//
// 说明:
//   - 支持精确匹配、模糊查询、IN查询、范围查询等多种查询方式
//   - 零值字段会被忽略，不会作为查询条件
func (dao *TFeedbackTaskDAO) buildTFeedbackTaskQueryCondition(db *gorm.DB, queryDTO *query.TFeedbackTaskDTO) *gorm.DB {
	if queryDTO == nil {
		return db
	}

	// 基础字段精确查询
	if queryDTO.Id != 0 {
		db = db.Where("id = ?", queryDTO.Id)
	}
	if queryDTO.FeedbackTaskId != "" {
		db = db.Where("feedbackTaskId = ?", queryDTO.FeedbackTaskId)
	}
	if queryDTO.AppId != "" {
		db = db.Where("appId = ?", queryDTO.AppId)
	}
	if queryDTO.AppName != "" {
		db = db.Where("appName = ?", queryDTO.AppName)
	}
	if queryDTO.Status != "" {
		db = db.Where("status = ?", queryDTO.Status)
	}
	if queryDTO.ApiName != nil && *queryDTO.ApiName != "" {
		db = db.Where("apiName = ?", *queryDTO.ApiName)
	}
	// bool类型字段：false也是有效值，这里简化处理，如需区分未设置和false，DTO应使用*bool
	if queryDTO.IsNeedFix {
		db = db.Where("isNeedFix = ?", queryDTO.IsNeedFix)
	}
	if queryDTO.FixPart != nil && *queryDTO.FixPart != "" {
		db = db.Where("fixPart = ?", *queryDTO.FixPart)
	}
	if queryDTO.Reason != nil && *queryDTO.Reason != "" {
		db = db.Where("reason = ?", *queryDTO.Reason)
	}
	if queryDTO.ConciseReason != nil && *queryDTO.ConciseReason != "" {
		db = db.Where("conciseReason = ?", *queryDTO.ConciseReason)
	}
	if queryDTO.BeforeCode != nil && *queryDTO.BeforeCode != "" {
		db = db.Where("beforeCode = ?", *queryDTO.BeforeCode)
	}
	if queryDTO.AfterCode != nil && *queryDTO.AfterCode != "" {
		db = db.Where("afterCode = ?", *queryDTO.AfterCode)
	}
	if queryDTO.BeforeMod != nil && *queryDTO.BeforeMod != "" {
		db = db.Where("beforeMod = ?", *queryDTO.BeforeMod)
	}
	if queryDTO.AfterMod != nil && *queryDTO.AfterMod != "" {
		db = db.Where("afterMod = ?", *queryDTO.AfterMod)
	}
	if queryDTO.ErrContext != nil && *queryDTO.ErrContext != "" {
		db = db.Where("errContext = ?", *queryDTO.ErrContext)
	}
	if queryDTO.HashCode != nil && *queryDTO.HashCode != "" {
		db = db.Where("hashCode = ?", *queryDTO.HashCode)
	}
	if queryDTO.Route != nil && *queryDTO.Route != "" {
		db = db.Where("route = ?", *queryDTO.Route)
	}
	if queryDTO.FeedbackSource != nil && *queryDTO.FeedbackSource != "" {
		db = db.Where("feedbackSource = ?", *queryDTO.FeedbackSource)
	}
	if queryDTO.UserFeedbackContent != nil && *queryDTO.UserFeedbackContent != "" {
		db = db.Where("userFeedbackContent = ?", *queryDTO.UserFeedbackContent)
	}
	if queryDTO.TaskGroupId != nil && *queryDTO.TaskGroupId != "" {
		db = db.Where("taskGroupId = ?", *queryDTO.TaskGroupId)
	}
	if !queryDTO.CreateTime.IsZero() {
		db = db.Where("createTime = ?", queryDTO.CreateTime)
	}
	if !queryDTO.UpdateTime.IsZero() {
		db = db.Where("updateTime = ?", queryDTO.UpdateTime)
	}

	// 模糊查询条件
	if queryDTO.FeedbackTaskIdFuzzy != "" {
		db = db.Where("feedbackTaskId LIKE ?", "%"+queryDTO.FeedbackTaskIdFuzzy+"%")
	}
	if queryDTO.AppIdFuzzy != "" {
		db = db.Where("appId LIKE ?", "%"+queryDTO.AppIdFuzzy+"%")
	}
	if queryDTO.AppNameFuzzy != "" {
		db = db.Where("appName LIKE ?", "%"+queryDTO.AppNameFuzzy+"%")
	}
	if queryDTO.StatusFuzzy != "" {
		db = db.Where("status LIKE ?", "%"+queryDTO.StatusFuzzy+"%")
	}
	if queryDTO.ApiNameFuzzy != nil && *queryDTO.ApiNameFuzzy != "" {
		db = db.Where("apiName LIKE ?", "%"+*queryDTO.ApiNameFuzzy+"%")
	}
	if queryDTO.FixPartFuzzy != nil && *queryDTO.FixPartFuzzy != "" {
		db = db.Where("fixPart LIKE ?", "%"+*queryDTO.FixPartFuzzy+"%")
	}
	if queryDTO.ReasonFuzzy != nil && *queryDTO.ReasonFuzzy != "" {
		db = db.Where("reason LIKE ?", "%"+*queryDTO.ReasonFuzzy+"%")
	}
	if queryDTO.ConciseReasonFuzzy != nil && *queryDTO.ConciseReasonFuzzy != "" {
		db = db.Where("conciseReason LIKE ?", "%"+*queryDTO.ConciseReasonFuzzy+"%")
	}
	if queryDTO.BeforeCodeFuzzy != nil && *queryDTO.BeforeCodeFuzzy != "" {
		db = db.Where("beforeCode LIKE ?", "%"+*queryDTO.BeforeCodeFuzzy+"%")
	}
	if queryDTO.AfterCodeFuzzy != nil && *queryDTO.AfterCodeFuzzy != "" {
		db = db.Where("afterCode LIKE ?", "%"+*queryDTO.AfterCodeFuzzy+"%")
	}
	if queryDTO.BeforeModFuzzy != nil && *queryDTO.BeforeModFuzzy != "" {
		db = db.Where("beforeMod LIKE ?", "%"+*queryDTO.BeforeModFuzzy+"%")
	}
	if queryDTO.AfterModFuzzy != nil && *queryDTO.AfterModFuzzy != "" {
		db = db.Where("afterMod LIKE ?", "%"+*queryDTO.AfterModFuzzy+"%")
	}
	if queryDTO.ErrContextFuzzy != nil && *queryDTO.ErrContextFuzzy != "" {
		db = db.Where("errContext LIKE ?", "%"+*queryDTO.ErrContextFuzzy+"%")
	}
	if queryDTO.HashCodeFuzzy != nil && *queryDTO.HashCodeFuzzy != "" {
		db = db.Where("hashCode LIKE ?", "%"+*queryDTO.HashCodeFuzzy+"%")
	}
	if queryDTO.RouteFuzzy != nil && *queryDTO.RouteFuzzy != "" {
		db = db.Where("route LIKE ?", "%"+*queryDTO.RouteFuzzy+"%")
	}
	if queryDTO.FeedbackSourceFuzzy != nil && *queryDTO.FeedbackSourceFuzzy != "" {
		db = db.Where("feedbackSource LIKE ?", "%"+*queryDTO.FeedbackSourceFuzzy+"%")
	}
	if queryDTO.UserFeedbackContentFuzzy != nil && *queryDTO.UserFeedbackContentFuzzy != "" {
		db = db.Where("userFeedbackContent LIKE ?", "%"+*queryDTO.UserFeedbackContentFuzzy+"%")
	}
	if queryDTO.TaskGroupIdFuzzy != nil && *queryDTO.TaskGroupIdFuzzy != "" {
		db = db.Where("taskGroupId LIKE ?", "%"+*queryDTO.TaskGroupIdFuzzy+"%")
	}

	// 日期范围查询
	if !queryDTO.CreateTimeStart.IsZero() {
		db = db.Where("createTime >= ?", queryDTO.CreateTimeStart)
	}
	if !queryDTO.CreateTimeEnd.IsZero() {
		db = db.Where("createTime < DATE_ADD(?, INTERVAL 1 DAY)", queryDTO.CreateTimeEnd)
	}
	if !queryDTO.UpdateTimeStart.IsZero() {
		db = db.Where("updateTime >= ?", queryDTO.UpdateTimeStart)
	}
	if !queryDTO.UpdateTimeEnd.IsZero() {
		db = db.Where("updateTime < DATE_ADD(?, INTERVAL 1 DAY)", queryDTO.UpdateTimeEnd)
	}

	// IN 查询条件
	if len(queryDTO.IdList) > 0 {
		db = db.Where("id IN ?", queryDTO.IdList)
	}
	if len(queryDTO.FeedbackTaskIdList) > 0 {
		db = db.Where("feedbackTaskId IN ?", queryDTO.FeedbackTaskIdList)
	}
	if len(queryDTO.AppIdList) > 0 {
		db = db.Where("appId IN ?", queryDTO.AppIdList)
	}
	if len(queryDTO.TaskGroupIdList) > 0 {
		db = db.Where("taskGroupId IN ?", queryDTO.TaskGroupIdList)
	}

	return db
}

// ==================== 基础查询方法 ====================

// SelectList 查询列表
// 参数:
//   - queryDTO: 查询条件DTO对象，支持分页、排序、多条件查询
//
// 返回:
//   - []*entity.TFeedbackTask: 查询结果列表
//   - error: 错误信息
func (dao *TFeedbackTaskDAO) SelectList(queryDTO *query.TFeedbackTaskDTO) ([]*entity.TFeedbackTask, error) {
	var resultList []*entity.TFeedbackTask
	db := dao.db.Model(&entity.TFeedbackTask{})

	// 应用查询条件
	db = dao.buildTFeedbackTaskQueryCondition(db, queryDTO)

	// 排序
	if queryDTO != nil && queryDTO.OrderBy != "" {
		if dao.isValidOrderBy(queryDTO.OrderBy) {
			db = db.Order(queryDTO.OrderBy)
		}
	}

	// 分页
	if queryDTO != nil && queryDTO.PageSize > 0 {
		db = db.Offset(queryDTO.PageOffset).Limit(queryDTO.PageSize)
	}

	err := db.Find(&resultList).Error
	return resultList, err
}

// SelectCount 查询数量
// 参数:
//   - queryDTO: 查询条件DTO对象
//
// 返回:
//   - int64: 符合条件的记录数量
//   - error: 错误信息
func (dao *TFeedbackTaskDAO) SelectCount(queryDTO *query.TFeedbackTaskDTO) (int64, error) {
	var count int64
	db := dao.db.Model(&entity.TFeedbackTask{})

	// 应用查询条件
	db = dao.buildTFeedbackTaskQueryCondition(db, queryDTO)

	err := db.Count(&count).Error
	return count, err
}

// ==================== 基础插入方法 ====================

// Insert 单行插入
// 参数:
//   - poBean: 要插入的PO对象
//
// 返回:
//   - error: 错误信息
//
// 说明:
//   - 插入所有字段，包括零值字段
//   - 自增主键会在插入后自动填充到poBean中
func (dao *TFeedbackTaskDAO) Insert(poBean *entity.TFeedbackTask) error {
	if poBean == nil {
		return fmt.Errorf("插入对象不能为空")
	}
	return dao.db.Create(poBean).Error
}

// InsertBatch 批量插入
// 参数:
//   - poBeanList: 要插入的PO对象列表
//
// 返回:
//   - error: 错误信息
//
// 说明:
//   - 批量插入所有记录，在一个事务中执行
//   - 自增主键会在插入后自动填充到各个poBean中
func (dao *TFeedbackTaskDAO) InsertBatch(poBeanList []*entity.TFeedbackTask) error {
	if len(poBeanList) == 0 {
		return fmt.Errorf("批量插入列表不能为空")
	}
	return dao.db.Create(&poBeanList).Error
}

// InsertOrUpdateNullable 插入或更新（会用零值覆盖）
// 参数:
//   - poBean: 要插入或更新的PO对象
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//  1. 如果记录不存在（根据主键判断），则执行插入操作
//  2. 如果记录已存在，则执行全字段更新操作
//  3. **重要**: 更新时会用传入对象的所有字段值覆盖数据库中的值，包括零值（nil、""、0、false等）
//     例如: 如果 poBean.Content = nil，会将数据库中的 content 字段更新为 NULL
//     例如: 如果 poBean.ArtifactName = ""，会将数据库中的 artifactName 字段更新为空字符串
//  4. 这种行为适用于需要"完整替换"记录的场景
//  5. 如果不希望零值覆盖数据库中的非零值，应使用 UpdateByXxx 等方法（内部使用 Updates）
func (dao *TFeedbackTaskDAO) InsertOrUpdateNullable(poBean *entity.TFeedbackTask) error {
	if poBean == nil {
		return fmt.Errorf("插入或更新对象不能为空")
	}
	// 使用 GORM 的 Save 方法:
	// - 根据主键判断记录是否存在
	// - 存在则更新所有字段（包括零值字段）
	// - 不存在则插入新记录
	return dao.db.Save(poBean).Error
}

// InsertOrUpdateBatchNullable 批量插入或更新（会用零值覆盖）
// 参数:
//   - poBeanList: 要插入或更新的PO对象列表
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//  1. 对列表中的每条记录，根据主键判断是插入还是更新
//  2. 如果记录不存在，则执行插入操作
//  3. 如果记录已存在，则执行全字段更新操作
//  4. **重要**: 更新时会用传入对象的所有字段值覆盖数据库中的值，包括零值（nil、""、0、false等）
//     这意味着如果某个字段在传入对象中为零值，会将数据库中对应字段更新为零值
//  5. 批量操作在一个事务中执行，要么全部成功，要么全部失败
//  6. 适用场景: 需要完整替换多条记录的场景
//  7. 性能提示: 批量操作比逐条调用 InsertOrUpdateNullable 效率更高
//  8. 如果不希望零值覆盖，建议逐条调用 UpdateByXxx 等方法
func (dao *TFeedbackTaskDAO) InsertOrUpdateBatchNullable(poBeanList []*entity.TFeedbackTask) error {
	if len(poBeanList) == 0 {
		return fmt.Errorf("批量插入或更新列表不能为空")
	}
	// 使用 GORM 的 Save 方法批量保存:
	// - 对每条记录根据主键判断是插入还是更新
	// - 更新时会覆盖所有字段（包括零值字段）
	// - 在一个事务中执行，保证原子性
	return dao.db.Save(&poBeanList).Error
}

// ==================== 主键索引方法 ====================

// SelectById 根据主键Id查询单条记录
// 参数:
//   - id: 主键值
//
// 返回:
//   - *entity.TFeedbackTask: 查询结果，如果不存在返回nil
//   - error: 错误信息
func (dao *TFeedbackTaskDAO) SelectById(id uint64) (*entity.TFeedbackTask, error) {
	var resultBean entity.TFeedbackTask
	err := dao.db.Where("id = ?", id).First(&resultBean).Error
	if err != nil {
		return nil, err
	}
	return &resultBean, nil
}

// UpdateById 根据主键Id更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - id: 主键值
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//  1. 根据指定的 id 更新记录
//  2. **重要**: 只更新非零值字段，零值字段会被忽略，不会覆盖数据库中的值
//     例如: 如果 poBean.Content = nil，不会更新数据库中的 content 字段
//     例如: 如果 poBean.ArtifactName = ""，不会更新数据库中的 artifactName 字段
//  3. 这种行为适用于"部分更新"场景，保留数据库中未传入的字段值
//  4. 如果需要将某个字段更新为零值，应使用 UpdateByIdWithMap 方法显式指定
//  5. 与 InsertOrUpdateNullable 的区别: InsertOrUpdateNullable 会用零值覆盖，UpdateById 不会
func (dao *TFeedbackTaskDAO) UpdateById(poBean *entity.TFeedbackTask, id uint64) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	// 使用 Updates 方法:
	// - 只更新结构体中的非零值字段
	// - 零值字段会被忽略，保留数据库中的原值
	// - 适合部分更新场景
	return dao.db.Model(&entity.TFeedbackTask{}).Where("id = ?", id).Updates(poBean).Error
}

// UpdateByIdWithMap 根据主键Id使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - id: 主键值
//   - updatedMap: 要更新的字段Map，key为字段名（数据库列名），value为字段值
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//  1. 根据指定的 id 更新记录
//  2. 使用 map 可以显式指定要更新的字段，包括零值字段
//  3. **重要**: 与 UpdateById 不同，使用 map 可以将字段更新为零值
//     例如: updatedMap["content"] = nil 会将 content 字段更新为 NULL
//     例如: updatedMap["artifactName"] = "" 会将 artifactName 字段更新为空字符串
//  4. 只更新 map 中指定的字段，未指定的字段保持不变
//  5. 适用场景: 需要精确控制更新哪些字段，包括需要将某些字段设置为零值的场景
//  6. 使用建议: 字段名必须与数据库列名一致（或使用 GORM 的字段映射名）
func (dao *TFeedbackTaskDAO) UpdateByIdWithMap(id uint64, updatedMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	// 使用 Updates 方法配合 map:
	// - 可以显式更新零值字段
	// - 只更新 map 中指定的字段
	// - 提供最精确的字段更新控制
	return dao.db.Model(&entity.TFeedbackTask{}).Where("id = ?", id).Updates(updatedMap).Error
}

// UpdateByIdWithCondition 根据主键Id和额外条件更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - id: 主键值
//   - conditionMap: 额外的查询条件Map，key为字段名，value为字段值
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//  1. 根据指定的 id 和额外的条件更新记录
//  2. 只更新非零值字段，零值字段会被忽略
//  3. 适用场景: 需要在主键基础上增加额外的更新条件，如乐观锁、状态检查等
//  4. 示例: conditionMap["version"] = 1 可以实现乐观锁，只有版本号匹配才更新
func (dao *TFeedbackTaskDAO) UpdateByIdWithCondition(poBean *entity.TFeedbackTask, id uint64, conditionMap map[string]interface{}) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TFeedbackTask{}).Where("id = ?", id)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(poBean).Error
}

// UpdateByIdWithMapAndCondition 根据主键Id和额外条件使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - id: 主键值
//   - updatedMap: 要更新的字段Map
//   - conditionMap: 额外的查询条件Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//  1. 根据指定的 id 和额外的条件更新记录
//  2. 使用 map 可以显式指定要更新的字段，包括零值字段
//  3. 提供最灵活的更新控制方式
func (dao *TFeedbackTaskDAO) UpdateByIdWithMapAndCondition(id uint64, updatedMap map[string]interface{}, conditionMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TFeedbackTask{}).Where("id = ?", id)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(updatedMap).Error
}

// DeleteById 根据主键Id删除
// 参数:
//   - id: 主键值
//
// 返回:
//   - error: 错误信息
func (dao *TFeedbackTaskDAO) DeleteById(id uint64) error {
	return dao.db.Where("id = ?", id).Delete(&entity.TFeedbackTask{}).Error
}

// ==================== 唯一索引 uk_feedbackTaskId 方法 ====================

// SelectByFeedbackTaskId 根据唯一索引uk_feedbackTaskId查询单条记录
// 参数:
//   - feedbackTaskId: 任务ID
//
// 返回:
//   - *entity.TFeedbackTask: 查询结果，如果不存在返回nil
//   - error: 错误信息
func (dao *TFeedbackTaskDAO) SelectByFeedbackTaskId(feedbackTaskId string) (*entity.TFeedbackTask, error) {
	var resultBean entity.TFeedbackTask
	err := dao.db.Where("feedbackTaskId = ?", feedbackTaskId).First(&resultBean).Error
	if err != nil {
		return nil, err
	}
	return &resultBean, nil
}

// UpdateByFeedbackTaskId 根据唯一索引uk_feedbackTaskId更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - feedbackTaskId: 任务ID
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 只更新非零值字段，零值字段会被忽略
func (dao *TFeedbackTaskDAO) UpdateByFeedbackTaskId(poBean *entity.TFeedbackTask, feedbackTaskId string) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	return dao.db.Model(&entity.TFeedbackTask{}).Where("feedbackTaskId = ?", feedbackTaskId).Updates(poBean).Error
}

// UpdateByFeedbackTaskIdWithMap 根据唯一索引uk_feedbackTaskId使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - feedbackTaskId: 任务ID
//   - updatedMap: 要更新的字段Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 使用 map 可以显式指定要更新的字段，包括零值字段
//   - 只更新 map 中指定的字段，未指定的字段保持不变
func (dao *TFeedbackTaskDAO) UpdateByFeedbackTaskIdWithMap(feedbackTaskId string, updatedMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	return dao.db.Model(&entity.TFeedbackTask{}).Where("feedbackTaskId = ?", feedbackTaskId).Updates(updatedMap).Error
}

// UpdateByFeedbackTaskIdWithCondition 根据唯一索引uk_feedbackTaskId和额外条件更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - feedbackTaskId: 任务ID
//   - conditionMap: 额外的查询条件Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 只更新非零值字段，零值字段会被忽略
//   - 适用场景: 需要在唯一键基础上增加额外的更新条件
func (dao *TFeedbackTaskDAO) UpdateByFeedbackTaskIdWithCondition(poBean *entity.TFeedbackTask, feedbackTaskId string, conditionMap map[string]interface{}) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TFeedbackTask{}).Where("feedbackTaskId = ?", feedbackTaskId)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(poBean).Error
}

// UpdateByFeedbackTaskIdWithMapAndCondition 根据唯一索引uk_feedbackTaskId和额外条件使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - feedbackTaskId: 任务ID
//   - updatedMap: 要更新的字段Map
//   - conditionMap: 额外的查询条件Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 使用 map 可以显式指定要更新的字段，包括零值字段
//   - 提供最灵活的更新控制方式
func (dao *TFeedbackTaskDAO) UpdateByFeedbackTaskIdWithMapAndCondition(feedbackTaskId string, updatedMap map[string]interface{}, conditionMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TFeedbackTask{}).Where("feedbackTaskId = ?", feedbackTaskId)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(updatedMap).Error
}

// DeleteByFeedbackTaskId 根据唯一索引uk_feedbackTaskId删除
// 参数:
//   - feedbackTaskId: 任务ID
//
// 返回:
//   - error: 错误信息
func (dao *TFeedbackTaskDAO) DeleteByFeedbackTaskId(feedbackTaskId string) error {
	return dao.db.Where("feedbackTaskId = ?", feedbackTaskId).Delete(&entity.TFeedbackTask{}).Error
}

// ==================== 普通索引 idx_taskGroupId 方法 ====================

// SelectByTaskGroupId 根据索引idx_taskGroupId查询列表
// 参数:
//   - taskGroupId: 任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同
//
// 返回:
//   - []*entity.TFeedbackTask: 查询结果列表
//   - error: 错误信息
//
// 说明:
//   - 该索引不是唯一索引，可能返回多条记录
func (dao *TFeedbackTaskDAO) SelectByTaskGroupId(taskGroupId *string) ([]*entity.TFeedbackTask, error) {
	var resultList []*entity.TFeedbackTask
	err := dao.db.Where("taskGroupId = ?", taskGroupId).Find(&resultList).Error
	return resultList, err
}

// SelectByTaskGroupIdList 根据索引idx_taskGroupId批量查询列表
// 参数:
//   - taskGroupIdList: 任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同列表
//
// 返回:
//   - []*entity.TFeedbackTask: 查询结果列表
//   - error: 错误信息
func (dao *TFeedbackTaskDAO) SelectByTaskGroupIdList(taskGroupIdList []*string) ([]*entity.TFeedbackTask, error) {
	if len(taskGroupIdList) == 0 {
		return []*entity.TFeedbackTask{}, nil
	}
	var resultList []*entity.TFeedbackTask
	err := dao.db.Where("taskGroupId IN ?", taskGroupIdList).Find(&resultList).Error
	return resultList, err
}

// UpdateByTaskGroupId 根据索引idx_taskGroupId更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - taskGroupId: 任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 只更新非零值字段，零值字段会被忽略
//   - 注意: taskGroupId 不是唯一键，可能会更新多条记录
func (dao *TFeedbackTaskDAO) UpdateByTaskGroupId(poBean *entity.TFeedbackTask, taskGroupId *string) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	return dao.db.Model(&entity.TFeedbackTask{}).Where("taskGroupId = ?", taskGroupId).Updates(poBean).Error
}

// UpdateByTaskGroupIdWithMap 根据索引idx_taskGroupId使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - taskGroupId: 任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同
//   - updatedMap: 要更新的字段Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 使用 map 可以显式指定要更新的字段，包括零值字段
//   - 只更新 map 中指定的字段，未指定的字段保持不变
//   - 注意: taskGroupId 不是唯一键，可能会更新多条记录
func (dao *TFeedbackTaskDAO) UpdateByTaskGroupIdWithMap(taskGroupId *string, updatedMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	return dao.db.Model(&entity.TFeedbackTask{}).Where("taskGroupId = ?", taskGroupId).Updates(updatedMap).Error
}

// UpdateByTaskGroupIdWithCondition 根据索引idx_taskGroupId和额外条件更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - taskGroupId: 任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同
//   - conditionMap: 额外的查询条件Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 只更新非零值字段，零值字段会被忽略
//   - 适用场景: 需要在 taskGroupId 基础上增加额外的更新条件，缩小更新范围
//   - 注意: 可能会更新多条记录
func (dao *TFeedbackTaskDAO) UpdateByTaskGroupIdWithCondition(poBean *entity.TFeedbackTask, taskGroupId *string, conditionMap map[string]interface{}) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TFeedbackTask{}).Where("taskGroupId = ?", taskGroupId)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(poBean).Error
}

// UpdateByTaskGroupIdWithMapAndCondition 根据索引idx_taskGroupId和额外条件使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - taskGroupId: 任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同
//   - updatedMap: 要更新的字段Map
//   - conditionMap: 额外的查询条件Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 使用 map 可以显式指定要更新的字段，包括零值字段
//   - 提供最灵活的更新控制方式
//   - 注意: 可能会更新多条记录
func (dao *TFeedbackTaskDAO) UpdateByTaskGroupIdWithMapAndCondition(taskGroupId *string, updatedMap map[string]interface{}, conditionMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TFeedbackTask{}).Where("taskGroupId = ?", taskGroupId)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(updatedMap).Error
}

// DeleteByTaskGroupId 根据索引idx_taskGroupId删除
// 参数:
//   - taskGroupId: 任务组ID，例如在一个用户反馈任务中，可能有多个要求改的api，那么会生成多个task，他们的taskGroupId相同
//
// 返回:
//   - error: 错误信息
//
// 说明:
//   - 注意: taskGroupId 不是唯一键，可能会删除多条记录
func (dao *TFeedbackTaskDAO) DeleteByTaskGroupId(taskGroupId *string) error {
	return dao.db.Where("taskGroupId = ?", taskGroupId).Delete(&entity.TFeedbackTask{}).Error
}

// ==================== 普通索引 idx_appId 方法 ====================

// SelectByAppId 根据索引idx_appId查询列表
// 参数:
//   - appId: 应用ID
//
// 返回:
//   - []*entity.TFeedbackTask: 查询结果列表
//   - error: 错误信息
//
// 说明:
//   - 该索引不是唯一索引，可能返回多条记录
func (dao *TFeedbackTaskDAO) SelectByAppId(appId string) ([]*entity.TFeedbackTask, error) {
	var resultList []*entity.TFeedbackTask
	err := dao.db.Where("appId = ?", appId).Find(&resultList).Error
	return resultList, err
}

// SelectByAppIdList 根据索引idx_appId批量查询列表
// 参数:
//   - appIdList: 应用ID列表
//
// 返回:
//   - []*entity.TFeedbackTask: 查询结果列表
//   - error: 错误信息
func (dao *TFeedbackTaskDAO) SelectByAppIdList(appIdList []string) ([]*entity.TFeedbackTask, error) {
	if len(appIdList) == 0 {
		return []*entity.TFeedbackTask{}, nil
	}
	var resultList []*entity.TFeedbackTask
	err := dao.db.Where("appId IN ?", appIdList).Find(&resultList).Error
	return resultList, err
}

// UpdateByAppId 根据索引idx_appId更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - appId: 应用ID
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 只更新非零值字段，零值字段会被忽略
//   - 注意: appId 不是唯一键，可能会更新多条记录
func (dao *TFeedbackTaskDAO) UpdateByAppId(poBean *entity.TFeedbackTask, appId string) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	return dao.db.Model(&entity.TFeedbackTask{}).Where("appId = ?", appId).Updates(poBean).Error
}

// UpdateByAppIdWithMap 根据索引idx_appId使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - appId: 应用ID
//   - updatedMap: 要更新的字段Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 使用 map 可以显式指定要更新的字段，包括零值字段
//   - 只更新 map 中指定的字段，未指定的字段保持不变
//   - 注意: appId 不是唯一键，可能会更新多条记录
func (dao *TFeedbackTaskDAO) UpdateByAppIdWithMap(appId string, updatedMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	return dao.db.Model(&entity.TFeedbackTask{}).Where("appId = ?", appId).Updates(updatedMap).Error
}

// UpdateByAppIdWithCondition 根据索引idx_appId和额外条件更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - appId: 应用ID
//   - conditionMap: 额外的查询条件Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 只更新非零值字段，零值字段会被忽略
//   - 适用场景: 需要在 appId 基础上增加额外的更新条件，缩小更新范围
//   - 注意: 可能会更新多条记录
func (dao *TFeedbackTaskDAO) UpdateByAppIdWithCondition(poBean *entity.TFeedbackTask, appId string, conditionMap map[string]interface{}) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TFeedbackTask{}).Where("appId = ?", appId)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(poBean).Error
}

// UpdateByAppIdWithMapAndCondition 根据索引idx_appId和额外条件使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - appId: 应用ID
//   - updatedMap: 要更新的字段Map
//   - conditionMap: 额外的查询条件Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 使用 map 可以显式指定要更新的字段，包括零值字段
//   - 提供最灵活的更新控制方式
//   - 注意: 可能会更新多条记录
func (dao *TFeedbackTaskDAO) UpdateByAppIdWithMapAndCondition(appId string, updatedMap map[string]interface{}, conditionMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TFeedbackTask{}).Where("appId = ?", appId)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(updatedMap).Error
}

// DeleteByAppId 根据索引idx_appId删除
// 参数:
//   - appId: 应用ID
//
// 返回:
//   - error: 错误信息
//
// 说明:
//   - 注意: appId 不是唯一键，可能会删除多条记录
func (dao *TFeedbackTaskDAO) DeleteByAppId(appId string) error {
	return dao.db.Where("appId = ?", appId).Delete(&entity.TFeedbackTask{}).Error
}

// ==================== 辅助方法 ====================

// getValidOrderByFields 获取允许排序的字段白名单
// 返回:
//   - map[string]bool: 字段白名单，key为字段名，value为true表示允许排序
func (dao *TFeedbackTaskDAO) getValidOrderByFields() map[string]bool {
	return map[string]bool{
		"id":                  true,
		"feedbackTaskId":      true,
		"appId":               true,
		"appName":             true,
		"status":              true,
		"apiName":             true,
		"isNeedFix":           true,
		"fixPart":             true,
		"reason":              true,
		"conciseReason":       true,
		"beforeCode":          true,
		"afterCode":           true,
		"beforeMod":           true,
		"afterMod":            true,
		"errContext":          true,
		"hashCode":            true,
		"route":               true,
		"feedbackSource":      true,
		"userFeedbackContent": true,
		"taskGroupId":         true,
		"createTime":          true,
		"updateTime":          true,
	}
}

// isValidOrderBy 验证排序字符串是否安全（基于字段白名单）
// 支持格式:
//   - 单字段: id DESC
//   - 多字段: id DESC, createTime ASC
//
// 参数:
//   - orderBy: 排序字符串
//
// 返回:
//   - true: 排序字符串合法且所有字段都在白名单中
//   - false: 排序字符串不合法或包含非白名单字段
func (dao *TFeedbackTaskDAO) isValidOrderBy(orderBy string) bool {
	if orderBy == "" {
		return false
	}

	// 获取字段白名单
	validFields := dao.getValidOrderByFields()

	// 按逗号分割多个排序字段
	orderParts := strings.Split(orderBy, ",")

	for _, part := range orderParts {
		part = strings.TrimSpace(part)
		if part == "" {
			return false
		}

		// 按空格分割字段名和排序方向
		tokens := strings.Fields(part)
		if len(tokens) == 0 || len(tokens) > 2 {
			// 格式错误: 必须是 "字段名" 或 "字段名 方向"
			return false
		}

		// 验证字段名是否在白名单中
		fieldName := tokens[0]
		if !validFields[fieldName] {
			// 字段不在白名单中
			return false
		}

		// 如果指定了排序方向，验证是否为 ASC 或 DESC
		if len(tokens) == 2 {
			direction := strings.ToUpper(tokens[1])
			if direction != "ASC" && direction != "DESC" {
				// 排序方向无效
				return false
			}
		}
	}

	return true
}
