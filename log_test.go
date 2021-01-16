package log

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	SetLogLevel(LogVerbose)
}

func TestConsoleLogger(t *testing.T) {
	logListeners = nil
	EnableConsoleLogging()
}

func TestDeregisterConsoleLogger(t *testing.T) {
	EnableConsoleLogging()
	WTF("TEST")
	DeregisterListener(consoleID)
}

func TestWTF(t *testing.T) {
	WTF("This is a %v", "test")
}

func TestFatal(t *testing.T) {

	if os.Getenv("EXIT") == "1" {
		Fatal("This is a %v", "test")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "EXIT=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestError(t *testing.T) {

	if os.Getenv("PANIC") == "1" {
		Error("This is a %v", "test")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestError")
	cmd.Env = append(os.Environ(), "PANIC=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestWarning(t *testing.T) {
	Warning("This is a %v", "test")
}

func TestDebug(t *testing.T) {
	Debug("This is a %v", "test")
}

func TestVerbose(t *testing.T) {
	Verbose("This is a %v", "test")
}

func TestExcludedLog(t *testing.T) {
	logLevel = LogFatal
	Verbose("This is a %v", "test")
}

func TestRegisterListener(t *testing.T) {
	testID := "test"
	testChan := RegisterListener(testID, nil)
	got := <-testChan
	if got != ConnectMessage {
		t.Errorf("Wanted: '%v'\t Received: '%v'", ConnectMessage, got)
	}
}

func TestDeregisterListener(t *testing.T) {
	testID := "test"
	logListeners[testID] = RegisterListener(testID, nil)
	testChan := logListeners[testID]
	<-testChan
	go DeregisterListener(testID)
	got := <-testChan
	if got != DisconnectMessage {
		t.Errorf("Wanted: '%v'\t Received: '%v'", DisconnectMessage, got)
	}
}

func TestEnd(t *testing.T) {
	EnableConsoleLogging()
	time.Sleep(5 * time.Millisecond)
	End()
}
func TestListenerUnavailable(t *testing.T) {
	WTF("Test")
}

func TestToLogLevelIntLow(t *testing.T) {
	want := LogFatal
	got := ToLogLevel(-1)
	if want != got {
		t.Errorf("Wanted '%v' but got '%v'", want, got)
	}
}

func TestToLogLevelIntHigh(t *testing.T) {
	want := LogFatal
	got := ToLogLevel(99999)
	if want != got {
		t.Errorf("Wanted '%v' but got '%v'", want, got)
	}
}

func TestToLogLevelValid(t *testing.T) {
	want := LogWarning
	got := ToLogLevel(int(LogWarning))
	if want != got {
		t.Errorf("Wanted '%v' but got '%v'", want, got)
	}
}
