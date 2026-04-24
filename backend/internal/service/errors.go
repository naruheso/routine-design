package service

import "errors"

// センチネルエラー定義
var (
	ErrNotFound      = errors.New("not found")
	ErrValidation    = errors.New("validation error")
	ErrConflict      = errors.New("conflict")
	ErrForbidden     = errors.New("forbidden")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrStatusConflict = errors.New("この申請は既に他のユーザーによって処理されています")
)
