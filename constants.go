package log

// LevelType is the standard Level type iota
type LevelType int

const (
	// LogWTF should only be used for debugging
	LogWTF LevelType = iota
	// LogFatal indicates Critical level logging
	LogFatal
	// LogError indicates Error level logging
	LogError
	// LogWarning indicates Warning level logging
	LogWarning
	// LogDebug indicates Debug level logging
	LogDebug
	// LogVerbose indicates Verbose level logging
	LogVerbose
)

const (
	consoleID = "CONSOLE"

	// ConnectMessage will be sent to any listeners that get registered
	ConnectMessage string = "Listening to log server"

	// DisconnectMessage will be sent to any listeners that get de-registered
	DisconnectMessage string = "Disconnected from log server"
)

func (l LevelType) toString() string {
	levels := []string{
		"WTF      ",
		"FATAL    ",
		"ERROR    ",
		"WARNING  ",
		"DEBUG    ",
		"VERBOSE  ",
	}
	return levels[l]
}

// ToLogLevel will return a Level from an int.
// Invalid values will return LogFatal
func ToLogLevel(i int) LevelType {
	levels := []LevelType{
		LogWTF,
		LogFatal,
		LogError,
		LogWarning,
		LogDebug,
		LogVerbose,
	}

	if i >= 0 && i < len(levels) {
		return levels[i]
	}
	return LogFatal
}
