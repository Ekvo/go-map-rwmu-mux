package transport

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ekvo/go-map-rwmu-mux/internal/config"
	"github.com/Ekvo/go-map-rwmu-mux/internal/db"
	"github.com/Ekvo/go-map-rwmu-mux/internal/model"
	"github.com/Ekvo/go-map-rwmu-mux/internal/service"
)

// набор цитат для записи в базу
// изменять лишь в случае полного понимания работы тестов в данном пакете
var quotesData = []model.Quote{
	{
		Author: `william james`,
		Body:   `the greatest weapon against stress is our ability to choose one thought over another`,
	},
	{
		Author: `napoleon bonaparte`,
		Body:   `my dictionary does not contain the word 'impossible'`,
	},
	{
		Author: `steve jobs`,
		Body:   `your time is limited, so don’t waste it living someone else’s life`,
	},
}

func Test_SaveOneQuote(t *testing.T) {
	testData := []struct {
		title              string
		datasForRequest    string
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			title:              `valid add quote`,
			datasForRequest:    `{"author":"William James","quote":"The greatest weapon against stress is our ability to choose one thought over another"}`,
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   "{}\n",
		},
		{
			title:              `invalid add quote (already exist)`,
			datasForRequest:    `{"author":"William James","quote":"The greatest weapon against stress is our ability to choose one thought over another"}`,
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   "{\"error\":\"quote already exists\"}\n",
		},
		{
			title:              `wrong quote, field "quote"  is empty`,
			datasForRequest:    `{"author":"William James","quote":"   "}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "{\"error\":\"invalid data\"}\n",
		},
		{
			title:              `dangerous) quote, alien field`,
			datasForRequest:    `{"author":"Pradator","quote":"Wins","Alien":"Not this time"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "{\"error\":\"json: unknown field \\\"Alien\\\"\"}\n",
		},
	}

	store := db.NewProvider()
	usecase := service.NewService(store)
	r := NewTransport(&config.Config{"not host", "not port"})
	r.Routes(usecase)

	for _, test := range testData {
		t.Run(test.title, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, `/quotes`, bytes.NewBuffer([]byte(test.datasForRequest)))
			if err != nil {
				t.Fatalf("http.NewRequest error - {%v};", err)
			}

			req.Header.Set("Content-Type", "application/json; charset=UTF-8")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != test.expectedStatusCode {
				t.Errorf("status not equal {got}:{want} {%d}:{%d}", w.Code, test.expectedStatusCode)
			}

			if w.Body != nil {
				if w.Body.String() != test.expectedResponse {
					t.Errorf("invalid response body {got}:{want} {%s}:{%s};", w.Body.String(), test.expectedResponse)
				}
			} else if test.expectedResponse != "" {
				t.Errorf("w.Body is nil, but expectedResponse - {%s};", test.expectedResponse)
			}
		})
	}
}

func Test_RetrieveRandomQuote(t *testing.T) {
	testData := []struct {
		title              string
		expectedStatusCode int
		expectedResponse   string
		testLogic          func() (*httptest.ResponseRecorder, error)
	}{
		{
			title:              `work with empty base`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   "{\"error\":\"quote list is empty\"}\n",
			testLogic: func() (*httptest.ResponseRecorder, error) {
				store := db.NewProvider()
				usecase := service.NewService(store)
				r := NewTransport(&config.Config{"not host", "not port"})
				r.Routes(usecase)

				req, err := http.NewRequest(http.MethodGet, `/quotes/random`, nil)
				if err != nil {
					return nil, err
				}

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				return w, nil
			},
		},
		{
			title:              `valid request for find random quote`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "{\"id\":\"1\",\"author\":\"steve jobs\",\"quote\":\"your time is limited, so don’t waste it living someone else’s life\"}\n",
			testLogic: func() (*httptest.ResponseRecorder, error) {
				store := db.NewProvider()

				// добавляем запись для поиска
				if err := store.NewQuote(
					context.TODO(),
					quotesData[2],
				); err != nil {
					return nil, err
				}

				usecase := service.NewService(store)
				r := NewTransport(&config.Config{"not host", "not port"})
				r.Routes(usecase)

				req, err := http.NewRequest(http.MethodGet, `/quotes/random`, nil)
				if err != nil {
					return nil, err
				}

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				return w, nil
			},
		},
	}

	for _, test := range testData {
		t.Run(test.title, func(t *testing.T) {
			w, err := test.testLogic()
			if err != nil {
				t.Fatalf("testLogic return error - {%v};", err)
			}

			if w.Code != test.expectedStatusCode {
				t.Errorf("status not equal {got}:{want} {%d}:{%d}", w.Code, test.expectedStatusCode)
			}

			if w.Body != nil {
				if w.Body.String() != test.expectedResponse {
					t.Errorf("invalid response body {got}:{want} {%s}:{%s};", w.Body.String(), test.expectedResponse)
				}
			} else if test.expectedResponse != "" {
				t.Errorf("w.Body is nil, but expectedResponse - {%s};", test.expectedResponse)
			}
		})
	}

}

func Test_RetrieveListOfQuote(t *testing.T) {
	testData := []struct {
		title              string
		expectedStatusCode int
		expectedResponse   string
		testLogic          func() (*httptest.ResponseRecorder, error)
	}{
		{
			title:              `work with empty base`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"error":"quote list is empty"}`,
			testLogic: func() (*httptest.ResponseRecorder, error) {
				store := db.NewProvider()
				usecase := service.NewService(store)
				r := NewTransport(&config.Config{"not host", "not port"})
				r.Routes(usecase)

				req, err := http.NewRequest(http.MethodGet, `/quotes`, nil)
				if err != nil {
					return nil, err
				}

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				return w, nil
			},
		},
		{
			title:              `work with empty base, url - contain author`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"error":"quote not found"}`,
			testLogic: func() (*httptest.ResponseRecorder, error) {
				store := db.NewProvider()
				usecase := service.NewService(store)
				r := NewTransport(&config.Config{"not host", "not port"})
				r.Routes(usecase)

				req, err := http.NewRequest(http.MethodGet, `/quotes?author=Alex`, nil)
				if err != nil {
					return nil, err
				}

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				return w, nil
			},
		},
		{
			title:              `valid request find list quote`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: `[{"id":"1","author":"william james","quote":"the greatest weapon against stress is our ability to choose one thought over another"},
{"id":"2","author":"napoleon bonaparte","quote":"my dictionary does not contain the word 'impossible'"},
{"id":"3","author":"steve jobs","quote":"your time is limited, so don’t waste it living someone else’s life"}]`,
			testLogic: func() (*httptest.ResponseRecorder, error) {
				store := db.NewProvider()

				// заполняем базу
				for _, quote := range quotesData {
					if err := store.NewQuote(context.TODO(), quote); err != nil {
						return nil, err
					}
				}

				usecase := service.NewService(store)
				r := NewTransport(&config.Config{"not host", "not port"})
				r.Routes(usecase)

				req, err := http.NewRequest(http.MethodGet, `/quotes`, nil)
				if err != nil {
					return nil, err
				}

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				return w, nil
			},
		},
		{
			title:              `valid request find list quote by author`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: `[{"id":"1","author":"william james","quote":"the greatest weapon against stress is our ability to choose one thought over another"}]
`,
			testLogic: func() (*httptest.ResponseRecorder, error) {
				store := db.NewProvider()

				for _, quote := range quotesData {
					if err := store.NewQuote(context.TODO(), quote); err != nil {
						return nil, err
					}
				}

				usecase := service.NewService(store)
				r := NewTransport(&config.Config{"not host", "not port"})
				r.Routes(usecase)

				req, err := http.NewRequest(http.MethodGet, `/quotes?author=William James`, nil)
				if err != nil {
					return nil, err
				}

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				return w, nil
			},
		},
	}

	for _, test := range testData {
		t.Run(test.title, func(t *testing.T) {
			w, err := test.testLogic()
			if err != nil {
				t.Fatalf("testLogic return error - {%v};", err)
			}

			if w.Code != test.expectedStatusCode {
				t.Errorf("status not equal {got}:{want} {%d}:{%d}", w.Code, test.expectedStatusCode)
			}

			if w.Body != nil {
				/*
					* add \n because

					* package json
					  func (enc *Encoder) Encode(v any) error {
							...
							e := newEncodeState()
							...
							e.WriteByte('\n')
							...
					}
				*/
				got := w.Body.String()
				want := strings.Replace(test.expectedResponse, "\n", "", -1) + "\n"
				if got != want {
					t.Errorf("invalid response body {got}:{want} {%s}:{%s};", got, want)
				}
			} else if test.expectedResponse != "" {
				t.Errorf("w.Body is nil, but expectedResponse - {%s};", test.expectedResponse)
			}
		})
	}
}

func Test_ExpelQuote(t *testing.T) {
	testData := []struct {
		title              string
		expectedStatusCode int
		expectedResponse   string
		testLogic          func() (*httptest.ResponseRecorder, error)
	}{
		{
			title:              `wrong delete, base is empty`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   "{\"error\":\"quote not found\"}\n",
			testLogic: func() (*httptest.ResponseRecorder, error) {
				store := db.NewProvider()
				usecase := service.NewService(store)
				r := NewTransport(&config.Config{"not host", "not port"})
				r.Routes(usecase)

				req, err := http.NewRequest(http.MethodDelete, `/quotes/1`, nil)
				if err != nil {
					return nil, err
				}

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				return w, nil
			},
		},
		{
			title:              `wrong delete, id not numeric`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "{\"error\":\"invalid data\"}\n",
			testLogic: func() (*httptest.ResponseRecorder, error) {
				store := db.NewProvider()
				usecase := service.NewService(store)
				r := NewTransport(&config.Config{"not host", "not port"})
				r.Routes(usecase)

				req, err := http.NewRequest(http.MethodDelete, `/quotes/{what}`, nil)
				if err != nil {
					return nil, err
				}

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				return w, nil
			},
		},
		{
			title:              `create quote and then delete`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "{}\n",
			testLogic: func() (*httptest.ResponseRecorder, error) {
				store := db.NewProvider()

				for _, quote := range quotesData {
					if err := store.NewQuote(context.TODO(), quote); err != nil {
						return nil, err
					}
				}

				usecase := service.NewService(store)
				r := NewTransport(&config.Config{"not host", "not port"})
				r.Routes(usecase)

				req, err := http.NewRequest(http.MethodDelete, `/quotes/1`, nil)
				if err != nil {
					return nil, err
				}

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				return w, nil
			},
		},
	}

	for _, test := range testData {
		t.Run(test.title, func(t *testing.T) {
			w, err := test.testLogic()
			if err != nil {
				t.Fatalf("testLogic return error - {%v};", err)
			}

			if w.Code != test.expectedStatusCode {
				t.Errorf("status not equal {got}:{want} {%d}:{%d}", w.Code, test.expectedStatusCode)
			}

			if w.Body != nil {
				got := w.Body.String()
				want := test.expectedResponse

				if got != want {
					t.Errorf("invalid response body {got}:{want} {%s}:{%s};", got, want)
				}
			} else if test.expectedResponse != "" {
				t.Errorf("w.Body is nil, but expectedResponse - {%s};", test.expectedResponse)
			}
		})
	}
}
