package zap

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志文件名格式
const logFileNameFormat = "2006-01-02.log"

// 按日期记录日志的Hook
type DateHook struct {
	lock        sync.Mutex
	file        *os.File
	day         int
	logPath     string
	serviceName string
}

// 实现WriteSyncer接口
func (h *DateHook) Write(p []byte) (n int, err error) {
	// 获取当前日期
	today := time.Now().Day()

	h.lock.Lock()
	defer h.lock.Unlock()

	// 如果日期发生变化，关闭旧文件，创建新文件
	if today != h.day {
		if h.file != nil {
			h.file.Close()
		}

		// 创建新文件
		filename := time.Now().Format(logFileNameFormat)
		logFilePath := fmt.Sprintf("%s/%s_%s", h.logPath, h.serviceName, filename)

		// 创建目录
		err := os.MkdirAll(h.logPath, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
			return 0, err
		}

		h.file, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log file: %v\n", err)
			return 0, err
		}

		h.day = today
	}

	// 写入日志到文件
	return h.file.Write(p)
}

func NewLogger(logPath, serviceName string, logLevel zapcore.Level) (*zap.Logger, error) {
	// 创建日志文件的Hook
	hook := &DateHook{
		logPath:     logPath,
		serviceName: serviceName,
	}

	// 设置日志输出格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 创建日志核心对象
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(hook),
		logLevel,
	)

	// 创建Logger对象
	logger := zap.New(core, zap.AddCaller())

	return logger, nil
}