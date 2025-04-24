package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	testSuite := []struct {
		name     string
		reqBody  string
		want     string
		wantCode int
	}{
		{
			name:     "basic case",
			reqBody:  `[29, 8, 3, 4]`,
			want:     "Odd numbers: [3 29]\nEven numbers: [4 8]\n",
			wantCode: http.StatusOK,
		},
		{
			name:     "empty array",
			reqBody:  `[]`,
			want:     "Odd numbers: []\nEven numbers: []\n",
			wantCode: http.StatusOK,
		},
		{
			name:     "only even numbers",
			reqBody:  `[2, 4, 6, 8]`,
			want:     "Odd numbers: []\nEven numbers: [2 4 6 8]\n",
			wantCode: http.StatusOK,
		},
		{
			name:     "only odd numbers",
			reqBody:  `[1, 3, 5, 7]`,
			want:     "Odd numbers: [1 3 5 7]\nEven numbers: []\n",
			wantCode: http.StatusOK,
		},
		{
			name:     "mixed even and odd with negatives",
			reqBody:  `[-2, -3, 0, 1, 4, -5]`,
			want:     "Odd numbers: [-5 -3 1]\nEven numbers: [-2 0 4]\n",
			wantCode: http.StatusOK,
		},
		{
			name:     "large numbers",
			reqBody:  `[2000000000, 3000000001]`,
			want:     "Odd numbers: [3000000001]\nEven numbers: [2000000000]\n",
			wantCode: http.StatusOK,
		},
		{
			name:     "single odd number",
			reqBody:  `[9]`,
			want:     "Odd numbers: [9]\nEven numbers: []\n",
			wantCode: http.StatusOK,
		},
		{
			name:     "single even number",
			reqBody:  `[10]`,
			want:     "Odd numbers: []\nEven numbers: [10]\n",
			wantCode: http.StatusOK,
		},
		{
			name:     "all zeroes",
			reqBody:  `[0, 0, 0]`,
			want:     "Odd numbers: []\nEven numbers: [0 0 0]\n",
			wantCode: http.StatusOK,
		},
		{
			name:     "repeated numbers",
			reqBody:  `[1, 1, 1, 2, 2, 2]`,
			want:     "Odd numbers: [1 1 1]\nEven numbers: [2 2 2]\n",
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid JSON",
			reqBody:  `[1, 2,`,
			want:     "Bad request\n",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid content",
			reqBody:  `{[1.1, 1.2, 2, 4]`,
			want:     "Bad request\n",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, singleTest := range testSuite {
		t.Run(singleTest.name, func(t *testing.T) {
			simRequest, err := http.NewRequest(http.MethodPut, "/numbers/sorted", strings.NewReader(singleTest.reqBody))
			if err != nil {
				t.Fatal(err)
			}

			simRequest.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handler)

			handler.ServeHTTP(rr, simRequest)

			if rr.Code != singleTest.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, singleTest.wantCode)
			}

			if rr.Body.String() != singleTest.want {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), singleTest.want)
			}
		})
	}
}
