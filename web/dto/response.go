package dto

import "xll-job/web/constant"

type Response[T any] struct {
	Code int32  `json:"code" form:"code" json:"code" uri:"code" xml:"code" yaml:"code" `
	Msg  string `json:"msg" form:"msg" json:"msg" uri:"msg" xml:"msg" yaml:"msg" `
	Data T      `json:"data" form:"data" json:"data" uri:"data" xml:"data" yaml:"data" `
}

func NewOkResponse[T any](t T) *Response[T] {
	resp := &Response[T]{
		Code: constant.HttpOk,
		Msg:  "ok",
		Data: t,
	}
	return resp
}
