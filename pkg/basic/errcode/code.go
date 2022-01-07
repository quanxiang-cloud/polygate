package errcode

import (
	"github.com/quanxiang-cloud/cabin/error/errdefiner"
)

// error code
var (
	//--------------------------------------------------------------------------
	// Signature

	ErrInvalidURI               = errdefiner.MustReg(140010010013, "无效的URI：%v")
	ErrInputArgValidateMismatch = errdefiner.MustReg(140010010012, "参数验证失败：%v.%v")
	ErrInputValueExpired        = errdefiner.MustReg(140010010011, "参数已过期：%v.%v(%s)")
	ErrInputValueInvalid        = errdefiner.MustReg(140010010010, "参数值不合法：%v.%v")
	ErrInputMissingArg          = errdefiner.MustReg(140010010009, "缺少参数：%v.%v")

	//--------------------------------------------------------------------------

	ErrTimestampFormat   = errdefiner.MustReg(140010010008, "时间格式'%v'错误, 仅支持%v")
	ErrParameterError    = errdefiner.MustReg(140010010007, "输入参数错误：%v")
	ErrDataFormatInvalid = errdefiner.MustReg(140010010006, "格式不合法：%v.%v(%v)")
	ErrInternal          = errdefiner.MustReg(140010010005, "内部错误:%v")
	ErrSystemBusy        = errdefiner.MustReg(140010010004, "系统繁忙，请稍候重试")
	ErrGateBlockedAPI    = errdefiner.MustReg(140010010003, "API服务繁忙，请稍候重试")
	ErrGateBlockedIP     = errdefiner.MustReg(140010010002, "访问受限，请联系管理员")
	ErrGateError         = errdefiner.MustReg(140010010001, "网关错误 %v:%v")
)

// exports
const (
	ErrParams = errdefiner.ErrParams
	Internal  = errdefiner.Internal
	Unknown   = errdefiner.Unknown
	Success   = errdefiner.Success
)

// func exports
var (
	Errorf             = errdefiner.Errorf
	NewErrorWithString = errdefiner.NewErrorWithString
)
