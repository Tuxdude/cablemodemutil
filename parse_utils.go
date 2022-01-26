package cablemodemutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// DAY MON DATE HH:MM:SS YYYY
	systemTimestampFormat = "Mon Jan 2 15:04:05 2006"
)

// Parses the specified string as an uint64 after stripping the suffix if required.
func parseUint64(str string, hasSuffix bool, suffix string, desc string) (uint64, error) {
	if hasSuffix {
		if !strings.HasSuffix(str, suffix) {
			return 0, fmt.Errorf("expected %s with %q suffix, but not available in %q", desc, suffix, str)
		}
		str = strings.TrimSuffix(str, suffix)
	}

	res, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to convert %q to uint64: %w", str, err)
	}
	return res, nil
}

// Parses the specified string as an int64 after stripping the suffix if required.
func parseInt64(str string, hasSuffix bool, suffix string, desc string) (int64, error) {
	if hasSuffix {
		if !strings.HasSuffix(str, suffix) {
			return 0, fmt.Errorf("expected %s with %q suffix, but not available in %q", desc, suffix, str)
		}
		str = strings.TrimSuffix(str, suffix)
	}

	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to convert %q to int64: %w", str, err)
	}
	return res, nil
}

// Parses the specified string as a channel frequency after stripping the ' Hz' suffix if required.
func parseFreqStr(str string, hasHzSuffix bool, desc string) (uint32, error) {
	f, err := parseUint64(str, hasHzSuffix, " Hz", desc)
	if err != nil {
		return 0, err
	}
	return uint32(f), nil
}

// Parses the specified string as a signal power integer value after stripping the ' dBmV' suffix if required.
func parseSignalPowerIntStr(str string, hasDBMVSuffix bool, desc string) (int32, error) {
	pow, err := parseInt64(str, hasDBMVSuffix, " dBmV", desc)
	if err != nil {
		return 0, err
	}
	return int32(pow), nil
}

// Parses the specified string as a signal SNR integer value after stripping the ' dB' suffix if required.
func parseSignalSNRStr(str string, hasDBSuffix bool, desc string) (int32, error) {
	snr, err := parseInt64(str, hasDBSuffix, " dB", desc)
	if err != nil {
		return 0, err
	}
	return int32(snr), nil
}

// Parses the specified string as a channel ID value.
func parseChannelIDStr(str string, desc string) (uint8, error) {
	id, err := parseUint64(str, false, "", desc)
	if err != nil {
		return 0, err
	}
	return uint8(id), nil
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
			return 0, fmt.Errorf("Unable to parse days in the specified %s timestamp %q, reason: %w", desc, str, err)
		}
		days = time.Duration(d)
		nextIndex = 1
	}
	hoursMinsSecs := strings.Split(components[nextIndex], ":")
	if len(hoursMinsSecs) != 3 {
		return 0, fmt.Errorf("Unable to split hours:mins:secs in the specified %s timestamp %q", desc, str)
	}
	hours, err := parseTimeElementWithSuffix(hoursMinsSecs[0], "h")
	if err != nil {
		return 0, fmt.Errorf("Unable to parse hours in the specified %s timestamp %q, reason: %w", desc, str, err)
	}
	mins, err := parseTimeElementWithSuffix(hoursMinsSecs[1], "m")
	if err != nil {
		return 0, fmt.Errorf("Unable to parse hours in the specified %s timestamp %q, reason: %w", desc, str, err)
	}
	secs, err := parseTimeElementWithSuffix(hoursMinsSecs[2], "s")
	if err != nil {
		return 0, fmt.Errorf("Unable to parse hours in the specified %s timestamp %q, reason: %w", desc, str, err)
	}
	return (time.Duration(days) * 24 * time.Hour) + (time.Duration(hours) * time.Hour) + (time.Duration(mins) * time.Minute) + (time.Duration(secs) * time.Second), nil
}

// Parses the time components with the specified suffix from the specified string.
func parseTimeElementWithSuffix(str string, suffix string) (uint32, error) {
	components := strings.Split(str, suffix)
	if len(components) != 2 {
		return 0, fmt.Errorf("Unable to parse string %q (split: %d), expected to have suffix %q but did not", str, len(components), suffix)
	}

	num, err := strconv.ParseUint(components[0], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Unable to parse component %q as a uint", num)
	}
	return uint32(num), nil
}

// Parses the value of the specified key as a string in the specified status information.
func parseString(data actionResponseBody, key string, desc string) (string, error) {
	s, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("unable to find key %q for parsing %q.data=%v", key, desc, data)
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

// Parses the value of the specified key as a signal power integer value in the specified status information.
func parseSignalPowerInt(data actionResponseBody, key string, hasDBMVSuffix bool, desc string) (int32, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return 0, err
	}
	return parseSignalPowerIntStr(s, hasDBMVSuffix, desc)
}

// Parses the value of the specified key as a signal SNR integer value in the specified status information.
func parseSignalSNR(data actionResponseBody, key string, hasDBSuffix bool, desc string) (int32, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return 0, err
	}
	return parseSignalSNRStr(s, hasDBSuffix, desc)
}

// Parses the value of the specified key as a channel ID integer value in the specified status information.
func parseChannelID(data actionResponseBody, key string, desc string) (uint8, error) {
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
