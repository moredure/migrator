package main

import (
	"strconv"
	"strings"
)

type RedisMetrics map[string]interface{}

func getValue(value string) interface{} {
	if !strings.Contains(value, ",") || !strings.Contains(value, "=") {
		if strings.Contains(value, ".") {
			result, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return value
			}
			return result
		}
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return value
		}
		return result
	} else {
		result := make(map[string]interface{})
		for _, item := range strings.Split(value, ",") {
			parts := strings.Split(item, "=")
			result[parts[0]] = getValue(parts[1])
		}
		return result
	}
}

func ParseMetrics(info string) RedisMetrics {
	m := make(map[string]interface{})
	for _, line := range strings.Split(info, "\r\n") {
		if len(line) != 0 && !strings.HasPrefix(line, "#") {
			result := strings.Split(line, ":")
			if len(result) == 2 {
				m[result[0]] = getValue(result[1])
			}
		}
	}
	return m
}

func ParseClient(client string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, item := range strings.Split(client, " ") {
		kv := strings.Split(item, "=")
		result[kv[0]] = kv[1]
	}
	return result
}