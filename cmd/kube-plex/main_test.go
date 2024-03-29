package main

import (
	"os"
	"testing"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

// disable logging in all tests
func TestMain(m *testing.M) {
	klog.SetLogger(logr.Discard())
	os.Exit(m.Run())
}

func Test_needBypass(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{"bypass aec3_eae", []string{"...", "-codec:1", "eac3_eae", "-eaeprefix:1", "..."}, true},
		{"bypass ac3_eae", []string{"...", "-codec:1", "ac3_eae", "-eaeprefix:1", "..."}, true},
		{"bypass truehd_eae", []string{"...", "-codec:1", "truehd_eae", "-eaeprefix:1", "..."}, true},
		{"bypass mlp_eae", []string{"...", "-codec:1", "mlp_eae", "-eaeprefix:1", "..."}, true},
		{"don't bypass with ac3", []string{"...", "-codec:1", "ac3", "-prefix:1", "..."}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := needBypass(tt.args); got != tt.want {
				t.Errorf("needBypass() = %v, want %v", got, tt.want)
			}
		})
	}
}
