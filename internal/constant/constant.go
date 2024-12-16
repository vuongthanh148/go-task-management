package constant

type SuccessMessage string

const (
	SuccessMessage_Created       SuccessMessage = "successfully created"
	SuccessMessage_Updated       SuccessMessage = "successfully updated"
	SuccessMessage_Deleted       SuccessMessage = "successfully deleted"
	SuccessMessage_ResetPassword SuccessMessage = "A password reset email has already been sent. Please check your inbox and follow the instructions to reset your password."

	SuccessCode                         = "0000"
	BusinessErrorCode_InvalidInputToMap = "9999"
)
