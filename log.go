package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	logLevel     LevelType = 0
	logListeners map[string]chan string
)

// SetLogLevel will determine what log entries are emitted
func SetLogLevel(level LevelType) {
	logLevel = level
}

// EnableConsoleLogging should be called in main.init() output all log lines to std.out
func EnableConsoleLogging() {
	logChan := RegisterListener(consoleID, nil)
	go func() {
		for {
			message := <-logChan
			if message == DisconnectMessage {
				return
			}
			log.Printf(message)
		}
	}()
}

// RegisterListener will start broadcasting all log entries to the provided channel
func RegisterListener(id string, messageChannel chan string) chan string {
	if logListeners == nil {
		logListeners = make(map[string]chan string)
	}

	if messageChannel == nil {
		messageChannel = make(chan string)
	}
	logListeners[id] = messageChannel
	go func() {
		for {
			select {
			case logListeners[id] <- ConnectMessage:
				return
			}
		}
	}()
	return messageChannel
}

// DeregisterListener remove the provided channel from the list of listeners
func DeregisterListener(id string) {
	if target := logListeners[id]; target != nil {
		logListeners[id] <- DisconnectMessage
	}
	delete(logListeners, id)
	return
}

// End will deregister all listeners
func End() {
	for id := range logListeners {
		fmt.Println(id)
		DeregisterListener(id)
	}
	logListeners = nil
}

// WTF should only be used for debugging and will ALWAYS generate a log entry
// despite the logging level
func WTF(msg string, args ...interface{}) {
	logImpl(LogWTF, msg, args...)
}

// Fatal will create a critical log entry and close the app with code 1
func Fatal(msg string, args ...interface{}) {
	logImpl(LogFatal, msg, args...)
	fmt.Printf(msg, args...)
	os.Exit(1)
}

// Error will create a error log entry and panic
func Error(msg string, args ...interface{}) {
	logImpl(LogError, msg, args...)
	panic(args)
}

// Warning will create a warning log entry
func Warning(msg string, args ...interface{}) {
	logImpl(LogWarning, msg, args...)
}

// Debug will create a debug log entry
func Debug(msg string, args ...interface{}) {
	logImpl(LogDebug, msg, args...)
}

// Verbose will create a verbose log entry
func Verbose(msg string, args ...interface{}) {
	logImpl(LogVerbose, msg, args...)
}

//Log provides a single point to control logging
func logImpl(l LevelType, msg string, args ...interface{}) {

	if logLevel < l {
		return
	}

	switch l {
	case LogFatal:
		fallthrough
	case LogError:
		fallthrough
	case LogWarning:
		fallthrough
	case LogDebug:
		fallthrough
	case LogWTF:
		fallthrough
	case LogVerbose:
		msg := fmt.Sprintf(l.toString()+getTag()+"\t"+msg+"\t", args...)
		notifyListeners(msg)
	}
}

func getTag() string {
	_, path, line, _ := runtime.Caller(3)

	rootPath, _ := os.Getwd()
	rootSegments := strings.Split(rootPath, string(os.PathSeparator))
	rootDir := rootSegments[len(rootSegments)-1]
	segments := strings.Split(path, rootDir+"/")
	file := segments[len(segments)-1]
	return fmt.Sprintf("%v:%v", file, line)
}

func notifyListeners(msg string) {
	for _, listener := range logListeners {
		select {
		case listener <- msg:
		}
	}
}
