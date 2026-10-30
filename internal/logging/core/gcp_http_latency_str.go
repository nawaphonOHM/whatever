package core

import (
	"math"
	"strconv"
	"time"
)

const (
	nanosPerMilli   = 1e6
	nanosDigitCount = 9
	nanosLastIndex  = nanosDigitCount - 1
)

func formatNanosDigits(nanos int64) [nanosDigitCount]byte {
	rem := nanos
	var buf [nanosDigitCount]byte
	for i := nanosLastIndex; i >= 0; i-- {
		buf[i] = byte('0' + rem%decimalBase)
		rem /= decimalBase
	}
	return buf
}

func formatNanosFraction(nanos int64) string {
	buf := formatNanosDigits(nanos)
	end := nanosLastIndex
	for end >= 0 && buf[end] == '0' {
		end--
	}
	return string(buf[:end+1])
}

func formatDuration(d time.Duration) string {
	sec := int64(d / time.Second)
	nanos := int64(d % time.Second)
	if nanos == 0 {
		return strconv.FormatInt(sec, decimalBase) + durationSuffix
	}
	return strconv.FormatInt(sec, decimalBase) + "." + formatNanosFraction(nanos) + durationSuffix
}

func isNonNegative(f, maxVal float64) bool {
	return f >= 0 && f <= maxVal
}

func isFiniteFloat(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}

func isFiniteNonNegativeFloat(f, maxVal float64) bool {
	return isNonNegative(f, maxVal) && isFiniteFloat(f)
}

func parseDurationString(s string) (string, bool) {
	d, err := time.ParseDuration(s)
	if err == nil && d >= 0 {
		return formatDuration(d), true
	}
	return "", false
}
