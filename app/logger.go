package app

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Logger 全局日志实例
var Logger *zap.Logger

// 5.这个 CustomLoggerEncoder 会实现日志双写功能
type CustomLoggerEncoder struct {
	zapcore.Encoder
	warnFile  *os.File
	errorFile *os.File
	infoFile  *os.File
	date      string
}

func (cle CustomLoggerEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := cle.Encoder.EncodeEntry(ent, fields)
	if err != nil {
		return buf, err
	}

	today := time.Now().Format(DATE_FMT)

	// 6.日志时间分片,按照每天一个目录来做区分
	// 创建今天的日志保存目录
	if cle.date != today {
		os.MkdirAll(fmt.Sprintf("logs/%s", today), 0o777)
		cle.date = today
	}

	// 4.修改日志内容: 统一增加前缀
	logPrefix := Config.Get("logger.log_prefix").(string)
	logString := buf.String()
	buf.Reset()
	buf.AppendString(logPrefix + logString)

	// 6.日志等级分片
	switch ent.Level {
	case zapcore.WarnLevel:
		if cle.warnFile == nil {
			filePath := fmt.Sprintf("logs/%s/warn.log", today)
			cle.warnFile, _ = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
		}
		cle.warnFile.Write(buf.Bytes())
	case zapcore.ErrorLevel:
		if cle.errorFile == nil {
			filePath := fmt.Sprintf("logs/%s/error.log", today)
			cle.errorFile, _ = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
		}
		cle.errorFile.Write(buf.Bytes())
	default:
		if cle.infoFile == nil {
			filePath := fmt.Sprintf("logs/%s/info.log", today)
			cle.infoFile, _ = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
		}
		cle.infoFile.Write(buf.Bytes())
	}

	return buf, nil
}

func InitLogger() {
	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(DATETIME_FMT) // 2.格式化时间

	// 3.使用自定义的 Encoder, 输出 json
	encoder := &CustomLoggerEncoder{
		Encoder: zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
	}

	// 构建自定义的 Core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(getLoggerLevelFromConfig()), // 1.设置日志级别
	)

	// 直接赋值
	Logger = zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(Logger) // 7.全局日志实例设置
}

func getLoggerLevelFromConfig() zapcore.Level {
	configLogLevel := Config.Get("logger.log_level")
	switch configLogLevel {
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}
