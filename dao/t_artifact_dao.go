package dao

import (
	"fmt"
	"model_infrax/output/model/entity"
	"model_infrax/output/model/query"
	"strings"

	"gorm.io/gorm"
)

// TArtifactDAO TArtifactDAO的实现
type TArtifactDAO struct {
	db *gorm.DB
}

// NewTArtifactDAO 创建TArtifactDAO实例
func NewTArtifactDAO(db *gorm.DB) *TArtifactDAO {
	return &TArtifactDAO{db: db}
}

// buildQueryCondition 构建查询条件
func (dao *TArtifactDAO) buildQueryCondition(db *gorm.DB, dto *query.TArtifactDTO) *gorm.DB {
	if dto == nil {
		return db
	}

	// 基础查询条件
	if dto.Id != 0 {
		db = db.Where("id = ?", dto.Id)
	}
	if dto.ArtifactId != "" {
		db = db.Where("artifactId = ?", dto.ArtifactId)
	}
	if dto.ArtifactName != "" {
		db = db.Where("artifactName = ?", dto.ArtifactName)
	}
	if dto.SessionId != "" {
		db = db.Where("sessionId = ?", dto.SessionId)
	}
	if dto.Step != 0 {
		db = db.Where("step = ?", dto.Step)
	}
	if dto.SubStep != "" {
		db = db.Where("subStep = ?", dto.SubStep)
	}
	if dto.Content != nil && *dto.Content != "" {
		db = db.Where("content = ?", *dto.Content)
	}
	if dto.Version != nil && *dto.Version != "" {
		db = db.Where("version = ?", *dto.Version)
	}
	if !dto.CreateTime.IsZero() {
		db = db.Where("createTime = ?", dto.CreateTime)
	}
	if !dto.UpdateTime.IsZero() {
		db = db.Where("updateTime = ?", dto.UpdateTime)
	}

	// 扩展查询条件（模糊查询）
	if dto.ArtifactIdFuzzy != "" {
		db = db.Where("artifactId LIKE ?", "%"+dto.ArtifactIdFuzzy+"%")
	}
	if dto.ArtifactNameFuzzy != "" {
		db = db.Where("artifactName LIKE ?", "%"+dto.ArtifactNameFuzzy+"%")
	}
	if dto.SessionIdFuzzy != "" {
		db = db.Where("sessionId LIKE ?", "%"+dto.SessionIdFuzzy+"%")
	}
	if dto.SubStepFuzzy != "" {
		db = db.Where("subStep LIKE ?", "%"+dto.SubStepFuzzy+"%")
	}
	if dto.ContentFuzzy != nil && *dto.ContentFuzzy != "" {
		db = db.Where("content LIKE ?", "%"+*dto.ContentFuzzy+"%")
	}
	if dto.VersionFuzzy != nil && *dto.VersionFuzzy != "" {
		db = db.Where("version LIKE ?", "%"+*dto.VersionFuzzy+"%")
	}

	// 日期范围查询
	if !dto.CreateTimeStart.IsZero() {
		db = db.Where("createTime >= ?", dto.CreateTimeStart)
	}
	if !dto.CreateTimeEnd.IsZero() {
		db = db.Where("createTime < DATE_ADD(?, INTERVAL 1 DAY)", dto.CreateTimeEnd)
	}
	if !dto.UpdateTimeStart.IsZero() {
		db = db.Where("updateTime >= ?", dto.UpdateTimeStart)
	}
	if !dto.UpdateTimeEnd.IsZero() {
		db = db.Where("updateTime < DATE_ADD(?, INTERVAL 1 DAY)", dto.UpdateTimeEnd)
	}

	// 列表查询条件（IN 查询）
	// 注意：这里使用 IF 模式，可能有多个 IN 条件同时生效
	if len(dto.IdList) > 0 {
		db = db.Where("id IN ?", dto.IdList)
	}
	if len(dto.ArtifactIdList) > 0 {
		db = db.Where("artifactId IN ?", dto.ArtifactIdList)
	}
	if len(dto.ArtifactNameList) > 0 {
		db = db.Where("artifactName IN ?", dto.ArtifactNameList)
	}
	if len(dto.SessionIdList) > 0 {
		db = db.Where("sessionId IN ?", dto.SessionIdList)
	}
	if len(dto.StepList) > 0 {
		db = db.Where("step IN ?", dto.StepList)
	}
	if len(dto.SubStepList) > 0 {
		db = db.Where("subStep IN ?", dto.SubStepList)
	}
	if len(dto.ContentList) > 0 {
		db = db.Where("content IN ?", dto.ContentList)
	}
	if len(dto.VersionList) > 0 {
		db = db.Where("version IN ?", dto.VersionList)
	}
	if len(dto.CreateTimeList) > 0 {
		db = db.Where("createTime IN ?", dto.CreateTimeList)
	}
	if len(dto.UpdateTimeList) > 0 {
		db = db.Where("updateTime IN ?", dto.UpdateTimeList)
	}

	return db
}

// SelectList 查询列表
func (dao *TArtifactDAO) SelectList(dto *query.TArtifactDTO) ([]*entity.TArtifact, error) {
	var result []*entity.TArtifact
	db := dao.db.Model(&entity.TArtifact{})

	// 应用查询条件
	db = dao.buildQueryCondition(db, dto)

	// 排序
	if dto != nil && dto.OrderBy != "" {
		// 使用已有的安全验证函数
		if isValidOrderBy(dto.OrderBy) {
			db = db.Order(dto.OrderBy)
		}
	}

	// 分页（使用扁平化字段）
	if dto != nil && dto.PageSize > 0 {
		db = db.Offset(dto.Offset).Limit(dto.PageSize)
	}

	err := db.Find(&result).Error
	return result, err
}

// SelectCount 查询数量
func (dao *TArtifactDAO) SelectCount(dto *query.TArtifactDTO) (int64, error) {
	var count int64
	db := dao.db.Model(&entity.TArtifact{})

	// 应用查询条件
	db = dao.buildQueryCondition(db, dto)

	err := db.Count(&count).Error
	return count, err
}

// Insert 单行插入
func (dao *TArtifactDAO) Insert(bean *entity.TArtifact) error {
	if bean == nil {
		return fmt.Errorf("插入对象不能为空")
	}
	return dao.db.Create(bean).Error
}

// InsertOrUpdateNullable 插入或更新
// 行为说明：
//  1. 如果记录不存在（根据主键判断），则执行插入操作
//  2. 如果记录已存在，则执行全字段更新操作
//  3. **重要**：更新时会用传入对象的所有字段值覆盖数据库中的值，包括零值（nil、""、0、false等）
//     例如：如果 bean.Content = nil，会将数据库中的 content 字段更新为 NULL
//     例如：如果 bean.ArtifactName = ""，会将数据库中的 artifactName 字段更新为空字符串
//  4. 这种行为适用于需要"完整替换"记录的场景
//  5. 如果不希望零值覆盖数据库中的非零值，应使用 UpdateById 等方法（内部使用 Updates）
func (dao *TArtifactDAO) InsertOrUpdateNullable(bean *entity.TArtifact) error {
	if bean == nil {
		return fmt.Errorf("插入或更新对象不能为空")
	}
	// 使用 GORM 的 Save 方法：
	// - 根据主键判断记录是否存在
	// - 存在则更新所有字段（包括零值字段）
	// - 不存在则插入新记录
	return dao.db.Save(bean).Error
}

// InsertBatch 批量插入
func (dao *TArtifactDAO) InsertBatch(list []*entity.TArtifact) error {
	if len(list) == 0 {
		return fmt.Errorf("批量插入列表不能为空")
	}
	return dao.db.Create(&list).Error
}

// InsertOrUpdateBatchNullable  批量插入或更新
// 行为说明：
//  1. 对列表中的每条记录，根据主键判断是插入还是更新
//  2. 如果记录不存在，则执行插入操作
//  3. 如果记录已存在，则执行全字段更新操作
//  4. **重要**：更新时会用传入对象的所有字段值覆盖数据库中的值，包括零值（nil、""、0、false等）
//     这意味着如果某个字段在传入对象中为零值，会将数据库中对应字段更新为零值
//  5. 批量操作在一个事务中执行，要么全部成功，要么全部失败
//  6. 适用场景：需要完整替换多条记录的场景
//  7. 性能提示：批量操作比逐条调用 InsertOrUpdate 效率更高
//  8. 如果不希望零值覆盖，建议逐条调用 UpdateById 等方法
func (dao *TArtifactDAO) InsertOrUpdateBatchNullable(list []*entity.TArtifact) error {
	if len(list) == 0 {
		return fmt.Errorf("批量插入或更新列表不能为空")
	}
	// 使用 GORM 的 Save 方法批量保存：
	// - 对每条记录根据主键判断是插入还是更新
	// - 更新时会覆盖所有字段（包括零值字段）
	// - 在一个事务中执行，保证原子性
	return dao.db.Save(&list).Error
}

// SelectById 根据Id查询
func (dao *TArtifactDAO) SelectById(id uint64) (*entity.TArtifact, error) {
	var result entity.TArtifact
	err := dao.db.Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateById 根据Id更新
// 行为说明：
//  1. 根据指定的 id 更新记录
//  2. **重要**：只更新非零值字段，零值字段会被忽略，不会覆盖数据库中的值
//     例如：如果 bean.Content = nil，不会更新数据库中的 content 字段
//     例如：如果 bean.ArtifactName = ""，不会更新数据库中的 artifactName 字段
//  3. 这种行为适用于"部分更新"场景，保留数据库中未传入的字段值
//  4. 如果需要将某个字段更新为零值，应使用 UpdateByIdWithMap 方法显式指定
//  5. 与 InsertOrUpdate 的区别：InsertOrUpdate 会用零值覆盖，UpdateById 不会
func (dao *TArtifactDAO) UpdateById(bean *entity.TArtifact, id uint64) error {
	if bean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	// 使用 Updates 方法：
	// - 只更新结构体中的非零值字段
	// - 零值字段会被忽略，保留数据库中的原值
	// - 适合部分更新场景
	return dao.db.Model(&entity.TArtifact{}).Where("id = ?", id).Updates(bean).Error
}

// UpdateByIdWithCondition 根据Id和额外条件更新
// 行为说明：
//  1. 根据指定的 id 和额外的条件更新记录
//  2. 只更新非零值字段，零值字段会被忽略
//  3. 适用场景：需要在主键基础上增加额外的更新条件，如乐观锁、状态检查等
//  4. 示例：condition["version"] = 1 可以实现乐观锁，只有版本号匹配才更新
func (dao *TArtifactDAO) UpdateByIdWithCondition(bean *entity.TArtifact, id uint64, condition map[string]interface{}) error {
	if bean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("id = ?", id)

	// 应用额外的条件
	for key, value := range condition {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(bean).Error
}

// UpdateByIdWithMapAndCondition 根据Id和额外条件使用Map更新指定字段
// 行为说明：
//  1. 根据指定的 id 和额外的条件更新记录
//  2. 使用 map 可以显式指定要更新的字段，包括零值字段
//  3. 适用场景：需要精确控制更新字段，并且需要额外的更新条件
//  4. 示例：updates["status"] = 0, condition["old_status"] = 1 实现状态流转控制
func (dao *TArtifactDAO) UpdateByIdWithMapAndCondition(id uint64, updates map[string]interface{}, condition map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("id = ?", id)

	// 应用额外的条件
	for key, value := range condition {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(updates).Error
}

// DeleteById 根据Id删除
func (dao *TArtifactDAO) DeleteById(id uint64) error {
	return dao.db.Where("id = ?", id).Delete(&entity.TArtifact{}).Error
}

// SelectByIdList 根据IdList列表查询
func (dao *TArtifactDAO) SelectByIdList(idList []uint64) ([]*entity.TArtifact, error) {
	if len(idList) == 0 {
		return []*entity.TArtifact{}, nil
	}
	var result []*entity.TArtifact
	err := dao.db.Where("id IN ?", idList).Find(&result).Error
	return result, err
}

// SelectByArtifactId 根据ArtifactId查询
func (dao *TArtifactDAO) SelectByArtifactId(artifactId string) (*entity.TArtifact, error) {
	var result entity.TArtifact
	err := dao.db.Where("artifactId = ?", artifactId).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateByArtifactId 根据ArtifactId更新
func (dao *TArtifactDAO) UpdateByArtifactId(bean *entity.TArtifact, artifactId string) error {
	if bean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	// 使用 Updates 方法，只更新非零值字段
	return dao.db.Model(&entity.TArtifact{}).Where("artifactId = ?", artifactId).Updates(bean).Error
}

// UpdateByArtifactIdWithMap 根据ArtifactId使用Map更新指定字段
// 行为说明：
//  1. 根据指定的 artifactId 更新记录
//  2. 使用 map 可以显式指定要更新的字段，包括零值字段
//  3. 只更新 map 中指定的字段，未指定的字段保持不变
func (dao *TArtifactDAO) UpdateByArtifactIdWithMap(artifactId string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	return dao.db.Model(&entity.TArtifact{}).Where("artifactId = ?", artifactId).Updates(updates).Error
}

// UpdateByArtifactIdWithCondition 根据ArtifactId和额外条件更新
// 行为说明：
//  1. 根据指定的 artifactId 和额外的条件更新记录
//  2. 只更新非零值字段，零值字段会被忽略
//  3. 适用场景：需要在唯一键基础上增加额外的更新条件
func (dao *TArtifactDAO) UpdateByArtifactIdWithCondition(bean *entity.TArtifact, artifactId string, condition map[string]interface{}) error {
	if bean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("artifactId = ?", artifactId)

	// 应用额外的条件
	for key, value := range condition {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(bean).Error
}

// UpdateByArtifactIdWithMapAndCondition 根据ArtifactId和额外条件使用Map更新指定字段
// 行为说明：
//  1. 根据指定的 artifactId 和额外的条件更新记录
//  2. 使用 map 可以显式指定要更新的字段，包括零值字段
//  3. 提供最灵活的更新控制方式
func (dao *TArtifactDAO) UpdateByArtifactIdWithMapAndCondition(artifactId string, updates map[string]interface{}, condition map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("artifactId = ?", artifactId)

	// 应用额外的条件
	for key, value := range condition {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(updates).Error
}

// DeleteByArtifactId 根据ArtifactId删除
func (dao *TArtifactDAO) DeleteByArtifactId(artifactId string) error {
	return dao.db.Where("artifactId = ?", artifactId).Delete(&entity.TArtifact{}).Error
}

// SelectByArtifactIdList 根据ArtifactIdList列表查询
func (dao *TArtifactDAO) SelectByArtifactIdList(artifactIdList []string) ([]*entity.TArtifact, error) {
	if len(artifactIdList) == 0 {
		return []*entity.TArtifact{}, nil
	}
	var result []*entity.TArtifact
	err := dao.db.Where("artifactId IN ?", artifactIdList).Find(&result).Error
	return result, err
}

// SelectBySessionIdAndVersion 根据SessionIdAndVersion查询
func (dao *TArtifactDAO) SelectBySessionIdAndVersion(sessionId string, version string) (*entity.TArtifact, error) {
	var result entity.TArtifact
	err := dao.db.Where("sessionId = ? AND version = ?", sessionId, version).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateBySessionIdAndVersion 根据SessionIdAndVersion更新
func (dao *TArtifactDAO) UpdateBySessionIdAndVersion(bean *entity.TArtifact, sessionId string, version string) error {
	if bean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	// 使用 Updates 方法，只更新非零值字段
	return dao.db.Model(&entity.TArtifact{}).Where("sessionId = ? AND version = ?", sessionId, version).Updates(bean).Error
}

// UpdateBySessionIdAndVersionWithMap 根据SessionIdAndVersion使用Map更新指定字段
// 行为说明：
//  1. 根据指定的 sessionId 和 version 更新记录
//  2. 使用 map 可以显式指定要更新的字段，包括零值字段
//  3. 只更新 map 中指定的字段，未指定的字段保持不变
func (dao *TArtifactDAO) UpdateBySessionIdAndVersionWithMap(sessionId string, version string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	return dao.db.Model(&entity.TArtifact{}).Where("sessionId = ? AND version = ?", sessionId, version).Updates(updates).Error
}

// UpdateBySessionIdAndVersionWithCondition 根据SessionIdAndVersion和额外条件更新
// 行为说明：
//  1. 根据指定的 sessionId、version 和额外的条件更新记录
//  2. 只更新非零值字段，零值字段会被忽略
//  3. 适用场景：需要在复合键基础上增加额外的更新条件
func (dao *TArtifactDAO) UpdateBySessionIdAndVersionWithCondition(bean *entity.TArtifact, sessionId string, version string, condition map[string]interface{}) error {
	if bean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("sessionId = ? AND version = ?", sessionId, version)

	// 应用额外的条件
	for key, value := range condition {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(bean).Error
}

// UpdateBySessionIdAndVersionWithMapAndCondition 根据SessionIdAndVersion和额外条件使用Map更新指定字段
// 行为说明：
//  1. 根据指定的 sessionId、version 和额外的条件更新记录
//  2. 使用 map 可以显式指定要更新的字段，包括零值字段
//  3. 提供最灵活的更新控制方式
func (dao *TArtifactDAO) UpdateBySessionIdAndVersionWithMapAndCondition(sessionId string, version string, updates map[string]interface{}, condition map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("sessionId = ? AND version = ?", sessionId, version)

	// 应用额外的条件
	for key, value := range condition {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(updates).Error
}

// DeleteBySessionIdAndVersion 根据SessionIdAndVersion删除
func (dao *TArtifactDAO) DeleteBySessionIdAndVersion(sessionId string, version string) error {
	return dao.db.Where("sessionId = ? AND version = ?", sessionId, version).Delete(&entity.TArtifact{}).Error
}

// SelectByStep 根据Step查询
func (dao *TArtifactDAO) SelectByStep(step int) ([]*entity.TArtifact, error) {
	var result []*entity.TArtifact
	err := dao.db.Where("step = ?", step).Find(&result).Error
	return result, err
}

// UpdateByStep 根据Step更新
func (dao *TArtifactDAO) UpdateByStep(bean *entity.TArtifact, step int) error {
	if bean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	// 使用 Updates 方法，只更新非零值字段
	return dao.db.Model(&entity.TArtifact{}).Where("step = ?", step).Updates(bean).Error
}

// UpdateByStepWithMap 根据Step使用Map更新指定字段
// 行为说明：
//  1. 根据指定的 step 更新记录（可能更新多条）
//  2. 使用 map 可以显式指定要更新的字段，包括零值字段
//  3. 只更新 map 中指定的字段，未指定的字段保持不变
//  4. 注意：step 不是唯一键，可能会更新多条记录
func (dao *TArtifactDAO) UpdateByStepWithMap(step int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	return dao.db.Model(&entity.TArtifact{}).Where("step = ?", step).Updates(updates).Error
}

// UpdateByStepWithCondition 根据Step和额外条件更新
// 行为说明：
//  1. 根据指定的 step 和额外的条件更新记录
//  2. 只更新非零值字段，零值字段会被忽略
//  3. 适用场景：需要在 step 基础上增加额外的更新条件，缩小更新范围
//  4. 注意：可能会更新多条记录
func (dao *TArtifactDAO) UpdateByStepWithCondition(bean *entity.TArtifact, step int, condition map[string]interface{}) error {
	if bean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("step = ?", step)

	// 应用额外的条件
	for key, value := range condition {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(bean).Error
}

// UpdateByStepWithMapAndCondition 根据Step和额外条件使用Map更新指定字段
// 行为说明：
//  1. 根据指定的 step 和额外的条件更新记录
//  2. 使用 map 可以显式指定要更新的字段，包括零值字段
//  3. 提供最灵活的更新控制方式
//  4. 注意：可能会更新多条记录
func (dao *TArtifactDAO) UpdateByStepWithMapAndCondition(step int, updates map[string]interface{}, condition map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	db := dao.db.Model(&entity.TArtifact{}).Where("step = ?", step)

	// 应用额外的条件
	for key, value := range condition {
		db = db.Where(key+" = ?", value)
	}

	return db.Updates(updates).Error
}

// DeleteByStep 根据Step删除
func (dao *TArtifactDAO) DeleteByStep(step int) error {
	return dao.db.Where("step = ?", step).Delete(&entity.TArtifact{}).Error
}

// SelectByStepList 根据StepList列表查询
func (dao *TArtifactDAO) SelectByStepList(stepList []int) ([]*entity.TArtifact, error) {
	if len(stepList) == 0 {
		return []*entity.TArtifact{}, nil
	}
	var result []*entity.TArtifact
	err := dao.db.Where("step IN ?", stepList).Find(&result).Error
	return result, err
}

// 以下是一些辅助方法，用于更灵活的查询

// SelectListWithPage 分页查询列表
func (dao *TArtifactDAO) SelectListWithPage(dto *query.TArtifactDTO, page, pageSize int) ([]*entity.TArtifact, int64, error) {
	var result []*entity.TArtifact
	var total int64

	db := dao.db.Model(&entity.TArtifact{})
	db = dao.buildQueryCondition(db, dto)

	// 先查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查询分页数据
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&result).Error
	return result, total, err
}

// UpdateByIdSelective 根据Id选择性更新（只更新非空字段）
func (dao *TArtifactDAO) UpdateByIdSelective(bean *entity.TArtifact, id uint64) error {
	if bean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	// 使用 Updates 方法，只更新非零值字段
	return dao.db.Model(&entity.TArtifact{}).Where("id = ?", id).Updates(bean).Error
}

// UpdateByIdWithMap 根据Id使用Map更新指定字段
// 行为说明：
//  1. 根据指定的 id 更新记录
//  2. 使用 map 可以显式指定要更新的字段，包括零值字段
//  3. **重要**：与 UpdateById 不同，使用 map 可以将字段更新为零值
//     例如：updates["content"] = nil 会将 content 字段更新为 NULL
//     例如：updates["artifactName"] = "" 会将 artifactName 字段更新为空字符串
//  4. 只更新 map 中指定的字段，未指定的字段保持不变
//  5. 适用场景：需要精确控制更新哪些字段，包括需要将某些字段设置为零值的场景
//  6. 使用建议：字段名必须与数据库列名一致（或使用 GORM 的字段映射名）
func (dao *TArtifactDAO) UpdateByIdWithMap(id uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("更新字段不能为空")
	}
	// 使用 Updates 方法配合 map：
	// - 可以显式更新零值字段
	// - 只更新 map 中指定的字段
	// - 提供最精确的字段更新控制
	return dao.db.Model(&entity.TArtifact{}).Where("id = ?", id).Updates(updates).Error
}

// MustUpdateById 必须更新一条记录，如果影响行数不为1则返回错误
func (dao *TArtifactDAO) MustUpdateById(bean *entity.TArtifact, id uint64) error {
	if bean == nil {
		return fmt.Errorf("更新对象不能为空")
	}
	result := dao.db.Model(&entity.TArtifact{}).Where("id = ?", id).Updates(bean)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("期望更新1条记录，实际更新了%d条", result.RowsAffected)
	}
	return nil
}

// SelectListWithOrder 带排序的查询列表
func (dao *TArtifactDAO) SelectListWithOrder(dto *query.TArtifactDTO, orderBy string) ([]*entity.TArtifact, error) {
	var result []*entity.TArtifact
	db := dao.db.Model(&entity.TArtifact{})

	// 应用查询条件
	db = dao.buildQueryCondition(db, dto)

	// 应用排序
	if orderBy != "" {
		// 防止SQL注入，验证orderBy字符串
		if isValidOrderBy(orderBy) {
			db = db.Order(orderBy)
		}
	}

	err := db.Find(&result).Error
	return result, err
}

// validOrderByFields 定义允许排序的字段白名单
// 只有在此列表中的字段才允许用于排序，防止SQL注入
var validOrderByFields = map[string]bool{
	"id":           true,
	"artifactId":   true,
	"artifactName": true,
	"sessionId":    true,
	"step":         true,
	"substep":      true,
	"content":      true,
	"version":      true,
	"createTime":   true,
	"updateTime":   true,
}

// isValidOrderBy 验证排序字符串是否安全（基于字段白名单）
// 支持格式：
//   - 单字段：id DESC
//   - 多字段：id DESC, createTime ASC
//
// 返回：
//   - true: 排序字符串合法且所有字段都在白名单中
//   - false: 排序字符串不合法或包含非白名单字段
func isValidOrderBy(orderBy string) bool {
	if orderBy == "" {
		return false
	}

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
			// 格式错误：必须是 "字段名" 或 "字段名 方向"
			return false
		}

		// 验证字段名是否在白名单中（不区分大小写）
		fieldName := strings.ToLower(tokens[0])
		if !validOrderByFields[fieldName] {
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
