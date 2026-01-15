package logger

type ILoggerConfig struct {
	LogPath    string // 日志文件路径
	MaxAge     int    // 日志文件最大保留天数
	MaxSize    int    // 单个日志文件最大大小(MB)
	MaxBackups int    // 最大保留的归档日志文件数
	Compress   bool   // 是否压缩归档日志
	Level      string // 日志级别: debug, info, warn, error
}
