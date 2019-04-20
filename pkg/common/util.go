package common

import (
	"fmt"
	"strings"
	"sync"
)

func GetMapValueString(ptr *sync.Map, key string) string {
	valueRaw, _ := ptr.Load(key)
	value, _ := valueRaw.(string)
	return value
}

func GetMapValueStringArr(ptr *sync.Map, key string) []string {
	valueRaw, _ := ptr.Load(key)
	value, _ := valueRaw.([]string)
	return value
}

func BuildTable(list [][]string) []string {
	if len(list) == 0 {
		return nil
	}
	_len := len(list[0])
	widths := make([]int, _len)
	for _, row := range list {
		if len(row) != _len {
			return nil
		}
		for idx, field := range row {
			if len(field) > widths[idx] {
				widths[idx] = len(field)
			}
		}
	}
	format := strings.TrimSpace(strings.Repeat("%s ", _len))
	var res []string
	for _, row := range list {
		var fields []interface{}
		for idx, field := range row {
			fields = append(fields, field+strings.Repeat(" ", widths[idx]-len(field)))
		}
		res = append(res, fmt.Sprintf(format, fields...))
	}
	return res
}
