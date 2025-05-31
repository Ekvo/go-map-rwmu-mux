package utils

import (
	"encoding/json"
	"errors"
	"log"
	"mime"
	"net/http"
)

var (
	// ErrUtilsInvalidMedia - неправильный тип media в запросе
	ErrUtilsInvalidMedia = errors.New("unexpected media type")

	// ErrUtilsEmptyBody - пустое тело в запросе, используется в 'DecodeJSON'
	ErrUtilsEmptyBody = errors.New("request body is empty")
)

// CommonError - для записи ошибок в 'http.ResponseWriter'
type CommonError struct {
	Message string `json:"error"`
}

func NewCommonError(err error) *CommonError {
	return &CommonError{Message: err.Error()}
}

// DecodeJSON - получение объекта 'obj' из запроса:
// проверка "Content-Type", тела ответа, запретить неизвестные поля 'DisallowUnknownFields'.
func DecodeJSON(req *http.Request, obj any) error {
	media := req.Header.Get("Content-Type")
	parse, _, err := mime.ParseMediaType(media)
	if err != nil || parse != "application/json" {
		return ErrUtilsInvalidMedia
	}

	if req.Body == nil {
		return ErrUtilsEmptyBody
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			log.Printf("utils: DecodeJSON req.Body.Close error - {%v};", err)
		}
	}()
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(obj)
}

// EncodeJSON - запись объекта в 'http.ResponseWriter',
// добавление статуса ответа.
// Задаем через Header - "Content-Type".
func EncodeJSON(w http.ResponseWriter, status int, obj any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		log.Printf("utils: EncodeJSON json.Encode error - {%v};", err)
	}
}

// IndexByValue - бинарный поиск индекса по значению 'target';
// если найден, возвращает индекс и true, в противном случае - 0 и false.
func IndexByValue(slice []uint, target uint) (uint, bool) {
	low := uint(0)
	high := uint(len(slice) - 1)

	for low <= high {
		mid := (low + high) / 2
		guess := slice[mid]

		if guess == target {
			// элемент найден
			return mid, true
		} else if guess > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return 0, false // элемент не найден
}
