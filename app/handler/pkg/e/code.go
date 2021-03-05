// 本包用于定义业务错误。
// 设计原则参见：[谷歌API设计指南|错误](https://www.bookstack.cn/read/API-design-guide/API-design-guide-07-%E9%94%99%E8%AF%AF.md)

package e

import (
	"go.uber.org/zap/zapcore"
	"net/http"
)

// Code 为业务错误码。注意：该错误码面向客户端，而非面向日志！
type Code uint8

const (
	// Not an error; returned on success
	//
	// 成功，没有错误
	//
	// HTTP Mapping: 200 OK
	CodeOK Code = iota

	// The operation was cancelled, typically by the caller.
	//
	// 客户端取消请求
	//
	// HTTP Mapping: 499 Client Closed Request
	CodeCancelled

	// Unknown error.  For example, this error may be returned when
	// a `Status` value received from another address space belongs to
	// an error space that is not known in this address space.  Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	//
	// 未知的服务器错误。 通常是服务器错误。
	//
	// HTTP Mapping: 500 Internal Server Error
	CodeUnknown

	// The client specified an invalid argument.  Note that this differs
	// from `FAILED_PRECONDITION`.  `INVALID_ARGUMENT` indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	//
	// 客户端指定了无效的参数。 检查错误消息和错误详细信息以获取更多信息。
	//
	// HTTP Mapping: 400 Bad Request
	CodeInvalidArgument

	// The deadline expired before the operation could complete. For operations
	// that change the state of the system, this error may be returned
	// even if the operation has completed successfully.  For example, a
	// successful response from a server could have been delayed long
	// enough for the deadline to expire.
	//
	// 已超过请求期限。如果重复发生，请考虑降低请求的复杂性。
	//
	// HTTP Mapping: 504 Gateway Timeout
	CodeDeadlineExceeded

	// Some requested entity (e.g., file or directory) was not found.
	//
	// Note to server developers: if a request is denied for an entire class
	// of users, such as gradual feature rollout or undocumented whitelist,
	// `NOT_FOUND` may be used. If a request is denied for some users within
	// a class of users, such as user-based access control, `PERMISSION_DENIED`
	// must be used.
	//
	// 找不到指定的资源，或者该请求被未公开的原因（例如白名单）拒绝。
	//
	// HTTP Mapping: 404 Not Found
	CodeNotFound

	// The entity that a client attempted to create (e.g., file or directory)
	// already exists.
	//
	// 客户端尝试创建的资源已存在。
	//
	// HTTP Mapping: 409 Conflict
	CodeAlreadyExists

	// The caller does not have permission to execute the specified
	// operation. `PERMISSION_DENIED` must not be used for rejections
	// caused by exhausting some resource (use `RESOURCE_EXHAUSTED`
	// instead for those errors). `PERMISSION_DENIED` must not be
	// used if the caller can not be identified (use `UNAUTHENTICATED`
	// instead for those errors). This error code does not imply the
	// request is valid or the requested entity exists or satisfies
	// other pre-conditions.
	//
	// 客户端没有足够的权限。
	// 这可能是因为OAuth令牌没有正确的范围，客户端没有权限，或者客户端项目尚未启用API。
	//
	// HTTP Mapping: 403 Forbidden
	CodePermissionDenied

	// The request does not have valid authentication credentials for the
	// operation.
	//
	// 由于遗失，无效或过期的OAuth令牌而导致请求未通过身份验证。
	//
	// HTTP Mapping: 401 Unauthorized
	CodeUnauthenticated

	// Some resource has been exhausted, perhaps a per-user quota, or
	// perhaps the entire file system is out of space.
	//
	// 资源配额达到速率限制。
	// 客户端应该查找 QuotaFailure 错误详细信息以获取更多信息。
	//
	// HTTP Mapping: 429 Too Many Requests
	CodeResourceExhausted

	// The operation was rejected because the system is not in a state
	// required for the operation's execution.  For example, the directory
	// to be deleted is non-empty, an rmdir operation is applied to
	// a non-directory, etc.
	//
	// Service implementors can use the following guidelines to decide
	// between `FAILED_PRECONDITION`, `ABORTED`, and `UNAVAILABLE`:
	//  (a) Use `UNAVAILABLE` if the client can retry just the failing call.
	//  (b) Use `ABORTED` if the client should retry at a higher level
	//      (e.g., when a client-specified test-and-set fails, indicating the
	//      client should restart a read-modify-write sequence).
	//  (c) Use `FAILED_PRECONDITION` if the client should not retry until
	//      the system state has been explicitly fixed.  E.g., if an "rmdir"
	//      fails because the directory is non-empty, `FAILED_PRECONDITION`
	//      should be returned since the client should not retry unless
	//      the files are deleted from the directory.
	//
	// 请求不能在当前系统状态下执行，例如删除非空目录。
	//
	// HTTP Mapping: 400 Bad Request
	CodeFailedPrecondition

	// The operation was aborted, typically due to a concurrency issue such as
	// a sequencer check failure or transaction abort.
	//
	// See the guidelines above for deciding between `FAILED_PRECONDITION`,
	// `ABORTED`, and `UNAVAILABLE`.
	//
	// 并发冲突，例如读-修改-写冲突。
	//
	// HTTP Mapping: 409 Conflict
	CodeAborted

	// The operation was attempted past the valid range.  E.g., seeking or
	// reading past end-of-file.
	//
	// Unlike `INVALID_ARGUMENT`, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate `INVALID_ARGUMENT` if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// `OUT_OF_RANGE` if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between `FAILED_PRECONDITION` and
	// `OUT_OF_RANGE`.  We recommend using `OUT_OF_RANGE` (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an `OUT_OF_RANGE` error to detect when
	// they are done.
	//
	// 客户端指定了无效的范围。
	//
	// HTTP Mapping: 400 Bad Request
	CodeOutOfRange

	// The operation is not implemented or is not supported/enabled in this
	// service.
	//
	// 服务器未实现该API方法。
	//
	// HTTP Mapping: 501 Not Implemented
	CodeUnimplemented

	// Internal errors.  This means that some invariants expected by the
	// underlying system have been broken.  This error code is reserved
	// for serious errors.
	//
	// 内部服务错误。 通常是服务器错误。
	//
	// HTTP Mapping: 500 Internal Server Error
	CodeInternal

	// The service is currently unavailable.  This is most likely a
	// transient condition, which can be corrected by retrying with
	// a backoff. Note that it is not always safe to retry
	// non-idempotent operations.
	//
	// See the guidelines above for deciding between `FAILED_PRECONDITION`,
	// `ABORTED`, and `UNAVAILABLE`.
	//
	// 暂停服务。通常是服务器已经关闭。
	//
	// HTTP Mapping: 503 Service Unavailable
	CodeUnavailable

	// Unrecoverable data loss or corruption.
	//
	// 不可恢复的数据丢失或数据损坏。 客户端应该向用户报告错误。
	//
	// HTTP Mapping: 500 Internal Server Error
	CodeDataLoss
)

// CodeDetail 类型定义业务错误码 Code 对应的详细信息
type CodeDetail struct {
	Code       Code          // 错误码
	Status     string        // 错误码名称
	Message    string        // 错误码描述
	HttpStatus int           // 错误码对应的 http status
	LogLevel   zapcore.Level // 错误码对应的 http 日志级别
}

// codeDetails 定义业务错误码 Code 对错误码详情 CodeDetail 的映射
var codeDetails map[Code]CodeDetail = map[Code]CodeDetail{
	CodeOK: {
		Code:       CodeOK,
		Status:     "OK",
		Message:    "成功",
		HttpStatus: http.StatusOK, // 200
		LogLevel:   zapcore.InfoLevel,
	},
	CodeCancelled: {
		Code:       CodeCancelled,
		Status:     "CANCELLED",
		Message:    "客户端取消请求",
		HttpStatus: 499,
		LogLevel:   zapcore.WarnLevel,
	},
	CodeUnknown: {
		Code:       CodeUnknown,
		Status:     "UNKNOWN",
		Message:    "未知的服务器错误",
		HttpStatus: http.StatusInternalServerError, // 500
		LogLevel:   zapcore.ErrorLevel,
	},
	CodeInvalidArgument: {
		Code:       CodeInvalidArgument,
		Status:     "INVALID_ARGUMENT",
		Message:    "客户端指定了无效的参数",
		HttpStatus: http.StatusBadRequest, // 400
		LogLevel:   zapcore.WarnLevel,
	},
	CodeDeadlineExceeded: {
		Code:       CodeDeadlineExceeded,
		Status:     "DEADLINE_EXCEEDED",
		Message:    "已超过请求期限",
		HttpStatus: http.StatusGatewayTimeout, // 504
		LogLevel:   zapcore.ErrorLevel,
	},
	CodeNotFound: {
		Code:       CodeNotFound,
		Status:     "NOT_FOUND",
		Message:    "找不到指定的资源",
		HttpStatus: http.StatusNotFound, // 404
		LogLevel:   zapcore.WarnLevel,
	},
	CodeAlreadyExists: {
		Code:       CodeAlreadyExists,
		Status:     "ALREADY_EXISTS",
		Message:    "客户端尝试创建的资源已存在",
		HttpStatus: http.StatusConflict, // 409
		LogLevel:   zapcore.WarnLevel,
	},
	CodePermissionDenied: {
		Code:       CodePermissionDenied,
		Status:     "PERMISSION_DENIED",
		Message:    "客户端没有足够的权限",
		HttpStatus: http.StatusForbidden, // 403
		LogLevel:   zapcore.WarnLevel,
	},
	CodeUnauthenticated: {
		Code:       CodeUnauthenticated,
		Status:     "UNAUTHENTICATED",
		Message:    "请求未通过身份验证",
		HttpStatus: http.StatusUnauthorized, // 401
		LogLevel:   zapcore.WarnLevel,
	},
	CodeResourceExhausted: {
		Code:       CodeResourceExhausted,
		Status:     "RESOURCE_EXHAUSTED",
		Message:    "资源配额达到速率限制",
		HttpStatus: http.StatusTooManyRequests, // 429
		LogLevel:   zapcore.WarnLevel,
	},
	CodeFailedPrecondition: {
		Code:       CodeFailedPrecondition,
		Status:     "FAILED_PRECONDITION",
		Message:    "请求不能在当前系统状态下执行",
		HttpStatus: http.StatusBadRequest, // 400
		LogLevel:   zapcore.WarnLevel,
	},
	CodeAborted: {
		Code:       CodeAborted,
		Status:     "ABORTED",
		Message:    "并发冲突",
		HttpStatus: http.StatusConflict, // 409
		LogLevel:   zapcore.WarnLevel,
	},
	CodeOutOfRange: {
		Code:       CodeOutOfRange,
		Status:     "OUT_OF_RANGE",
		Message:    "客户端指定了无效的范围",
		HttpStatus: http.StatusBadRequest, // 400
		LogLevel:   zapcore.WarnLevel,
	},
	CodeUnimplemented: {
		Code:       CodeUnimplemented,
		Status:     "NOT_IMPLEMENTED",
		Message:    "服务器未实现该API方法",
		HttpStatus: http.StatusNotImplemented, // 501
		LogLevel:   zapcore.WarnLevel,
	},
	CodeInternal: {
		Code:       CodeInternal,
		Status:     "INTERNAL",
		Message:    "内部服务错误",
		HttpStatus: http.StatusInternalServerError, // 500
		LogLevel:   zapcore.ErrorLevel,
	},
	CodeUnavailable: {
		Code:       CodeUnavailable,
		Status:     "UNAVAILABLE",
		Message:    "暂停服务",
		HttpStatus: http.StatusServiceUnavailable, // 503
		LogLevel:   zapcore.WarnLevel,
	},
	CodeDataLoss: {
		Code:       CodeDataLoss,
		Status:     "DATA_LOSS",
		Message:    "不可恢复的数据丢失或数据损坏",
		HttpStatus: http.StatusInternalServerError, // 500
		LogLevel:   zapcore.ErrorLevel,
	},
}

// GetCodeDetail 通过业务错误码 Code 获取对应的业务错误详情 CodeDetail
func GetCodeDetail(code Code) CodeDetail {
	return codeDetails[code]
}
