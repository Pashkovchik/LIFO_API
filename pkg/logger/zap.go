package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZapLogger(logLevel string) {
	level, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		log.Fatalf("logger - NewZapLogger - zap.ParseAtomicLevel: %s", err.Error())
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.Level = level
	zapConfig.DisableStacktrace = true

	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("logger - NewZapLogger - zapConfig.Build: %s", err.Error())
	}

	zap.ReplaceGlobals(zapLog)
	zap.S().Infof("Log Level: [%s]", level.Level().CapitalString())

}

//
//func ChangeLogLevel(logLevel string) error {
//	level, err := zap.ParseAtomicLevel(logLevel)
//	if err != nil {
//		zap.S().Infof("logger - ChangeLogLevel - zap.ParseAtomicLevel: %s", err.Error())
//
//		return err
//	}
//
//	newLvl := zap.NewAtomicLevel()
//	newLvl.SetLevel(level.Level())
//	zap.S().Infof("Zap Level updated to: [%s]", level.Level().CapitalString())
//
//	return nil
//}
