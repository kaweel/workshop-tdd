//go:build unit_test
// +build unit_test

package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kaweel/workshop-tdd/payment/service"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	Calls []service.RequestPayment
	err   error
}

func (m *mockService) SetPayment(r service.RequestPayment, err error) {
	m.Calls = append(m.Calls, r)
	m.err = err
}

func (m *mockService) Payment(r service.RequestPayment) error {
	return m.err
}

func TestPaymentHandler(t *testing.T) {

	var (
		m  *mockService
		h  PaymentHandler
		rr *httptest.ResponseRecorder
		r  *mux.Router
	)

	setup := func() {
		m = &mockService{}
		h = NewHandler(m)
		rr = httptest.NewRecorder()
		r = mux.NewRouter()
		r.HandleFunc("/payment", h.Payment())
	}

	t.Run("invalid request should return bad request", func(t *testing.T) {
		setup()
		reqBody := bytes.NewBufferString(``)
		req, err := http.NewRequest(http.MethodPost, "/payment", reqBody)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, 0, len(m.Calls))
	})

	t.Run("failed payment should return unprocessable entity when", func(t *testing.T) {
		data := []struct {
			r   string
			err error
		}{
			{`{"orderID":1,"channel":"zebit","amount":100}`, errors.New("invalid payment channel")},
			{`{"orderID":2,"channel":"debit","amount":100}`, errors.New("order status is not request payment")},
			{`{"orderID":3,"channel":"debit","amount":100}`, errors.New("customer status is not active")},
			{`{"orderID":4,"channel":"debit","amount":100}`, errors.New("customer amount is not enough")},
			{`{"orderID":5,"channel":"debit","amount":100}`, errors.New("merchant status is not active")},
		}

		for _, v := range data {
			t.Run(fmt.Sprintf("%v", v.err.Error()), func(t *testing.T) {
				setup()
				expected := service.RequestPayment{}
				json.Unmarshal([]byte(v.r), &expected)
				m.SetPayment(expected, v.err)
				reqBody := bytes.NewBufferString(v.r)
				req, err := http.NewRequest(http.MethodPost, "/payment", reqBody)
				if err != nil {
					t.Fatal(err)
				}

				r.ServeHTTP(rr, req)

				assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
				assert.Equal(t, 1, len(m.Calls))
				assert.Equal(t, expected, m.Calls[0])
			})
		}

	})

	t.Run("failed payment should return internal error when unknown error occurred", func(t *testing.T) {
		setup()
		reqStr := `{"orderID":1,"channel":"zebit","amount":100}`
		expected := service.RequestPayment{}
		json.Unmarshal([]byte(reqStr), &expected)
		m.SetPayment(expected, errors.New("Unknow error"))
		reqBody := bytes.NewBufferString(reqStr)
		req, err := http.NewRequest(http.MethodPost, "/payment", reqBody)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, 1, len(m.Calls))
		assert.Equal(t, expected, m.Calls[0])
	})

}
