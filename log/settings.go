package log

import "os"

const (
	// The key used in contexts to retrieve the logger that should be used
	ContextLogKey NmcContextKey = "nmc_logger"

	// Lumberjack settings
	MaxLogSize    int = 20
	MaxLogBackups int = 3
	MaxLogAge     int = 28

	logDirMode  os.FileMode = 0755
	logFileMode os.FileMode = 0644
)

type NmcContextKey string
