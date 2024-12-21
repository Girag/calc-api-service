package application

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplication(t *testing.T) {
	testCases := []struct {
		request  string
		expected string
	}{
		{`{"expression":"5*(7+9)"}`, `{"result":80}`},
		{`{"expression":"(14+6)/5"}`, `{"result":4}`},
		{`{"expression":"3+3*3"}`, `{"result":12}`},
		{`{"expression":"(2+3)*4"}`, `{"result":20}`},
		{`{"expression":"7+7/7"}`, `{"result":8}`},
		{`{"expression":"1+2+3-4+5-6"}`, `{"result":1}`},
		{`{"expression":"50*3-30"}`, `{"result":120}`},
		{`{"expression":"16/2*4"}`, `{"result":32}`},
		{`{"expression":"15/3"}`, `{"result":5}`},
		{`{"expression":"10+20+33+40+55"}`, `{"result":158}`},
		{`{"expression":"22-(15-5)"}`, `{"result":12}`},
		{`{"expression":"3*5+(9+10)"}`, `{"result":34}`},
		{`{"expression":"7-4"}`, `{"result":3}`},
		{`{"expression":"22-2+5"}`, `{"result":25}`},
		{`{"expression":"5+5"}`, `{"result":10}`},
		{`{"expression":"11/11*4"}`, `{"result":4}`},
		{`{"expression":"(2+2)*2"}`, `{"result":8}`},
		{`{"expression":"100/10/2"}`, `{"result":5}`},
	}

	app := NewApp()
	handler := http.HandlerFunc(app.CalcHandler)
	for _, testCase := range testCases {
		t.Run(testCase.request, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, app.config.Path, bytes.NewBuffer([]byte(testCase.request)))
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusOK {
				t.Errorf("Wrong status code: %d, wanted 200", response.Code)
			} else if response.Body.String() != testCase.expected {
				t.Errorf("Wrong data: %s, wanted %s", response.Body.String(), testCase.expected)
			}
		})
	}
}

func TestApplicationCalcWithErrors(t *testing.T) {
	testCases := []string{
		`{"expression":"fgsfds"}`,
		`{"expression":"11^0"}`,
		`{"expression":"10/0"}`,
		`{"expression":"(23+32"}`,
		`{"expression":"32+23)"}`,
		`{"expression":"*6"}`,
		`{"expression":"0^11"}`,
		`{"expression":"5 + h"}`,
		`{"expression":"-2^4"}`,
		`{"expression":"(5+7*9"}`,
		`{"expression":"     "}`,
		`{"expression":"4+"}`,
		`{"expression":"9+t"}`,
		`{"expression":"6*"}`,
		`{"expression":"3+2*(1"}`,
		`{"expression":"7+"}`,
		`{"expression":")29+3("}`,
	}

	app := NewApp()
	handler := http.HandlerFunc(app.CalcHandler)
	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, app.config.Path, bytes.NewBuffer([]byte(testCase)))
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusUnprocessableEntity {
				t.Errorf("Wrong status code: %d, wanted 422", response.Code)
			}
		})
	}
}

func TestApplicationWithWrongBody(t *testing.T) {
	testCases := []string{
		`5+5`,
		`{"expression":true}`,
		`{"expression":[4,2]}`,
		`true`,
		`{"expression": "5+5" ,}`,
		`{"expression": 5+5}`,
		`{"expression":"5 + 5",}`,
		`{"expression": []}`,
		`[9,8,7]`,
		`{"expression":"5+5", some: "another"}`,
		`{"expression":null}`,
		`{"a":33, "b":44}`,
		`"5+5"`,
		`{"expression": "5+5", "}`,
		`{"expression" : 5+5}`,
		`{"expression":5+5"}`,
		`{"expression":5 + 5}`,
		`5`,
		`{"expression":,`,
		`null`,
		`{"expression": 5+5}`,
		`{"expression": "5+5",}`,
		`[]`,
		`{},`,
		`{"expression": "5+5" ,},`,
		`{"expression": {},}`,
		`{"exp":"5+5"}`,
		`{expression:"5 + 5"},`,
		`{"expression":},`,
		`{}`,
	}

	app := NewApp()
	handler := http.HandlerFunc(app.CalcHandler)
	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, app.config.Path, bytes.NewBuffer([]byte(testCase)))
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Errorf("Wrong status code: %d, wanted 400", response.Code)
			}
		})
	}
}

func TestApplicationWithWrongPath(t *testing.T) {
	app := NewApp()
	handler := http.HandlerFunc(app.CalcHandler)
	request := httptest.NewRequest(http.MethodPost, "/cucumber", bytes.NewBuffer([]byte("")))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Errorf("Wrong status code: %d, wanted 404", response.Code)
	}
}

func TestApplicationWithWrongMethod(t *testing.T) {
	app := NewApp()
	handler := http.HandlerFunc(app.CalcHandler)
	request := httptest.NewRequest(http.MethodGet, app.config.Path, bytes.NewBuffer([]byte("")))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf("Wrong status code: %d, wanted 405", response.Code)
	}
}
