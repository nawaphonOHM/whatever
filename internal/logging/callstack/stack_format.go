package callstack

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func writeBuilderString(sb *strings.Builder, s string) {
	if _, err := sb.WriteString(s); err != nil {
		return
	}
}

func formatSingleFrame(sb *strings.Builder, i int, f *Frame) {
	if f == nil {
		return
	}
	pkg := f.Package
	if pkg != "" {
		pkg += "."
	}
	line := fmt.Sprintf("  [%d] %s%s\n      at %s:%d\n      time: %s\n",
		i, pkg, f.Function, f.File, f.Line, f.Timestamp.Format(time.RFC3339Nano))
	writeBuilderString(sb, line)
}

// String returns a formatted multiline representation of the entire call stack.
func (cs *CallStack) String() string {
	if cs == nil {
		return "<nil stack>"
	}
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	if len(cs.frames) == 0 {
		return "<empty stack>"
	}
	return formatFrames(cs.frames)
}

func formatFrames(frames []*Frame) string {
	var sb strings.Builder
	writeBuilderString(&sb, "CallStack:\n")
	for i, f := range frames {
		formatSingleFrame(&sb, i, f)
	}
	return sb.String()
}

func applyFuncInfo(fn *runtime.Func, frame *Frame) {
	pkg, fnName := splitPackageAndFunc(fn.Name())
	frame.Package = pkg
	if frame.Function == "" {
		frame.Function = fnName
	}
}

func resolveCallerInfo(pc uintptr, frame *Frame) {
	if frame == nil {
		return
	}
	if fn := runtime.FuncForPC(pc); fn != nil {
		applyFuncInfo(fn, frame)
	}
}

func captureCallerFrame(name string, skip int) *Frame {
	frame := &Frame{
		Function:  name,
		Timestamp: time.Now(),
	}
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		frame.File = file
		frame.Line = line
		resolveCallerInfo(pc, frame)
	}
	return frame
}

func splitPackageAndFunc(fullName string) (prefix, rest string) {
	lastSlash := strings.LastIndex(fullName, "/")
	prefix = ""
	rest = fullName
	if lastSlash >= 0 {
		prefix = fullName[:lastSlash+1]
		rest = fullName[lastSlash+1:]
	}
	firstDot := strings.Index(rest, ".")
	if firstDot >= 0 {
		return prefix + rest[:firstDot], rest[firstDot+1:]
	}
	return prefix, rest
}
