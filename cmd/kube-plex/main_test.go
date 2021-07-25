package main

import (
	"os"
	"testing"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

// EmptyLogger implements logr.Logging
type EmptyLogger struct{}

// Enabled always returns false
func (e *EmptyLogger) Enabled() bool {
	return false
}

// Info does nothing
func (e *EmptyLogger) Info(msg string, keysAndValues ...interface{}) {}

// Error does nothing
func (e *EmptyLogger) Error(err error, msg string, keysAndValues ...interface{}) {}

// V returns itself
func (e *EmptyLogger) V(level int) logr.Logger {
	return e
}

// WithValues returns itself
func (e *EmptyLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	return e
}

// WithName returns itself
func (e *EmptyLogger) WithName(name string) logr.Logger {
	return e
}

// disable logging in all tests
func TestMain(m *testing.M) {
	klog.SetLogger(&EmptyLogger{})
	os.Exit(m.Run())
}

func Test_needBypass(t *testing.T) {
	type args struct {
		args []string
		m    PmsMetadata
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"bypass aec3_eae", args{args: []string{"...", "-codec:1", "eac3_eae", "-eaeprefix:1", "..."}}, true},
		{"bypass ac3_eae", args{args: []string{"...", "-codec:1", "ac3_eae", "-eaeprefix:1", "..."}}, true},
		{"bypass truehd_eae", args{args: []string{"...", "-codec:1", "truehd_eae", "-eaeprefix:1", "..."}}, true},
		{"bypass mlp_eae", args{args: []string{"...", "-codec:1", "mlp_eae", "-eaeprefix:1", "..."}}, true},
		{"don't bypass with ac3", args{args: []string{"...", "-codec:1", "ac3", "-prefix:1", "..."}}, false},
		{"don't bypass if explicit eae root is set", args{args: []string{"...", "-codec:1", "mlp_eae", "-eaeprefix:1", "..."}, m: PmsMetadata{EaeRootDir: "/eaeroot"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := needBypass(tt.args.args, tt.args.m); got != tt.want {
				t.Errorf("needBypass() = %v, want %v", got, tt.want)
			}
		})
	}
}
