package cablemodemutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Log timestamps are in the format "DD/MM/YYYY HH:MM:SS".
	eventLogTimestampFormat = "2/1/2006 15:04:05"
	// System timestamps are in the format "DAY MON DATE HH:MM:SS YYYY".
	systemTimestampFormat = "Mon Jan 2 15:04:05 2006"
)

// Parses the specified string as an uint32 after stripping the suffix if required.
func parseUint32(str string, hasSuffix bool, suffix string, desc string) (uint32, error) {
	if hasSuffix {
		if !strings.HasSuffix(str, suffix) {
			return 0, fmt.Errorf(
				"expected %s with %q suffix, but not available in %q",
				desc,
				suffix,
				str,
			)
		}
		str = strings.TrimSuffix(str, suffix)
	}

	res, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parsing %q, unable to convert %q to uint32: %w", desc, str, err)
	}
	return uint32(res), nil
}

// Parses the specified string as a float64 after stripping the suffix if required.
func parseFloat32(str string, hasSuffix bool, suffix string, desc string) (float32, error) {
	if hasSuffix {
		if !strings.HasSuffix(str, suffix) {
			return 0, fmt.Errorf("expected %s with %q suffix, but not available in %q", desc, suffix, str)
		}
		str = strings.TrimSuffix(str, suffix)
	}

	res, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, fmt.Errorf("parsing %q, unable to convert %q to float32: %w", desc, str, err)
	}
	return float32(res), nil
}

// Parses the specified string as a channel frequency after stripping the ' Hz' suffix if required.
func parseFreqStr(str string, hasHzSuffix bool, desc string) (uint32, error) {
	return parseUint32(str, hasHzSuffix, " Hz", desc)
}

// Parses the specified string as a signal power floating point value after stripping the ' dBmV' suffix if required.
func parseSignalPowerStr(str string, hasDBMVSuffix bool, desc string) (float32, error) {
	return parseFloat32(str, hasDBMVSuffix, " dBmV", desc)
}

// Parses the specified string as a signal SNR floating point value after stripping the ' dB' suffix if required.
func parseSignalSNRStr(str string, hasDBSuffix bool, desc string) (float32, error) {
	return parseFloat32(str, hasDBSuffix, " dB", desc)
}

// Parses the specified string as a signal errors integer value.
func parseSignalErrorsStr(str string, desc string) (uint32, error) {
	return parseUint32(str, false, "", desc)
}

// Parses the specified string as a channel ID value.
func parseChannelIDStr(str string, desc string) (uint32, error) {
	return parseUint32(str, false, "", desc)
}

// Parses the log timestamp from the specified date and time string values.
func parseLogTimestamp(dateStr string, timeStr string) (time.Time, error) {
	timestamp := fmt.Sprintf("%s %s", dateStr, timeStr)
	// TODO: Do this one time instead of doing this for every call of this function.
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, fmt.Errorf("error loading \"Local\" location, reason: %w", err)
	}
	t, err := time.ParseInLocation(eventLogTimestampFormat, timestamp, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing timestamp %q, reason: %w", timestamp, err)
	}
	return t, nil
}

// Parses the system timestamp from the specified timestamp string.
func parseSystemTimestampStr(timestamp string, desc string) (time.Time, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, fmt.Errorf("error loading \"Local\" location, reason: %w", err)
	}
	t, err := time.ParseInLocation(systemTimestampFormat, timestamp, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing %s timestamp %q, reason: %w", desc, timestamp, err)
	}
	return t, nil
}

// Parses the time duration information from the specified duration string.
func parseDurationStr(str string, desc string) (time.Duration, error) {
	// Eg. "3 days 14h:15m:33s"
	days := time.Duration(0)
	components := strings.Split(str, " days ")
	nextIndex := 0
	if len(components) > 1 {
		d, err := strconv.ParseUint(components[0], 10, 32)
		if err != nil {
			return 0, fmt.Errorf("unable to parse days in the specified %s timestamp %q, reason: %w", desc, str, err)
		}
		days = time.Duration(d)
		nextIndex = 1
	}
	hoursMinsSecs := strings.Split(components[nextIndex], ":")
	if len(hoursMinsSecs) != 3 {
		return 0, fmt.Errorf("unable to split hours:mins:secs in the specified %s timestamp %q", desc, str)
	}
	hours, err := parseTimeElementWithSuffix(hoursMinsSecs[0], "h")
	if err != nil {
		return 0, fmt.Errorf("unable to parse hours in the specified %s timestamp %q, reason: %w", desc, str, err)
	}
	mins, err := parseTimeElementWithSuffix(hoursMinsSecs[1], "m")
	if err != nil {
		return 0, fmt.Errorf("unable to parse hours in the specified %s timestamp %q, reason: %w", desc, str, err)
	}
	secs, err := parseTimeElementWithSuffix(hoursMinsSecs[2], "s")
	if err != nil {
		return 0, fmt.Errorf("unable to parse hours in the specified %s timestamp %q, reason: %w", desc, str, err)
	}
	res := days * 24 * time.Hour
	res += time.Duration(hours) * time.Hour
	res += time.Duration(mins) * time.Minute
	res += time.Duration(secs) * time.Second
	return res, nil
}

// Parses the time components with the specified suffix from the specified string.
func parseTimeElementWithSuffix(str string, suffix string) (uint32, error) {
	components := strings.Split(str, suffix)
	if len(components) != 2 {
		return 0, fmt.Errorf(
			"unable to parse string %q (split: %d), expected to have suffix %q but did not",
			str,
			len(components),
			suffix,
		)
	}

	num, err := strconv.ParseUint(components[0], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("unable to parse component %q as a uint", num)
	}
	return uint32(num), nil
}

// Parses the specified log entry string.
func parseLogEntry(str string) string {
	// Just replace two spaces with one (seen commonly with login log entries).
	return strings.ReplaceAll(str, "  ", " ")
}

// Parses the value of the specified key as a string in the specified status information.
func parseString(data actionResponseBody, key string, desc string) (string, error) {
	s, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("unable to find key %q while parsing %q.\ndata=%v", key, desc, data)
	}
	return s, nil
}

// Parses the value of the specified key as a bool in the specified status information.
func parseBool(data actionResponseBody, key string, trueVal string, desc string) (bool, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return false, err
	}
	return s == trueVal, nil
}

// Parses the value of the specified key as a channel frequency in the specified status information.
func parseFreq(data actionResponseBody, key string, hasHzSuffix bool, desc string) (uint32, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return 0, err
	}
	return parseFreqStr(s, hasHzSuffix, desc)
}

// Parses the value of the specified key as a signal power floating point value in the specified status information.
func parseSignalPower(data actionResponseBody, key string, hasDBMVSuffix bool, desc string) (float32, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return 0, err
	}
	return parseSignalPowerStr(s, hasDBMVSuffix, desc)
}

// Parses the value of the specified key as a signal SNR floating point value in the specified status information.
func parseSignalSNR(data actionResponseBody, key string, hasDBSuffix bool, desc string) (float32, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return 0, err
	}
	return parseSignalSNRStr(s, hasDBSuffix, desc)
}

// Parses the value of the specified key as a channel ID integer value in the specified status information.
func parseChannelID(data actionResponseBody, key string, desc string) (uint32, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return 0, err
	}
	return parseChannelIDStr(s, desc)
}

// Parses the value of the specified key as a system timestamp in the specified status information.
func parseSystemTimestamp(data actionResponseBody, key string, desc string) (time.Time, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return time.Time{}, err
	}
	return parseSystemTimestampStr(s, desc)
}

// Parses the value of the specified key as a time duration in the specified status information.
func parseDuration(data actionResponseBody, key string, desc string) (time.Duration, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return 0, err
	}
	return parseDurationStr(s, desc)
}
