package callstack

import (
	"fmt"
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
