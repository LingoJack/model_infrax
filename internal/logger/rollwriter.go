package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Option 配置选项函数类型
type Option func(*config)

// config
type config struct {
	maxSize    int  // 单个日志文件最大大小(MB)
	maxAge     int  // 日志文件最大保留天数
	maxBackups int  // 最大保留的归档日志文件数
	compress   bool // 是否压缩归档日志
}

// WithMaxSize 设置单个日志文件最大大小(MB)
func WithMaxSize(maxSize int) Option {
	return func(c *config) {
		c.maxSize = maxSize
	}
}

// WithMaxAge 设置日志文件最大保留天数
func WithMaxAge(maxAge int) Option {
	return func(c *config) {
		c.maxAge = maxAge
	}
}

// WithMaxBackups 设置最大保留的归档日志文件数
func WithMaxBackups(maxBackups int) Option {
	return func(c *config) {
		c.maxBackups = maxBackups
	}
}

// WithCompress 设置是否压缩归档日志
func WithCompress(compress bool) Option {
	return func(c *config) {
		c.compress = compress
	}
}

// RollWriter 日志滚动写入器
type RollWriter struct {
	filename   string // 当前日志文件路径
	maxSize    int64  // 单个日志文件最大大小(bytes)
	maxAge     int    // 日志文件最大保留天数
	maxBackups int    // 最大保留的归档日志文件数
	compress   bool   // 是否压缩归档日志

	currentSize int64    // 当前文件大小
	file        *os.File // 当前文件句柄

	mu sync.Mutex // 互斥锁，保证并发安全
}

// NewRollWriter 创建新的 RollWriter 实例
func NewRollWriter(filename string, options []Option) (*RollWriter, error) {
	// 默认配置
	cfg := &config{
		maxSize:    100,   // 默认100MB
		maxAge:     30,    // 默认保留30天
		maxBackups: 10,    // 默认保留10个文件
		compress:   false, // 默认不压缩
	}

	// 应用配置选项
	for _, opt := range options {
		opt(cfg)
	}

	// 确保日志目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	rw := &RollWriter{
		filename:   filename,
		maxSize:    int64(cfg.maxSize) * 1024 * 1024, // 转换为bytes
		maxAge:     cfg.maxAge,
		maxBackups: cfg.maxBackups,
		compress:   cfg.compress,
	}

	// 打开或创建日志文件
	if err := rw.open(); err != nil {
		return nil, err
	}

	return rw, nil
}

// open 打开或创建日志文件
func (rw *RollWriter) open() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	// 如果已有打开的文件，先关闭
	if rw.file != nil {
		rw.file.Close()
	}

	// 检查文件是否存在，如果存在则获取当前大小
	info, err := os.Stat(rw.filename)
	if err == nil {
		rw.currentSize = info.Size()
	} else if os.IsNotExist(err) {
		rw.currentSize = 0
	} else {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 打开文件，如果不存在则创建
	file, err := os.OpenFile(rw.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	rw.file = file
	return nil
}

// Write 实现 io.Writer 接口
func (rw *RollWriter) Write(p []byte) (n int, err error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	// 检查是否需要滚动文件
	if rw.needRotate() {
		if err := rw.rotate(); err != nil {
			return 0, err
		}
	}

	// 写入数据
	n, err = rw.file.Write(p)
	if err != nil {
		return n, err
	}

	rw.currentSize += int64(n)
	return n, nil
}

// needRotate 检查是否需要滚动文件
func (rw *RollWriter) needRotate() bool {
	return rw.currentSize >= rw.maxSize
}

// rotate 滚动日志文件
func (rw *RollWriter) rotate() error {
	// 关闭当前文件
	if rw.file != nil {
		rw.file.Close()
	}

	// 生成归档文件名
	baseName := filepath.Base(rw.filename)
	dirName := filepath.Dir(rw.filename)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)

	// 找到下一个可用的序号
	maxIndex := 0
	files, _ := filepath.Glob(filepath.Join(dirName, nameWithoutExt+"-*"+ext))
	for _, file := range files {
		parts := strings.Split(filepath.Base(file), "-")
		if len(parts) >= 2 {
			var index int
			fmt.Sscanf(parts[1], "%d", &index)
			if index > maxIndex {
				maxIndex = index
			}
		}
	}

	// 归档当前文件
	archiveName := fmt.Sprintf("%s-%d%s", nameWithoutExt, maxIndex+1, ext)
	archivePath := filepath.Join(dirName, archiveName)

	if err := os.Rename(rw.filename, archivePath); err != nil {
		return fmt.Errorf("归档日志文件失败: %v", err)
	}

	// 如果启用压缩，压缩归档文件
	if rw.compress {
		go rw.compressFile(archivePath)
	}

	// 清理过期的归档文件
	go rw.cleanOldFiles()

	// 重新创建新的日志文件
	file, err := os.OpenFile(rw.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("创建新日志文件失败: %v", err)
	}

	rw.file = file
	rw.currentSize = 0

	return nil
}

// compressFile 压缩文件（简化实现，实际可使用 gzip）
func (rw *RollWriter) compressFile(filepath string) {
	os.Rename(filepath, filepath+".gz")
}

// cleanOldFiles 清理过期的归档文件
func (rw *RollWriter) cleanOldFiles() {
	dir := filepath.Dir(rw.filename)
	baseName := filepath.Base(rw.filename)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)

	// 获取所有归档文件
	pattern := filepath.Join(dir, nameWithoutExt+"-*"+ext)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	// 按修改时间排序
	type fileInfo struct {
		path    string
		modTime time.Time
	}

	var fileInfos []fileInfo
	for _, file := range files {
		if info, err := os.Stat(file); err == nil {
			fileInfos = append(fileInfos, fileInfo{
				path:    file,
				modTime: info.ModTime(),
			})
		}
	}

	// 按修改时间降序排序（最新的在前）
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].modTime.After(fileInfos[j].modTime)
	})

	// 删除超过 maxBackups 数量的文件
	if len(fileInfos) > rw.maxBackups {
		for i := rw.maxBackups; i < len(fileInfos); i++ {
			os.Remove(fileInfos[i].path)
		}
	}

	// 删除超过 maxAge 天数的文件
	cutoff := time.Now().AddDate(0, 0, -rw.maxAge)
	for _, f := range fileInfos {
		if f.modTime.Before(cutoff) {
			os.Remove(f.path)
		}
	}
}

func (rw *RollWriter) Close() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.file != nil {
		return rw.file.Close()
	}
	return nil
}

func (rw *RollWriter) Sync() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.file != nil {
		return rw.file.Sync()
	}
	return nil
}
