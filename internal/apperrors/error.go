package apperrors

type SentinelError struct {
	code       Code
	messageKey string
	err        error
}

func New(code Code, messageKey string) *SentinelError {
	return &SentinelError{code: code, messageKey: messageKey}
}

func (e *SentinelError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.messageKey
}

func (e *SentinelError) Unwrap() error {
	return e.err
}

func (e *SentinelError) Code() int {
	return int(e.code)
}

func (e *SentinelError) MessageKey() string {
	return e.messageKey
}

func (e *SentinelError) WithErr(err error) *SentinelError {
	return &SentinelError{code: e.code, messageKey: e.messageKey, err: err}
}
