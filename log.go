package log

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

type level int

const (
	LevelWTF level = iota - 1
	LevelInfo
	LevelFatal
	LevelError
	LevelWarning
	LevelDebug
	LevelVerbose
)

var (
	levelStrings = map[level]string{
		LevelWTF:     "WTF",
		LevelInfo:    "Info",
		LevelFatal:   "Fatal",
		LevelError:   "Error",
		LevelWarning: "Warning",
		LevelDebug:   "Debug",
		LevelVerbose: "Verbose",
	}

	tagPattern = *regexp.MustCompile(`.*\/(.*)\.\(.*\)\.(.*)`)
)

func (l level) toString() string {
	return fmt.Sprintf("%-9s", levelStrings[l])
}

func GetLevelFromString(levelString string) level {
	for level, str := range levelStrings {
		if strings.EqualFold(levelString, str) {
			return level
		}
	}
	return LevelWTF
}

var (
	logLevel  level                  = LevelWarning
	toConsole bool                   = true
	listeners map[string]chan string = make(map[string]chan string)
	timeout   time.Duration          = 5 * time.Second
)

// SetLogLevel will set the level of logging verbosity. Default = LevelWarning
func SetLogLevel(newLevel level) {
	logLevel = newLevel
}

// ToConsole controls whether to log to console or not
func ToConsole(logToConsole bool) {
	toConsole = logToConsole
}

// SetTimeout sets the time until we stop trying to notify a listener and discard the log
func SetTimeout(newTimeout time.Duration) {
	timeout = newTimeout
}

// AddListener will start broadcasting all log entries to the provided channe
func AddListener(id string, messageChannel chan string) error {
	if messageChannel == nil {
		return errors.New("attempted to register a log receiver without a channel")
	}

	if listeners[id] != nil {
		return errors.New("attempted to register a log receiver without a unique identifier")
	}

	listeners[id] = messageChannel
	return nil
}

// RemoveListener removes the provided ID from the list of listeners
func RemoveListener(id string) {
	delete(listeners, id)
}

// WTF should only be used for debugging and will ALWAYS generate a log entry
// despite the logging level
func WTF(args ...interface{}) {
	logImpl(LevelWTF, args...)
}

// Info will create a log entry you intentionally want to display to users
func Info(args ...interface{}) {
	logImpl(LevelInfo, args...)
}

// Fatal will create a fatal log entry and panic
func Fatal(args ...interface{}) {
	defer panic(args)
	logImpl(LevelFatal, args...)
}

// Error will create an error log entry
func Error(args ...interface{}) {
	logImpl(LevelError, args...)
}

// Warning will create a warning log entry
func Warning(args ...interface{}) {
	logImpl(LevelWarning, args...)
}

// Debug will create a debug log entry
func Debug(args ...interface{}) {
	logImpl(LevelDebug, args...)
}

// Verbose will create a verbose log entry
func Verbose(args ...interface{}) {
	logImpl(LevelVerbose, args...)
}

func logImpl(lvl level, args ...interface{}) {
	if logLevel < lvl && lvl != LevelWTF {
		return
	}

	msg := formatLogEntry(lvl, args...)

	switch lvl {
	case LevelFatal:
		msg = color.HiRedString("%v", msg)
	case LevelError:
		msg = color.RedString("%v", msg)
	case LevelWarning:
		msg = color.YellowString("%v", msg)
	}

	go notifyListeners(msg)
	if toConsole {
		log.Println(msg)
	}
}

func formatLogEntry(lvl level, args ...interface{}) string {
	return lvl.toString() + getTag() + "\t" + argsToString(args...)
}

func argsToString(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	format := "%v"
	if strings.ContainsRune(fmt.Sprint(args[0]), '%') && len(args) > 1 {
		format = args[0].(string)
		args = args[1:]
	} else {
		format += strings.Repeat(" %v", len(args)-1)
	}
	return fmt.Sprintf(format, args...)
}

func getTag() string {
	pc, _, line, _ := runtime.Caller(4)

	caller := tagPattern.ReplaceAllString(runtime.FuncForPC(pc).Name(), `$1.$2`)

	return fmt.Sprintf("%v:%v", caller, line)
}

func notifyListeners(msg string) {
	for id, listener := range listeners {
		go notifyListener(id, msg, listener)
	}
}

func notifyListener(id string, msg string, listener chan string) {
	timeout := time.NewTimer(timeout)
	select {
	case listener <- msg:
	case <-timeout.C:
		RemoveListener(id)
		timeout.Stop()
	}
}
