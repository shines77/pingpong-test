package main

import (
	"strconv"
	"strings"
)

func formatString(str *string) string {
	if str == nil {
		return "<nil>"
	} else if strings.TrimSpace(*str) == "" {
		return "<empty>"
	} else {
		return *str
	}
}

func formatBool(b *bool) string {
	if b == nil {
		return "<nil>"
	} else {
		bval := *b
		if bval {
			return "true"
		} else {
			return "false"
		}
	}
}

func parseBool(str *string, defaultVal bool) bool {
	if str == nil {
		return defaultVal
	} else if strings.TrimSpace(*str) == "" {
		return defaultVal
	} else {
		boolstr := strings.TrimSpace(*str)
		boolstr = strings.ToLower(boolstr)
		if boolstr == "true" || boolstr == "1" {
			return true
		} else if boolstr == "false" || boolstr == "0" {
			return false
		} else {
			val, err := strconv.ParseInt(boolstr, 10, 32)
			if err != nil && val != 0 {
				return true
			} else {
				return defaultVal
			}
		}
	}
}
