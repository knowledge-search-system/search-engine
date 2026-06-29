package apperrors

import "testing"

func TestTranslatorTranslate(t *testing.T) {
	translator := NewTranslator()

	tests := []struct {
		name       string
		messageKey string
		lang       string
		want       string
	}{
		{name: "known key ru", messageKey: "search.empty_query", lang: "ru", want: "Поисковый запрос не может быть пустым"},
		{name: "known key en", messageKey: "search.empty_query", lang: "en", want: "Search query must not be empty"},
		{name: "unknown lang falls back to ru", messageKey: "search.empty_query", lang: "fr", want: "Поисковый запрос не может быть пустым"},
		{name: "unknown key falls back to internal error", messageKey: "does.not.exist", lang: "ru", want: "Внутренняя ошибка сервера"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := translator.Translate(tt.messageKey, tt.lang)
			if got != tt.want {
				t.Errorf("Translate(%q, %q) = %q, want %q", tt.messageKey, tt.lang, got, tt.want)
			}
		})
	}
}

func TestSentinelErrorCodeAndMessageKey(t *testing.T) {
	err := ErrEmptyQuery

	if err.Code() != int(CodeInvalidArgument) {
		t.Errorf("Code() = %d, want %d", err.Code(), CodeInvalidArgument)
	}
	if err.MessageKey() != "search.empty_query" {
		t.Errorf("MessageKey() = %q, want %q", err.MessageKey(), "search.empty_query")
	}

	wrapped := err.WithErr(errIO)
	if wrapped.Unwrap() != errIO {
		t.Errorf("Unwrap() did not return the wrapped error")
	}
}

var errIO = &customErr{"io failure"}

type customErr struct{ msg string }

func (e *customErr) Error() string { return e.msg }
