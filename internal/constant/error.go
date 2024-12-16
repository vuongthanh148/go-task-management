package constant

type ErrorCode int

const (
	ErrCodeInvalidRequest ErrorCode = iota + 1
	ErrCodeUnauthorized
	ErrCodeForbidden
	ErrCodeNotFound
	ErrCodeInternalServer
	ErrCodeConflict
	ErrCodeUpdateTaskStatus
	ErrCodeGetTasks
	ErrCodeGetTaskSummary
	ErrCodeGenerateToken
	ErrCodeDuplicateUser
	ErrCodeInvalidCredential
)

const (
	ErrMsgInvalidRequest    = "Invalid request payload"
	ErrMsgUnauthorized      = "Unauthorized"
	ErrMsgForbidden         = "Forbidden"
	ErrMsgNotFound          = "Not found"
	ErrMsgInternalServer    = "Internal server error"
	ErrMsgConflict          = "Conflict"
	ErrMsgUpdateTaskStatus  = "Failed to update task status"
	ErrMsgGetTasks          = "Failed to get tasks"
	ErrMsgGetTaskSummary    = "Failed to get task summary"
	ErrMsgGenerateToken     = "Failed to generate token"
	ErrMsgDuplicateUser     = "User existed"
	ErrMsgInvalidCredential = "Invalid credential"
)

func (e ErrorCode) String() string {

	var errorMessages = map[ErrorCode]string{
		ErrCodeInvalidRequest:    ErrMsgInvalidRequest,
		ErrCodeUnauthorized:      ErrMsgUnauthorized,
		ErrCodeForbidden:         ErrMsgForbidden,
		ErrCodeNotFound:          ErrMsgNotFound,
		ErrCodeInternalServer:    ErrMsgInternalServer,
		ErrCodeConflict:          ErrMsgConflict,
		ErrCodeUpdateTaskStatus:  ErrMsgUpdateTaskStatus,
		ErrCodeGetTasks:          ErrMsgGetTasks,
		ErrCodeGetTaskSummary:    ErrMsgGetTaskSummary,
		ErrCodeGenerateToken:     ErrMsgGenerateToken,
		ErrCodeDuplicateUser:     ErrMsgDuplicateUser,
		ErrCodeInvalidCredential: ErrMsgInvalidCredential,
	}

	return errorMessages[e]
}
