package models

import "errors"

var ErrInvalidInput = errors.New("invalid input")

const (
	ErrUserAlreadyExists = "User already exists"
	ErrUserDoesNotExist  = "User does not exist"
	ErrUserBlocked       = "User is blocked"
	ErrInvalidID         = "Unauthorized or invalid user ID"

	MsgLoginSuccessful           = "Login successful"
	MsgLogoutSuccessful          = "Logout successful"
	MsgSignupSuccessful          = "User signed up successfully!"
	MsgEmailVerifiedSuccessfully = "Email verified successfully"
	MsgVerificationEmailResent   = "Verification email resent"
	MsgPasswordResetEmailSent    = "Password reset email sent"
	MsgPasswordResetSuccessfully = "Password reset successfully"

	MsgProfileUpdatedSuccessfully = "Profile updated successfully"
	MsgProfilePictureUploaded     = "Profile picture uploaded successfully"

	ErrRequiredFieldsEmpty = "Required fields cannot be empty"
	ErrInvalidEmailFormat  = "Invalid email format"
	ErrNegativeAge         = "Age must be positive"
	ErrPasswordComplexity  = "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
	ErrPasswordLength      = "Password must be between %d and %d characters"
	ErrInvalidPhoneNumber  = "Invalid phone number format"

	MinPasswordLength = 8
	MaxPasswordLength = 72
)
