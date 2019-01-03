package main

import (
	"strconv"
	"strings"
)

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

func ParseInfo(info string) map[string]interface{} {
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

func parseClient(client string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, item := range strings.Split(client, " ") {
		kv := strings.SplitN(item, "=", 2)
		result[kv[0]] = getClientValue(kv[1])
	}
	return result
}

func getClientValue(value string) interface{} {
	if result, err := strconv.ParseInt(value, 10, 64); err != nil {
		return value
	} else {
		return result
	}
}

func ParseClientList(clientList string) []map[string]interface{} {
	var list []map[string]interface{}
	for _, client := range strings.Split(clientList, "\r\n") {
		list = append(list, parseClient(client))
	}
	return list
}
