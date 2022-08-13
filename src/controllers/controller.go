//  Chaipay Api:
//   version: 1.0.0
//   title:Chai Pay Assignment Api
//  Schemes: http
//  Host: localhost:5000
//  BasePath: /api/v1
//  Produces:
//    - application/json
package controllers

import (
	"log"
	"net/http"

	"github.com/hexdecimal16/chaipay-assignment/config"
	"github.com/hexdecimal16/chaipay-assignment/database"
	"github.com/hexdecimal16/chaipay-assignment/src/models"

	"github.com/gin-gonic/gin"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/refund"
)

// Get the stripe key from the config file.
var key = config.Config("STRIPE_KEY")

// Health get the health of the server.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome to the Chaipay Assignment Server",
	})
	return
}

// GetPaymentIntents gets all the payment intents in the database.
func GetPaymentIntents(c *gin.Context) {
	var paymentIntents []models.PaymentIntent

	// Get all the payment intents from the database
	result := database.DB.Model(&models.PaymentIntent{}).Find(&paymentIntents)

	// If there is an error while fetching the payment intents from the database
	if result.Error != nil {
		log.Printf("Error while fetching payment intents\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error while fetching payment intents",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   paymentIntents,
	})
	return

}

// CreatePaymentIntent creates a payment intent and adds it to the database.
// It also creates a customer and attaches it to the payment intent.
// It takes amount as an argument and returns the payment intent id.
func CreateIntent(c *gin.Context) {
	var paymentReq models.PaymentIntent
	c.BindJSON(&paymentReq)
	amount := paymentReq.Amount

	stripe.Key = key

	params := &stripe.CustomerParams{
		Description:      stripe.String("Stripe Developer"),
		Email:            stripe.String("gostripe@stripe.com"),
		PreferredLocales: stripe.StringSlice([]string{"en", "es"}),
	}
	// Create a new Customer
	customer, err := customer.New(params)
	if err != nil {
		log.Printf("Error while creating a new user\n")
		code := stripeErrorLogger(err)
		resp := models.GenericIDResponse{Error: err.Error()}
		c.JSON(code, resp)
		return
	}

	paymentIntentParams := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(amount),
		Currency:      stripe.String(string(stripe.CurrencyINR)),
		Customer:      stripe.String(customer.ID),
		Description:   stripe.String("test payment intent"),
		PaymentMethod: stripe.String("pm_card_amex_threeDSecureNotSupported"),
		Confirm:       stripe.Bool(true),
	}

	// Create a Payment Intent
	resp, err := paymentintent.New(paymentIntentParams)
	if err != nil {
		log.Printf("Error while creating a new payment\n")
		code := stripeErrorLogger(err)
		resp := models.GenericIDResponse{Error: err.Error()}
		c.JSON(code, resp)
		return
	}

	newPaymentIntent := models.PaymentIntent{
		ID:        resp.ID,
		Amount:    resp.Amount,
		CreatedAt: resp.Created,
		Refunded:  false,
		Captured:  false,
	}

	// Add the payment intent to the database
	result := database.DB.Create(&newPaymentIntent)

	if result.Error != nil {
		log.Printf("Error while inserting new payment intent into db but was created with ID,%s\n", resp.ID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error while inserting new payment intent into db but was created with ID " + resp.ID,
		})
		return
	}

	res := models.GenericIDResponse{
		ID: resp.ID,
	}

	c.JSON(http.StatusCreated, res)
	return

}

// CapturePaymentIntent captures the payment intent and updates the database.
// It takes payment intent id as an argument and returns a success message.
func CaptureIntent(c *gin.Context) {
	stripe.Key = key
	var id = c.Param("id")
	if id == "" {
		log.Printf("Error while capturing payment intent id\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Error while capturing payment intent",
		})
		return
	}
	var payment = &models.PaymentIntentRequest{
		ID: id,
	}

	pi, err := paymentintent.Get(
		payment.ID,
		nil,
	)
	if err != nil {
		log.Printf("Error while fetching payment details\n")
		code := stripeErrorLogger(err)
		resp := models.GenericIDResponse{Error: err.Error()}
		c.JSON(code, resp)
		return
	}

	// if the payment intent is not captured
	if pi.Status != "succeeded" {
		pi, err = paymentintent.Capture(
			payment.ID,
			&stripe.PaymentIntentCaptureParams{},
		)

		if err != nil {
			log.Printf("Error while capturing a payment\n")
			code := stripeErrorLogger(err)
			resp := models.GenericIDResponse{Error: err.Error()}
			c.JSON(code, resp)
			return
		}

	}

	// Update the database with the captured payment intent
	database.DB.Model(&models.PaymentIntent{}).Where("id = ?", payment.ID).Update(
		models.PaymentIntent{
			Captured: true,
			ChargeId: pi.Charges.Data[0].ID,
		},
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Payment Captured Successfully",
	})

	return

}

// CreateRefund creates a refund for the payment intent and updates the database.
// It takes payment intent id as an argument and returns a success
// message along with the refund id.
func CreateRefund(c *gin.Context) {
	stripe.Key = key
	var id = c.Param("id")
	if id == "" {
		log.Printf("Error while capturing payment intent id\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Error while capturing payment intent",
		})
		return
	}
	var payment = &models.PaymentIntentRequest{
		ID: id,
	}
	var pi models.PaymentIntent
	// Get the payment intent details from stripe
	result := database.DB.Model(&models.PaymentIntent{}).Where("id = ?", payment.ID).Find(&pi)

	if result.Error != nil {
		log.Printf("Error while fetching payment intent details\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error while fetching payment intent details",
		})
		return
	}

	refundParams := stripe.RefundParams{
		Charge: stripe.String(pi.ChargeId),
	}

	// Refund the payment charge
	refund, err := refund.New(&refundParams)
	if err != nil {
		log.Printf("Error while refunding the payment\n")
		code := stripeErrorLogger(err)
		resp := gin.H{"error": err.Error()}
		c.JSON(code, resp)
		return
	}

	// Update the database with the payment id
	result = database.DB.Model(&models.PaymentIntent{}).Where("id = ?", payment.ID).Update(
		models.PaymentIntent{Refunded: true, RefundId: refund.ID})

	if result.Error != nil {
		log.Printf("Error while updating the refunded status\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error while updating the refunded status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Payment Refund Started with ID " + refund.ID,
	})
}

// Error404 handles 404 errors.
func Error404(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  404,
		"message": "Route Not Found",
	})
	return
}

// stripeErrorLogger is a generic function to log the stripe errors.
func stripeErrorLogger(err error) int {

	if stripeErr, ok := err.(*stripe.Error); ok {
		// The Code field will contain a basic identifier for the failure.
		switch stripeErr.Code {
		case stripe.ErrorCodeCardDeclined:
		case stripe.ErrorCodeExpiredCard:
		case stripe.ErrorCodeIncorrectCVC:
		case stripe.ErrorCodeIncorrectZip:
			// etc.
		}

		// The Err field can be coerced to a more specific error type with a type
		// assertion. This technique can be used to get more specialized
		// information for certain errors.
		if cardErr, ok := stripeErr.Err.(*stripe.CardError); ok {
			log.Printf("Card was declined with code: %v\n", cardErr.DeclineCode)
		} else {
			log.Printf("Other Stripe error occurred: %v\n", stripeErr.Error())
		}
		return stripeErr.HTTPStatusCode
	} else {
		log.Printf("Other error occurred: %v\n", err.Error())
		return 500
	}
}
