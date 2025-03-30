package authlete

import (
	"fmt"
)

// AuthleteError は、Authleteクライアントで発生するエラーを表します。
type AuthleteError struct {
	Code    string
	Message string
	Err     error
}

func (e *AuthleteError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("AuthleteError: %s - %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("AuthleteError: %s - %s", e.Code, e.Message)
}
