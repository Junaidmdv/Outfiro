package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"gorm.io/gorm"
)

func RenderRazorpay(c *gin.Context) {
	c.HTML(http.StatusOK, "app.html", nil)
}

var RazorPayOrderID string

func RazapayPayment(c *gin.Context) {
	
	orderID := OID
	fmt.Println(OID)
	var order models.Order
	result := database.DB.Model(&models.Order{}).Preload("Payment").First(&order, orderID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	if order.Payment.PaymentStatus == models.PaymentPaid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment already completed"})
		return
	}

	AmountInPaisa := utils.RoundFloat(order.TotalAmount) * 100
	OrderIDstr := strconv.Itoa(int(orderID))

	data := map[string]interface{}{
		"amount":          AmountInPaisa,
		"currency":        "INR",
		"receipt":         OrderIDstr,
		"payment_capture": 1,
		"notes": map[string]interface{}{
			"description": fmt.Sprintf("Payment of order %s", OrderIDstr),
		},
	}

	cleint := razorpay.NewClient(os.Getenv("RAZORPAY_KEYID"), os.Getenv("RAZORPAY_SECRETE_KEY"))
	body, err := cleint.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	fmt.Println(cleint)

	Responce := models.RazorpayResponce{
		OrderID:  body["id"].(string),
		Amount:   AmountInPaisa,
		Currency: "INR",
		KeyID:    os.Getenv("RAZORPAY_KEYID"),
	}
	fmt.Println(Responce)
	RazorPayOrderID = body["receipt"].(string)
	c.JSON(http.StatusOK, Responce)

}

func VerifyPaymentHandler(c *gin.Context) {
	var req models.PaymentVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		PaymentFailed(c)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := req.RazorPayOrderID + "|" + req.RazorPayPaymentID
	h := hmac.New(sha256.New, []byte(os.Getenv("RAZORPAY_SECRETE_KEY")))
	h.Write([]byte(data))
	calculatedSignature := hex.EncodeToString(h.Sum(nil))

	// Verify signature
	if calculatedSignature != req.RazorpaySignature {
		PaymentFailed(c)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}
	OrderID, err := strconv.Atoi(RazorPayOrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed fetch razorpay id"})
		return
	}
	payment := models.Payment{
		RazorPayOrderID:   req.RazorPayOrderID,
		RazorpayPaymentID: req.RazorPayPaymentID,
		RazorPaySignature: req.RazorpaySignature,
		PaymentStatus:     models.PaymentPaid,
	}
	if err := database.DB.Model(&models.Payment{}).Where("order_id=?", OrderID).Updates(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed update the payment data"})
		return
	}
	if err := database.DB.Model(&models.Order{}).Where("id=?", OrderID).Update("order_status", models.Processing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order stauts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Payment verified successfully",
	})
}
func PaymentFailed(c *gin.Context) {
	OrderID, err := strconv.Atoi(RazorPayOrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data"})
	}
	if err := database.DB.Model(&models.Payment{}).Where("order_id=?", OrderID).Update("payment_status", models.PaymentFailed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the  payment details"})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"details": "Razorpay payment failed",
	})
}
