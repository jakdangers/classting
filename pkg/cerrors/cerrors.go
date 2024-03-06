package cerrors

import (
	"bytes"
	"errors"
	"fmt"
)

type Op string
type Kind uint8

var Separator = ":\n\t"

const (
	Other      Kind = iota // 분류되지 않은 오류의 경우 (이 값은 메시지에 포함하지 않음).
	Invalid                // 유효하지 않은 행위를 한 경우.
	Auth                   // 비인증.
	Permission             // 권한이 옳바르지 않은 경우.
	Exist                  // 이미 존재하는 경우.
	NotExist               // 존재하지 않는 경우.
	Internal               // 로직 오류의 경우.
)

type Error struct {
	Op             Op     // 도메인/액션
	Kind           Kind   // 에러 종류
	Err            error  // 에러
	ServiceMessage string // 클라이언트 전용 메시지
}

// Error 관련
func (e *Error) Error() string {
	b := new(bytes.Buffer)

	if e.Op != "" {
		pad(b, ": ")
		b.WriteString(string(e.Op))
	}

	if e.Kind != 0 {
		pad(b, ": ")
		b.WriteString(e.Kind.String())
	}

	if e.Err != nil {
		var prevErr *Error
		if errors.As(e.Err, &prevErr) {
			pad(b, Separator)
			b.WriteString(e.Err.Error())
		} else {
			pad(b, ": ")
			b.WriteString(e.Err.Error())
		}
	}

	if b.Len() == 0 {
		return "no error"
	}

	return b.String()
}

func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to cerrors.E with no arguments")
	}
	e := &Error{}

	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case string:
			e.ServiceMessage = arg
		case Kind:
			e.Kind = arg
		case error:
			e.Err = arg
		case *Error:
			copy := *arg
			e.Err = &copy
		default:
			return errors.New(fmt.Sprintf("unknown type %T, value %v in error call", arg, arg))
		}
	}

	return e
}

func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case Invalid:
		return "invalid operation"
	case Auth:
		return "unauthorized"
	case Permission:
		return "permission denied"
	case Exist:
		return "item already exists"
	case NotExist:
		return "item does not exist"
	case Internal:
		return "internal error"
	}
	return "unknown error kind"
}

func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}

	b.WriteString(str)
}
