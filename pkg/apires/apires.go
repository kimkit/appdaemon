package apires

import (
	"fmt"
	"regexp"
	"strconv"
)

type Reply struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Error struct {
	Reply *Reply
}

func (err *Error) Error() string {
	return fmt.Sprintf("apires.Error: (%d) %s", err.Reply.Code, err.Reply.Message)
}

func (err *Error) Clone() error {
	return NewError(err.Reply.Code, err.Reply.Message, err.Reply.Data)
}

func NewReply(code int, message string, data interface{}) *Reply {
	return &Reply{code, message, data}
}

func NewError(code int, message string, data interface{}) error {
	return &Error{&Reply{code, message, data}}
}

var (
	errorRegexp = regexp.MustCompile(`^apires\.Error: \((\-?[0-9]+)\) (.*)$`)
)

func ParseError(str string) error {
	arr := errorRegexp.FindStringSubmatch(str)
	if len(arr) == 3 {
		code, _ := strconv.Atoi(arr[1])
		return NewError(code, arr[2], nil)
	} else {
		return fmt.Errorf("%s", str)
	}
}
