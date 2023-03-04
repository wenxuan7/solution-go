package log

import "go.uber.org/zap"

var Sugar *zap.SugaredLogger

func InitSugaredLogger() {
	var (
		err    error
		logger *zap.Logger
	)
	if logger, err = zap.NewProduction(); err != nil {
		panic(err)
	}

	Sugar = logger.Sugar()
}

// SugaredLoggerSync choosing logger
func SugaredLoggerSync() {
	var err error
	if err = Sugar.Sync(); err != nil {
		panic(err)
	}
}
