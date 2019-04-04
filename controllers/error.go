package controllers

import (
	"fmt"
)

type ErrorCode int32

const (
	// 100~200
	ErrorCodeAccountNameError   ErrorCode = 101
	ErrorCodeAccountPwdError    ErrorCode = 102
	ErrorCodeAccountNotExist    ErrorCode = 103
	ErrorCodeAccountInsertError ErrorCode = 104

	// 200~300
	// 300~400
	ErrorCodeSystemDecodeError ErrorCode = 301
)

func NewErrorMsg(msg string, code ErrorCode) string {
	return fmt.Sprintf("%s error, error code: %d", msg, code)
}
