package gormzap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Logger logger for gorm2
type Logger struct {
	log *zap.Logger
	logger.Config
	customFields []func(ctx context.Context) zap.Field
}

// Option logger/recover option
type Option func(l *Logger)

// WithCustomFields optional custom field
func WithCustomFields(fields ...func(ctx context.Context) zap.Field) Option {
	return func(l *Logger) {
		l.customFields = fields
	}
}

// WithConfig optional custom logger.Config
func WithConfig(cfg logger.Config) Option {
	return func(l *Logger) {
		l.Config = cfg
	}
}

// SetGormDBLogger set db logger
func SetGormDBLogger(db *gorm.DB, l logger.Interface) {
	db.Logger = l
}

// New logger form gorm2
func New(zapLogger *zap.Logger, opts ...Option) logger.Interface {
	l := &Logger{
		log: zapLogger,
		Config: logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Warn,
		},
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// LogMode log mode
func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.log.Sugar().Debugf(msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
	}
}

// Warn print warn messages
func (l Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.log.Sugar().Warnf(msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
	}
}

// Error print error messages
func (l Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.log.Sugar().Errorf(msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
	}
}

// Trace print sql message
func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	fields := make([]zap.Field, 0, 6+len(l.customFields))
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		for _, customField := range l.customFields {
			fields = append(fields, customField(ctx))
		}
		fields = append(fields,
			zap.Error(err),
			zap.String("file", utils.FileWithLineNum()),
			zap.Duration("latency", elapsed),
		)

		sql, rows := fc()
		if rows == -1 {
			fields = append(fields, zap.String("rows", "-"))
		} else {
			fields = append(fields, zap.Int64("rows", rows))
		}
		fields = append(fields, zap.String("sql", sql))
		l.log.Error("trace", fields...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		for _, customField := range l.customFields {
			fields = append(fields, customField(ctx))
		}
		fields = append(fields,
			zap.Error(err),
			zap.String("file", utils.FileWithLineNum()),
			zap.String("slow!!!", fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)),
			zap.Duration("latency", elapsed),
		)

		sql, rows := fc()
		if rows == -1 {
			fields = append(fields, zap.String("rows", "-"))
		} else {
			fields = append(fields, zap.Int64("rows", rows))
		}
		fields = append(fields, zap.String("sql", sql))
		l.log.Warn("trace", fields...)
	case l.LogLevel == logger.Info:
		for _, customField := range l.customFields {
			fields = append(fields, customField(ctx))
		}
		fields = append(fields,
			zap.Error(err),
			zap.String("file", utils.FileWithLineNum()),
			zap.Duration("latency", elapsed),
		)

		sql, rows := fc()
		if rows == -1 {
			fields = append(fields, zap.String("rows", "-"))
		} else {
			fields = append(fields, zap.Int64("rows", rows))
		}
		fields = append(fields, zap.String("sql", sql))
		l.log.Info("trace", fields...)
	}
}

// Immutable custom immutable field
func Immutable(key string, value interface{}) func(ctx context.Context) zap.Field {
	field := zap.Any(key, value)
	return func(ctx context.Context) zap.Field { return field }
}
