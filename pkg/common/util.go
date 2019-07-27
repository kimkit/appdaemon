package common

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
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

func GetMapValueInt(ptr *sync.Map, key string) int {
	valueRaw, _ := ptr.Load(key)
	value, _ := valueRaw.(int)
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
			if idx == len(row)-1 {
				fields = append(fields, field)
			} else {
				fields = append(fields, field+strings.Repeat(" ", widths[idx]-len(field)))
			}
		}
		res = append(res, fmt.Sprintf(format, fields...))
	}
	return res
}

func Time2str(t int) string {
	if t <= 0 {
		return "0000-00-00 00:00:00"
	}
	return time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
}

type MapSlice struct {
	Array []map[string]string
	Key   string
}

func (p *MapSlice) Len() int {
	return len(p.Array)
}

func (p *MapSlice) Less(i, j int) bool {
	if strings.Compare(p.Array[i][p.Key], p.Array[j][p.Key]) < 0 {
		return true
	} else {
		return false
	}
}

func (p *MapSlice) Swap(i, j int) {
	p.Array[i], p.Array[j] = p.Array[j], p.Array[i]
}

func SortMaps(maps []map[string]string, key string) {
	p := &MapSlice{maps, key}
	sort.Sort(p)
}

func Addslashes(str string) string {
	str = strings.Replace(str, "\\", "\\\\", -1)
	str = strings.Replace(str, "\"", "\\\"", -1)
	str = strings.Replace(str, "'", "\\'", -1)
	return str
}

var (
	shellSCRegexp = regexp.MustCompile("[ \t\r\n`$\\\";&|<>]")
)

func Args2str(args []string) string {
	var _args []string
	for _, arg := range args {
		if arg == "" {
			_args = append(_args, "\"\"")
			continue
		}
		num := 0
		str := shellSCRegexp.ReplaceAllStringFunc(arg, func(m string) string {
			num++
			switch m {
			case " ":
				return " "
			case "\t":
				return "\\t"
			case "\r":
				return "\\r"
			case "\n":
				return "\\n"
			case "`":
				return "\\`"
			case "$":
				return "\\$"
			case "\\":
				return "\\\\"
			case "\"":
				return "\\\""
			default:
				return m
			}
		})
		if num == 0 {
			_args = append(_args, str)
		} else {
			_args = append(_args, fmt.Sprintf("\"%s\"", str))
		}
	}
	return strings.Join(_args, " ")
}
