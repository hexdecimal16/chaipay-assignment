package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hexdecimal16/chaipay-assignment/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/chaipay/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "Welcome to the Chaipay Assignment Server",
		})
	})
	r.GET("/chaipay/api/v1/get_intents", func(c *gin.Context) {
		var paymentIntents []models.PaymentIntent
		paymentIntents = append(paymentIntents, models.PaymentIntent{
			ID:        "1",
			Amount:    100,
			CreatedAt: 1579098983,
			Captured:  false,
			Refunded:  false,
			RefundId:  "",
			ChargeId:  "",
		})

		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": paymentIntents,
		})
	})
	r.POST("/chaipay/api/v1/create_intent", func(c *gin.Context) {
		c.JSON(http.StatusCreated, models.GenericIDResponse{
			ID: "pi_sdadtarad",
		})
	})
	r.POST("/chaipay/api/v1/capture_intent", func(c *gin.Context) {
		var paymentReq models.PaymentIntentRequest
		c.BindJSON(&paymentReq)
		log.Println(paymentReq)
		if paymentReq.ID == "" {
			c.JSON(http.StatusBadRequest, models.GenericIDResponse{
				Error: "ID is required",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "Payment Captured Successfully",
		})
	})

	r.POST("/chaipay/api/v1/create_refund", func(c *gin.Context) {
		var paymentReq models.PaymentIntentRequest
		c.BindJSON(&paymentReq)
		log.Println(paymentReq)
		if paymentReq.ID == "" {
			c.JSON(http.StatusBadRequest, models.GenericIDResponse{
				Error: "ID is required",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "Payment Refunded Started with ID: re_sdadtarad",
		})
	})
	return r
}

func TestHealth(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/chaipay/api/v1/health", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"Welcome to the Chaipay Assignment Server\",\"status\":200}", w.Body.String())

}

func TestGetPaymentIntents(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/chaipay/api/v1/get_intents", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":[{\"id\":\"1\",\"amount\":100,\"created_at\":1579098983,\"captured\":false,\"refunded\":false,\"refund_id\":\"\",\"charge_id\":\"\"}],\"status\":200}", w.Body.String())
}

func TestCreateIntent(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/chaipay/api/v1/create_intent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "{\"id\":\"pi_sdadtarad\",\"error\":\"\"}", w.Body.String())
}

func TestCaptureIntent(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	data, b := map[string]string{
		"id": "pi_sdadtarad",
	}, new(bytes.Buffer)

	json.NewEncoder(b).Encode(data)
	// tx := json.NewEncoder(w).Encode(data)
	req, _ := http.NewRequest(http.MethodPost, "/chaipay/api/v1/capture_intent", b)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"Payment Captured Successfully\",\"status\":200}", w.Body.String())
}

func TestCreatRefund(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	data, b := map[string]string{
		"id": "pi_sdadtarad",
	}, new(bytes.Buffer)

	json.NewEncoder(b).Encode(data)
	// tx := json.NewEncoder(w).Encode(data)
	req, _ := http.NewRequest(http.MethodPost, "/chaipay/api/v1/create_refund", b)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"Payment Refunded Started with ID: re_sdadtarad\",\"status\":200}", w.Body.String())

}
