package apperrors

type Code int

const (
	CodeUnknown Code = iota
	CodeInvalidArgument
	CodeNotFound
	CodeAlreadyExists
	CodeInternal
	CodeUnavailable
)
