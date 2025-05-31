package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_DecodeJSON_EncodeJSON(t *testing.T) {
	type user struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Age     uint   `json:"age"`
	}

	dataForRequest := []struct {
		description string
		body        string
		resCode     int
		resBody     string
		msg         string
	}{
		{
			description: `valid user`,
			body:        `{"name":"Alex","surname":"","age":26}`,
			resCode:     http.StatusOK,
			resBody:     "{\"user\":\"approve\"}\n",
			msg:         `valid Decode and Encode`,
		},
		{
			description: `wrong user`,
			body:        `{"name":"Alex","surname":"","age":26,"avp":"alien"}`,
			resCode:     http.StatusBadRequest,
			resBody:     "{\"error\":\"json: unknown field \\\"avp\\\"\"}\n",
			msg:         `invalid Decode and valid Encode`,
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /test", func(w http.ResponseWriter, r *http.Request) {
		var u user
		if err := DecodeJSON(r, &u); err != nil {
			EncodeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		EncodeJSON(w, http.StatusOK, map[string]string{"user": "approve"})
	})

	for _, test := range dataForRequest {
		t.Run(test.description, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte(test.body)))
			if err != nil {
				t.Fatalf("http.NewRequest error - {%v};", err)
			}

			req.Header.Set("Content-Type", "application/json; charset=UTF-8")

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != test.resCode {
				t.Errorf("http status code should be equal got:want {%d}:{%d}", w.Code, test.resCode)
			}

			if w.Body == nil {
				t.Fatal("w.Body is nil")
			}
			
			gotBody := w.Body.String()
			wantBody := test.resBody

			if gotBody != wantBody {
				t.Errorf("body got:want {%s}:{%s}", gotBody, wantBody)
			}
		})
	}
}

func Test_IndexByValue(t *testing.T) {
	nums := []uint{0, 5, 8, 19, 234, 567, 890, 1234}

	testData := []struct {
		title     string
		num       uint
		want      uint
		wantExist bool
	}{
		{
			title:     `find index for 5`,
			num:       5,
			want:      1,
			wantExist: true,
		},
		{
			title:     `find index for 0`,
			num:       0,
			want:      0,
			wantExist: true,
		},
		{
			title:     `find index for 19`,
			num:       19,
			want:      3,
			wantExist: true,
		},
		{
			title:     `find index for 567`,
			num:       567,
			want:      5,
			wantExist: true,
		},
		{
			title:     `find index for 1234`,
			num:       1234,
			want:      7,
			wantExist: true,
		},
		{
			title:     `find index for 789`,
			num:       789,
			want:      0,
			wantExist: false,
		},
	}

	for _, test := range testData {
		t.Run(test.title, func(t *testing.T) {
			got, gotExist := IndexByValue(nums, test.num)

			if gotExist != test.wantExist {
				t.Errorf("IndexByValue: gotExist - {%t} not equal wantExist - {%t};", gotExist, test.wantExist)
			}

			if got != test.want {
				t.Errorf("IndexByValue: got - {%d} not equal wantExist - {%d};", got, test.want)
			}
		})
	}
}
