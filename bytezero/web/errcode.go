package web

type ErrCode int32

const (
	ErrCode_success                               ErrCode = 0
	ErrCode_error                                 ErrCode = 1
	ErrCode_http                                  ErrCode = 100000
	ErrCode_ec_module_action_miss                 ErrCode = 100001
)
