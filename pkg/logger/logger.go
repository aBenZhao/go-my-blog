package logger

import (
	"go-my-blog/config/priority_config"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2" // 日志切割工具（生产环境用）
)

// 全局日志对象：整个项目通过此对象打印日志
var log *zap.Logger

// Init 初始化日志：根据配置（开发/生产环境）创建不同的日志实例
func Init() {
	// 从全局配置中获取日志相关参数（如级别、输出方式）
	logConf := priority_config.PriorityConf.Log

	// 1. 确定日志级别（Debug < Info < Warn < Error < Fatal）
	level := getZapLevel(logConf.Level)

	// 2. 确定日志输出介质（开发环境：控制台；生产环境：文件）
	var core zapcore.Core
	if priority_config.PriorityConf.Gin.Debug { // 开发环境（Gin 调试模式开启）
		core = newConsoleCore(level)
	} else { // 生产环境
		core = newFileCore(level, logConf.FilePath)
	}

	// 3. 创建日志实例：开发环境添加调用栈信息，生产环境简化
	if priority_config.PriorityConf.Gin.Debug {
		// 开发环境：打印文件名+行号，Debug 级别以上打印调用栈
		log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DebugLevel))
	} else {
		// 生产环境：仅打印文件名+行号，Error 级别以上打印调用栈（减少性能损耗）
		log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	}

	// 4. 替换 zap 全局日志（可选，方便其他包直接用 zap.L() 调用）
	zap.ReplaceGlobals(log)

	Info("日志初始化成功", zap.String("log_level", logConf.Level), zap.String("env", getEnvName()))
}

// -------------- 以下是日志对外暴露的打印接口（项目中直接调用）--------------
// Debug 是一个封装了zap日志库的调试级别日志记录函数
// 它提供了一个简化的接口来记录带有字段的调试信息
//
// 参数:
//
//	msg: 日志消息内容，字符串类型
//	fields: 可变参数，用于传递额外的字段信息，使用zap.Field类型
//	       这些字段可以提供结构化的上下文信息，如键值对
//
// 函数功能:
//
//	调用底层zap日志库的Debug方法，将传入的消息和字段记录到日志中
//	这是zap日志库的一个简单包装，保持了原有功能的同时提供了更简洁的调用方式
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info 是一个封装了日志记录功能的函数，用于记录带有额外字段的普通信息
// 它接收一个消息字符串和可变数量的zap.Field字段参数
//
// 参数:
//
//	msg: string类型，要记录的消息内容
//	fields: 可变参数，类型为zap.Field，用于记录额外的结构化日志字段
//
// 该函数内部调用了log.Info方法，将传入的消息和字段传递给底层的日志记录器
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 是一个封装了日志警告功能的函数
// 它接收一个消息字符串和可选的字段参数
//
// @param msg: 警告消息的文本内容
// @param fields: 可变参数，用于传递额外的日志字段，使用 zap.Field 类型
func Warn(msg string, fields ...zap.Field) {
	// 调用底层日志库的 Warn 方法记录警告日志
	// 将传入的消息和字段参数直接传递给底层实现
	log.Warn(msg, fields...)
}

// Error 是一个封装好的错误日志记录函数
// 它接收一个错误消息字符串和可选的字段列表
// 参数:
//
//	msg: 错误消息的字符串内容
//	fields: 可变参数，用于传递额外的错误上下文字段
func Error(msg string, fields ...zap.Field) {
	// 调用底层日志库的Error方法记录错误日志
	// 传入错误消息和可选的字段
	log.Error(msg, fields...)
}

// Fatal 是一个日志记录函数，用于记录致命级别的日志信息
// 当调用此函数时，程序会终止执行
//
// 参数:
//
//	msg: 日志消息字符串，用于描述错误或异常情况
//	fields: 可变参数，使用zap.Field类型，用于记录额外的结构化日志字段
//
// 该函数内部调用了log.Fatal方法，传入消息和字段后终止程序
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Sync 同步日志缓冲区：程序退出前调用，确保日志全部写入（如文件日志）
// Sync 函数用于同步日志
// 该函数会调用日志的 Sync 方法，确保所有缓冲的日志条目都被写入到底层存储
// 返回值:
//   - error: 如果同步过程中发生错误，则返回相应的错误信息；否则返回 nil
func Sync() error {
	return log.Sync()
}

// -------------- 以下是内部辅助函数（不对外暴露）--------------

// getZapLevel 将配置中的日志级别（字符串）转为 zap 库的级别类型
func getZapLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		// 配置错误时，默认用 Info 级别
		Warn("日志级别配置无效，默认使用 info 级别", zap.String("invalid_level", levelStr))
		return zapcore.InfoLevel
	}
}

// newConsoleCore 创建开发环境的控制台日志核心：彩色输出、详细时间格式
func newConsoleCore(level zapcore.Level) zapcore.Core {
	// 控制台编码器：开发环境显示彩色级别、详细时间（如 2024-06-01 15:30:00.123）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",                         // 显示文件名+行号（如 pkg/logger/logger.go:20）
		EncodeLevel:  zapcore.CapitalColorLevelEncoder, // 彩色级别（如 INFO 绿色，ERROR 红色）
		EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeCaller: zapcore.ShortCallerEncoder, // 短路径（长路径用 FullCallerEncoder）
		LineEnding:   zapcore.DefaultLineEnding,
	}

	// 输出到控制台（标准输出流 os.Stdout）
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)
}

// newFileCore 创建生产环境的文件日志核心：JSON 格式、日志切割
func newFileCore(level zapcore.Level, logPath string) zapcore.Core {
	// 确保日志目录存在（如 /var/log/blog/，不存在则创建）
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		Warn("创建日志目录失败，默认输出到当前目录", zap.Error(err), zap.String("target_path", logPath))
		logPath = "blog-backend.log" // 降级到当前目录
	}

	// 日志切割配置（生产环境必备，避免日志文件过大）
	writer := &lumberjack.Logger{
		Filename:   logPath, // 日志文件路径（如 /var/log/blog/backend.log）
		MaxSize:    100,     // 单个日志文件最大 100MB
		MaxBackups: 7,       // 保留 7 个备份文件（超过自动删除）
		MaxAge:     30,      // 日志文件保留 30 天
		Compress:   true,    // 备份文件压缩（节省磁盘空间）
		LocalTime:  true,    // 用本地时间命名备份文件（如 backend-20240601.log）
	}

	// 文件编码器：生产环境用 JSON 格式（便于日志分析工具解析）、UTC 时间
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.CapitalLevelEncoder, // 无颜色（JSON 格式不需要）
		EncodeTime:   zapcore.ISO8601TimeEncoder,  // UTC 时间（如 2024-06-01T07:30:00.123Z）
		EncodeCaller: zapcore.FullCallerEncoder,   // 完整路径（便于定位问题）
		LineEnding:   zapcore.DefaultLineEnding,
	}

	// 输出到文件（带切割功能的 writer）
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		level,
	)
}

// getEnvName 获取当前环境名称（开发/生产），用于日志打印
func getEnvName() string {
	if priority_config.PriorityConf.Gin.Debug {
		return "dev"
	}
	return "prod"
}
