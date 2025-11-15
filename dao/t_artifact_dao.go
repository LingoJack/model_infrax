package dao

import (
	"fmt"
	"model_infrax/output/model/entity"
	"model_infrax/output/model/query"
	"strings"

	"gorm.io/gorm"
)

// TArtifactDAO 任务执行流程中生成的中间产物表的DAO实现
type TArtifactDAO struct {
	db *gorm.DB
}

// NewTArtifactDAO 创建TArtifactDAO实例
// 参数:
//   - db: GORM数据库连接实例
//
// 返回:
//   - *TArtifactDAO: DAO实例
func NewTArtifactDAO(db *gorm.DB) *TArtifactDAO {
	return &TArtifactDAO{db: db}
}

// buildTArtifactQueryCondition 构建查询条件
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
func (dao *TArtifactDAO) buildTArtifactQueryCondition(db *gorm.DB, queryDTO *query.TArtifactDTO) *gorm.DB {
	if queryDTO == nil {
		return db
	}

	// 基础字段精确查询
	if queryDTO.Id != 0 {
		db = db.Where("id = ?", queryDTO.Id)
	}
	if queryDTO.ArtifactId != "" {
		db = db.Where("artifactId = ?", queryDTO.ArtifactId)
	}
	if queryDTO.ArtifactName != "" {
		db = db.Where("artifactName = ?", queryDTO.ArtifactName)
	}
	if queryDTO.SessionId != "" {
		db = db.Where("sessionId = ?", queryDTO.SessionId)
	}
	if queryDTO.Step != 0 {
		db = db.Where("step = ?", queryDTO.Step)
	}
	if queryDTO.SubStep != "" {
		db = db.Where("subStep = ?", queryDTO.SubStep)
	}
	if queryDTO.Content != nil && *queryDTO.Content != "" {
		db = db.Where("content = ?", *queryDTO.Content)
	}
	if queryDTO.Version != nil && *queryDTO.Version != "" {
		db = db.Where("version = ?", *queryDTO.Version)
	}
	if !queryDTO.CreateTime.IsZero() {
		db = db.Where("createTime = ?", queryDTO.CreateTime)
	}
	if !queryDTO.UpdateTime.IsZero() {
		db = db.Where("updateTime = ?", queryDTO.UpdateTime)
	}

	// 模糊查询条件
	if queryDTO.ArtifactIdFuzzy != "" {
		db = db.Where("artifactId LIKE ?", "%"+queryDTO.ArtifactIdFuzzy+"%")
	}
	if queryDTO.ArtifactNameFuzzy != "" {
		db = db.Where("artifactName LIKE ?", "%"+queryDTO.ArtifactNameFuzzy+"%")
	}
	if queryDTO.SessionIdFuzzy != "" {
		db = db.Where("sessionId LIKE ?", "%"+queryDTO.SessionIdFuzzy+"%")
	}
	if queryDTO.SubStepFuzzy != "" {
		db = db.Where("subStep LIKE ?", "%"+queryDTO.SubStepFuzzy+"%")
	}
	if queryDTO.ContentFuzzy != nil && *queryDTO.ContentFuzzy != "" {
		db = db.Where("content LIKE ?", "%"+*queryDTO.ContentFuzzy+"%")
	}
	if queryDTO.VersionFuzzy != nil && *queryDTO.VersionFuzzy != "" {
		db = db.Where("version LIKE ?", "%"+*queryDTO.VersionFuzzy+"%")
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
	if len(queryDTO.ArtifactIdList) > 0 {
		db = db.Where("artifactId IN ?", queryDTO.ArtifactIdList)
	}
	if len(queryDTO.SessionIdList) > 0 {
		db = db.Where("sessionId IN ?", queryDTO.SessionIdList)
	}
	if len(queryDTO.StepList) > 0 {
		db = db.Where("step IN ?", queryDTO.StepList)
	}
	if len(queryDTO.VersionList) > 0 {
		db = db.Where("version IN ?", queryDTO.VersionList)
	}

	return db
}

// ==================== 基础查询方法 ====================

// SelectList 查询列表
// 参数:
//   - queryDTO: 查询条件DTO对象，支持分页、排序、多条件查询
//
// 返回:
//   - []*entity.TArtifact: 查询结果列表
//   - error: 错误信息
func (dao *TArtifactDAO) SelectList(queryDTO *query.TArtifactDTO) ([]*entity.TArtifact, error) {
	var resultList []*entity.TArtifact
	db := dao.db.Model(&entity.TArtifact{})

	// 应用查询条件
	db = dao.buildTArtifactQueryCondition(db, queryDTO)

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
func (dao *TArtifactDAO) SelectCount(queryDTO *query.TArtifactDTO) (int64, error) {
	var count int64
	db := dao.db.Model(&entity.TArtifact{})

	// 应用查询条件
	db = dao.buildTArtifactQueryCondition(db, queryDTO)

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
func (dao *TArtifactDAO) Insert(poBean *entity.TArtifact) error {
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
func (dao *TArtifactDAO) InsertBatch(poBeanList []*entity.TArtifact) error {
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
func (dao *TArtifactDAO) InsertOrUpdateNullable(poBean *entity.TArtifact) error {
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
func (dao *TArtifactDAO) InsertOrUpdateBatchNullable(poBeanList []*entity.TArtifact) error {
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
//   - *entity.TArtifact: 查询结果，如果不存在返回nil
//   - error: 错误信息
func (dao *TArtifactDAO) SelectById(id uint64) (*entity.TArtifact, error) {
	var resultBean entity.TArtifact
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
func (dao *TArtifactDAO) UpdateById(poBean *entity.TArtifact, id uint64) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	// 使用 Updates 方法:
	// - 只更新结构体中的非零值字段
	// - 零值字段会被忽略，保留数据库中的原值
	// - 适合部分更新场景
	return dao.db.Model(&entity.TArtifact{}).Where("id = ?", id).Updates(poBean).Error
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
func (dao *TArtifactDAO) UpdateByIdWithMap(id uint64, updatedMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	// 使用 Updates 方法配合 map:
	// - 可以显式更新零值字段
	// - 只更新 map 中指定的字段
	// - 提供最精确的字段更新控制
	return dao.db.Model(&entity.TArtifact{}).Where("id = ?", id).Updates(updatedMap).Error
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
func (dao *TArtifactDAO) UpdateByIdWithCondition(poBean *entity.TArtifact, id uint64, conditionMap map[string]interface{}) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("id = ?", id)

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
func (dao *TArtifactDAO) UpdateByIdWithMapAndCondition(id uint64, updatedMap map[string]interface{}, conditionMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("id = ?", id)

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
func (dao *TArtifactDAO) DeleteById(id uint64) error {
	return dao.db.Where("id = ?", id).Delete(&entity.TArtifact{}).Error
}

// ==================== 唯一索引 uk_artifactId 方法 ====================

// SelectByArtifactId 根据唯一索引uk_artifactId查询单条记录
// 参数:
//   - artifactId: 产物ID
//
// 返回:
//   - *entity.TArtifact: 查询结果，如果不存在返回nil
//   - error: 错误信息
func (dao *TArtifactDAO) SelectByArtifactId(artifactId string) (*entity.TArtifact, error) {
	var resultBean entity.TArtifact
	err := dao.db.Where("artifactId = ?", artifactId).First(&resultBean).Error
	if err != nil {
		return nil, err
	}
	return &resultBean, nil
}

// UpdateByArtifactId 根据唯一索引uk_artifactId更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - artifactId: 产物ID
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 只更新非零值字段，零值字段会被忽略
func (dao *TArtifactDAO) UpdateByArtifactId(poBean *entity.TArtifact, artifactId string) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	return dao.db.Model(&entity.TArtifact{}).Where("artifactId = ?", artifactId).Updates(poBean).Error
}

// UpdateByArtifactIdWithMap 根据唯一索引uk_artifactId使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - artifactId: 产物ID
//   - updatedMap: 要更新的字段Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 使用 map 可以显式指定要更新的字段，包括零值字段
//   - 只更新 map 中指定的字段，未指定的字段保持不变
func (dao *TArtifactDAO) UpdateByArtifactIdWithMap(artifactId string, updatedMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	return dao.db.Model(&entity.TArtifact{}).Where("artifactId = ?", artifactId).Updates(updatedMap).Error
}

// UpdateByArtifactIdWithCondition 根据唯一索引uk_artifactId和额外条件更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - artifactId: 产物ID
//   - conditionMap: 额外的查询条件Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 只更新非零值字段，零值字段会被忽略
//   - 适用场景: 需要在唯一键基础上增加额外的更新条件
func (dao *TArtifactDAO) UpdateByArtifactIdWithCondition(poBean *entity.TArtifact, artifactId string, conditionMap map[string]interface{}) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("artifactId = ?", artifactId)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(poBean).Error
}

// UpdateByArtifactIdWithMapAndCondition 根据唯一索引uk_artifactId和额外条件使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - artifactId: 产物ID
//   - updatedMap: 要更新的字段Map
//   - conditionMap: 额外的查询条件Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 使用 map 可以显式指定要更新的字段，包括零值字段
//   - 提供最灵活的更新控制方式
func (dao *TArtifactDAO) UpdateByArtifactIdWithMapAndCondition(artifactId string, updatedMap map[string]interface{}, conditionMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("artifactId = ?", artifactId)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(updatedMap).Error
}

// DeleteByArtifactId 根据唯一索引uk_artifactId删除
// 参数:
//   - artifactId: 产物ID
//
// 返回:
//   - error: 错误信息
func (dao *TArtifactDAO) DeleteByArtifactId(artifactId string) error {
	return dao.db.Where("artifactId = ?", artifactId).Delete(&entity.TArtifact{}).Error
}

// ==================== 普通索引 idx_step 方法 ====================

// SelectByStep 根据索引idx_step查询列表
// 参数:
//   - step: 大的步骤点
//
// 返回:
//   - []*entity.TArtifact: 查询结果列表
//   - error: 错误信息
//
// 说明:
//   - 该索引不是唯一索引，可能返回多条记录
func (dao *TArtifactDAO) SelectByStep(step int) ([]*entity.TArtifact, error) {
	var resultList []*entity.TArtifact
	err := dao.db.Where("step = ?", step).Find(&resultList).Error
	return resultList, err
}

// SelectByStepList 根据索引idx_step批量查询列表
// 参数:
//   - stepList: 大的步骤点列表
//
// 返回:
//   - []*entity.TArtifact: 查询结果列表
//   - error: 错误信息
func (dao *TArtifactDAO) SelectByStepList(stepList []int) ([]*entity.TArtifact, error) {
	if len(stepList) == 0 {
		return []*entity.TArtifact{}, nil
	}
	var resultList []*entity.TArtifact
	err := dao.db.Where("step IN ?", stepList).Find(&resultList).Error
	return resultList, err
}

// UpdateByStep 根据索引idx_step更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - step: 大的步骤点
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 只更新非零值字段，零值字段会被忽略
//   - 注意: step 不是唯一键，可能会更新多条记录
func (dao *TArtifactDAO) UpdateByStep(poBean *entity.TArtifact, step int) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	return dao.db.Model(&entity.TArtifact{}).Where("step = ?", step).Updates(poBean).Error
}

// UpdateByStepWithMap 根据索引idx_step使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - step: 大的步骤点
//   - updatedMap: 要更新的字段Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 使用 map 可以显式指定要更新的字段，包括零值字段
//   - 只更新 map 中指定的字段，未指定的字段保持不变
//   - 注意: step 不是唯一键，可能会更新多条记录
func (dao *TArtifactDAO) UpdateByStepWithMap(step int, updatedMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	return dao.db.Model(&entity.TArtifact{}).Where("step = ?", step).Updates(updatedMap).Error
}

// UpdateByStepWithCondition 根据索引idx_step和额外条件更新（不会用零值覆盖）
// 参数:
//   - poBean: 包含更新数据的PO对象
//   - step: 大的步骤点
//   - conditionMap: 额外的查询条件Map
//
// 返回:
//   - error: 错误信息
//
// 行为说明:
//   - 只更新非零值字段，零值字段会被忽略
//   - 适用场景: 需要在 step 基础上增加额外的更新条件，缩小更新范围
//   - 注意: 可能会更新多条记录
func (dao *TArtifactDAO) UpdateByStepWithCondition(poBean *entity.TArtifact, step int, conditionMap map[string]interface{}) error {
	if poBean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("step = ?", step)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(poBean).Error
}

// UpdateByStepWithMapAndCondition 根据索引idx_step和额外条件使用Map更新指定字段（可以用零值覆盖）
// 参数:
//   - step: 大的步骤点
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
func (dao *TArtifactDAO) UpdateByStepWithMapAndCondition(step int, updatedMap map[string]interface{}, conditionMap map[string]interface{}) error {
	if len(updatedMap) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("step = ?", step)

	// 应用额外的条件
	for key, value := range conditionMap {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(updatedMap).Error
}

// DeleteByStep 根据索引idx_step删除
// 参数:
//   - step: 大的步骤点
//
// 返回:
//   - error: 错误信息
//
// 说明:
//   - 注意: step 不是唯一键，可能会删除多条记录
func (dao *TArtifactDAO) DeleteByStep(step int) error {
	return dao.db.Where("step = ?", step).Delete(&entity.TArtifact{}).Error
}

// ==================== 辅助方法 ====================

// getValidOrderByFields 获取允许排序的字段白名单
// 返回:
//   - map[string]bool: 字段白名单，key为字段名，value为true表示允许排序
func (dao *TArtifactDAO) getValidOrderByFields() map[string]bool {
	return map[string]bool{
		"id":           true,
		"artifactId":   true,
		"artifactName": true,
		"sessionId":    true,
		"step":         true,
		"subStep":      true,
		"content":      true,
		"version":      true,
		"createTime":   true,
		"updateTime":   true,
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
func (dao *TArtifactDAO) isValidOrderBy(orderBy string) bool {
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
