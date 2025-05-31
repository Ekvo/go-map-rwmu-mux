// реализация запросов
package transport

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ekvo/go-map-rwmu-mux/internal/db"
	"github.com/Ekvo/go-map-rwmu-mux/internal/service"
	"github.com/Ekvo/go-map-rwmu-mux/pkg/utils"
)

// добавление новой цитаты
// все хорошо -> возвращаем struct{}{}
func SaveOneQuote(usecase service.AddQuote) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("transport: SaveOneQuote member - {%s}, path - {%s};", r.Method, r.URL.Path)

		deserialize := service.NewQuoteDeserializer()
		if err := deserialize.Decode(r); err != nil {
			utils.EncodeJSON(w, http.StatusBadRequest, utils.NewCommonError(err))
			return
		}

		if err := usecase.CreateQuote(r.Context(), deserialize.Model()); err != nil {
			status := 0
			if errors.Is(err, db.ErrDBAlreadyExists) {
				status = http.StatusConflict
			} else {
				status = http.StatusInternalServerError
			}
			utils.EncodeJSON(w, status, utils.NewCommonError(err))
			return
		}

		utils.EncodeJSON(w, http.StatusCreated, struct{}{})
	}
}

// ищем случайную цитату
// нет ошибок -> возвращаем 'QuoteResponse'
func RetrieveRandomQuote(usecase service.FindRandomQuote) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("transport: RetrieveRandomQuote member - {%s}, path - {%s};", r.Method, r.URL.Path)

		quoteResponse, err := usecase.ReadRandomQuote(r.Context())
		if err != nil {
			status := 0
			if errors.Is(err, db.ErrDBEmpty) {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}

			utils.EncodeJSON(w, status, utils.NewCommonError(err))
			return
		}

		utils.EncodeJSON(w, http.StatusOK, quoteResponse)
	}
}

// полчение списка цитат
//
// 1. если параметр из url "author"
// запускае логику 'ReadQuoteListByAuthor'
//
// 2. нет параметра из url "author"
// запускае логику 'ReadQuoteList'
//
// нет ошибок -> возвращаем '[]QuoteResponse'
func RetrieveListOfQuote(usecase service.FindList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("transport: RetrieveListOfQuote member - {%s}, path - {%s};", r.Method, r.URL.Path)

		ctx := r.Context()
		param := r.URL.Query()

		if _, ex := param["author"]; ex {
			author := strings.TrimSpace(param.Get("author"))
			if author == "" {
				utils.EncodeJSON(w, http.StatusBadRequest, utils.NewCommonError(service.ErrServiceInvalidData))
				return
			}

			quotesResponse, err := usecase.ReadQuoteListByAuthor(ctx, author)
			if err != nil {
				status := 0
				if errors.Is(err, db.ErrDBNotFound) {
					status = http.StatusNotFound
				} else {
					status = http.StatusInternalServerError
				}

				utils.EncodeJSON(w, status, utils.NewCommonError(err))
				return
			}

			utils.EncodeJSON(w, http.StatusOK, quotesResponse)
			return
		}

		quotesResponse, err := usecase.ReadQuoteList(ctx)
		if err != nil {
			status := 0
			if errors.Is(err, db.ErrDBEmpty) {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}

			utils.EncodeJSON(w, status, utils.NewCommonError(err))
			return
		}

		utils.EncodeJSON(w, http.StatusOK, quotesResponse)
	}
}

// удаление цитаты по id полученого из пути url
func ExpelQuote(usecase service.RemoveQuote) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("transport: ExpelQuote member - {%s}, path - {%s};", r.Method, r.URL.Path)

		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Printf("transport: ExpelQuote index - {%s} not numeric;", idStr)
			utils.EncodeJSON(w, http.StatusBadRequest, utils.NewCommonError(service.ErrServiceInvalidData))
			return
		}
		if id == 0 {
			log.Print("transport: ExpelQuote index is zero")
			utils.EncodeJSON(w, http.StatusBadRequest, utils.NewCommonError(service.ErrServiceInvalidData))
			return
		}

		if err := usecase.DeleteQuote(r.Context(), uint(id)); err != nil {
			status := 0
			if errors.Is(err, db.ErrDBNotFound) {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}

			utils.EncodeJSON(w, status, utils.NewCommonError(err))
			return
		}

		utils.EncodeJSON(w, http.StatusOK, struct{}{})
	}
}
