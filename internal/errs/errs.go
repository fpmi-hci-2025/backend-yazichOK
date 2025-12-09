package errs

import (
	"fmt"
	"runtime"
	"strconv"
)

const (
	unknown         = "UNKNOWN"
	skipDirsInStack = 1

	errorTmpl = "%s; %s\n file: %s\n line: %s"
	wrapTmpl  = "%s : %w"
)

type Error struct {
	err  error
	msg  string
	file string
	line string
}

func New(err error, msg string) Error {
	e := Error{
		err: err,
		msg: msg,
	}

	_, file, line, ok := runtime.Caller(skipDirsInStack)
	if !ok {
		e.file = unknown
		e.line = unknown

		return e
	}

	e.file = file
	e.line = strconv.Itoa(line)

	return e
}

func (e Error) Error() string {
	return fmt.Sprintf(
		errorTmpl,
		e.err.Error(),
		e.msg,
		e.file,
		e.line,
	)
}

func (e Error) Unwrap() error {
	return e.err
}

func Wrap(msg string, err error) error {
	return fmt.Errorf(wrapTmpl, msg, err)
}
