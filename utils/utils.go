package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	SysWindows        = "windows"
	SysMac            = "darwin"
	TimeLayout string = "2006-01-02 15:04:05"
)

// GetGOPATHs returns all paths in GOPATH variable.
func GetGOPATHs() []string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" && strings.Compare(runtime.Version(), "go1.8") >= 0 {
		gopath = defaultGOPATH()
	}
	return filepath.SplitList(gopath)
}

func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == SysWindows {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		return filepath.Join(home, "go")
	}
	return ""
}

func IsWindows() bool {
	return runtime.GOOS == SysWindows
}

func IsMac() bool {
	return runtime.GOOS == SysMac
}
