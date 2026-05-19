package system

import (
	"fmt"
	"os/exec"
	"runtime"
)

type ToolInfo struct {
	Available bool   `json:"available"`
	Path      string `json:"path,omitempty"`
}

type Environment struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type DoctorResult struct {
	OK          bool        `json:"ok"`
	Command     string      `json:"command"`
	Environment Environment `json:"environment"`
	Tools       ToolMap     `json:"tools"`
	Version     string      `json:"version"`
}

type ToolMap struct {
	IPATool ToolInfo `json:"ipatool"`
	Plutil  ToolInfo `json:"plutil"`
	Strings ToolInfo `json:"strings"`
}

func CheckTool(name string) ToolInfo {
	path, err := exec.LookPath(name)
	if err != nil {
		return ToolInfo{Available: false}
	}
	return ToolInfo{Available: true, Path: path}
}

func RunDoctor(version string) DoctorResult {
	env := Environment{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	tools := ToolMap{
		IPATool: CheckTool("ipatool"),
		Plutil:  CheckTool("plutil"),
		Strings: CheckTool("strings"),
	}

	allOK := tools.IPATool.Available && tools.Plutil.Available && tools.Strings.Available

	return DoctorResult{
		OK:          allOK,
		Command:     "doctor",
		Environment: env,
		Tools:       tools,
		Version:     version,
	}
}

func MissingToolError(tool string) string {
	switch tool {
	case "ipatool":
		return fmt.Sprintf("%s is not installed. Install it via: brew install majd/repo/ipatool", tool)
	case "plutil":
		return fmt.Sprintf("%s is not found. It should be available on macOS by default at /usr/bin/plutil", tool)
	case "strings":
		return fmt.Sprintf("%s is not found. Install Xcode Command Line Tools: xcode-select --install", tool)
	default:
		return fmt.Sprintf("%s is not installed. Please install it before proceeding.", tool)
	}
}
