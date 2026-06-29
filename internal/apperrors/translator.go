package apperrors

const defaultLang = "ru"

var messages = map[string]map[string]string{
	"ru": {
		"search.empty_query":               "Поисковый запрос не может быть пустым",
		"search.invalid_pagination":        "Некорректные параметры пагинации",
		"index.no_chunks":                  "Список чанков для индексации пуст",
		"search.elasticsearch_unavailable": "Поисковый сервис временно недоступен",
		"search.history_not_found":         "История поиска не найдена",
		"common.internal_error":            "Внутренняя ошибка сервера",
	},
	"en": {
		"search.empty_query":               "Search query must not be empty",
		"search.invalid_pagination":        "Invalid pagination parameters",
		"index.no_chunks":                  "Chunk list for indexing is empty",
		"search.elasticsearch_unavailable": "Search service is temporarily unavailable",
		"search.history_not_found":         "Search history not found",
		"common.internal_error":            "Internal server error",
	},
}

type Translator struct{}

func NewTranslator() *Translator {
	return &Translator{}
}

func (t *Translator) Translate(messageKey, lang string) string {
	byLang, ok := messages[lang]
	if !ok {
		byLang = messages[defaultLang]
	}

	if msg, ok := byLang[messageKey]; ok {
		return msg
	}

	return messages[defaultLang]["common.internal_error"]
}
