package log

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func init() {
	timeout = 10 * time.Millisecond
}

func TestInterface(t *testing.T) {
	defer catchPanic()
	SetLogLevel(LevelWarning)
	ToConsole(true)
	WTF("test")
	Info("test")
	Error("test")
	Warning("test")
	Debug("test")
	Verbose("test")
	Fatal("test")
}

func TestGetLevelFromStringSuccess(t *testing.T) {
	want := LevelVerbose
	got := GetLevelFromString("VeRbOsE")

	if want != got {
		t.FailNow()
	}
}

func TestGetLevelFromStringInvalid(t *testing.T) {
	want := LevelWTF
	got := GetLevelFromString("Invalid")

	if want != got {
		t.FailNow()
	}
}

func TestSetTimeout(t *testing.T) {
	SetTimeout(10 * time.Millisecond)
}

func TestLogLevelToString(t *testing.T) {
	level := LevelWTF
	want := "WTF      "
	got := level.toString()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}

func TestAddListenerSuccess(t *testing.T) {
	got := AddListener("TestAddListenerSuccess", make(chan string))

	if diff := cmp.Diff(nil, got); diff != "" {
		t.Error(diff)
	}
}

func TestAddListenerNoChannel(t *testing.T) {
	got := AddListener("TestAddListenerNoChannel", nil)

	want := errors.New("attempted to register a log receiver without a channel")

	if diff := cmp.Diff(want.Error(), got.Error()); diff != "" {
		t.Error(diff)
	}
}

func TestAddListenerDuplicateId(t *testing.T) {
	AddListener("TestAddListenerDuplicateId", make(chan string))
	got := AddListener("TestAddListenerDuplicateId", make(chan string))

	want := errors.New("attempted to register a log receiver without a unique identifier")

	if diff := cmp.Diff(want.Error(), got.Error()); diff != "" {
		t.Error(diff)
	}
}

func TestRemoveListener(t *testing.T) {
	AddListener("TestRemoveListener", make(chan string))
	RemoveListener("TestRemoveListener")
}

func TestLogImplToConsole(t *testing.T) {
	logImpl(LevelVerbose, "Test")
}

func TestArgsToStringNoArgs(t *testing.T) {
	if argsToString() != "" {
		t.FailNow()
	}
}

func TestArgsToStringFormatting(t *testing.T) {
	got := argsToString("%+v", 10)
	if Diff := cmp.Diff(got, "10"); Diff != "" {
		t.Error(Diff)
	}
}

func TestArgsToStringMultipleArgs(t *testing.T) {
	if Diff := cmp.Diff(argsToString("test", 1, 2, 3), "test 1 2 3"); Diff != "" {
		t.Error(Diff)
	}
}

func TestGetTag(t *testing.T) {
	getTag()
}

func TestNotifyListeners(t *testing.T) {
	logChan1 := make(chan string)
	logChan2 := make(chan string)
	listeners["TestNotifyListeners1"] = logChan1
	listeners["TestNotifyListeners2"] = logChan2

	notifyListeners("TestNotifyListeners")

	got := <-logChan1
	got += <-logChan2
	want := "TestNotifyListenersTestNotifyListeners"

	if Diff := cmp.Diff(want, got); Diff != "" {
		t.Error(Diff)
	}
}

func TestNotifyListenerSuccess(t *testing.T) {
	logChan := make(chan string)
	listeners["TestNotifyListenerSuccess"] = logChan

	go notifyListener("TestNotifyListenerSuccess", "test", logChan)

	got := <-logChan
	want := "test"

	if Diff := cmp.Diff(want, got); Diff != "" {
		t.Error(Diff)
	}
}

func TestNotifyListenerTimeout(t *testing.T) {
	logChan := make(chan string)
	listeners["TestNotifyListenerTimeout"] = logChan

	notifyListener("TestNotifyListenerTimeout", "test", logChan)

	if listeners["TestNotifyListenerTimeout"] != nil {
		t.FailNow()
	}
}

func catchPanic() {
	if r := recover(); r != nil {
		fmt.Println("Recovered", r)
	}
}
