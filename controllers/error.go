package controllers

import (
	"fmt"

	"github.com/go-sfox-lib/sfox/log"
)

type ErrorCode int32

var log = logger.FindOrCreateLoggerInstance(logger.NewLoggerConfig("debug", "", "", ""))

const (
	ErrorNo ErrorCode = 0
	// 100~200
	ErrorCodeAccountNameError   ErrorCode = 101
	ErrorCodeAccountPwdError    ErrorCode = 102
	ErrorCodeAccountNotExist    ErrorCode = 103
	ErrorCodeAccountInsertError ErrorCode = 104

	// 200~300
	ErrorCodeTxInsertError ErrorCode = 201
	// 300~400
	ErrorCodeSystemDecodeError ErrorCode = 301
)

func NewErrorMsg(msg string, code ErrorCode) string {
	if code == ErrorNo {
		return fmt.Sprintf("{\"state\": \"success\"}")
	}
	return fmt.Sprintf("{\"err_msg\": \"%s\", \"err_code\": %d}", msg, code)
}
