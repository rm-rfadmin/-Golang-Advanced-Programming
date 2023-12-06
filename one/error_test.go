package one

import (
	"fmt"
	"github.com/chai2010/errors"
	"io"
	"os"
	"testing"
)

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}

	written, err = io.Copy(dst, src)
	dst.Close()
	src.Close()
	return
}

func CopyFile1(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}

	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}

	defer dst.Close()

	return io.Copy(dst, src)
}

//func ParseJSON(input string) (s *Syntax, err error){
//	defer func() {
//		if p:= recover(); p != nil {
//			err = fmt.Errorf("JSON: internal error: %v",p)
//		}
//	}()
//}

type Error interface {
	Caller() []CallerInfo
	Wrapewd() []error
	Code() int
	error

	private()
}

type CallerInfo struct {
	FuncName string
	FileName string
	FileLine string
}

//func New(msg string) error
//func NewWithCode(code int, msg string) error
//
//func Wrap(err error, msg string) error

func TestError1(t *testing.T) {
	//err := syscall.Chmod(":invalid path:")
	err := errors.NewWithCode(404, "http error code")
	fmt.Println(err)
	fmt.Println(err.(errors.Error).Code())
}
