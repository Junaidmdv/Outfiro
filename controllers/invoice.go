package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"gorm.io/gorm"
)

func GenerateOrderInvoice(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.Users
	if err := database.DB.Model(&user).Where("id=?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
	}
	orderId := c.Param("id")
	orderID, err := strconv.Atoi(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid param order id"})
		return
	}
	var order models.Order
	result := database.DB.Model(&models.Order{}).Where("id=?", orderID).First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order details is not found",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "databaser error"})
			return
		}
	}

	var OrderItem []models.OrderItem
	if err := database.DB.Model(&models.OrderItem{}).Where("order_id=? AND order_item_status<>?", orderID, models.Cancelled).Preload("Product").Find(&OrderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	m, err := OrderInvoice(OrderItem, order, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate order invoice"})
		return
	}
	invoicedocs, err := m.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed the create the pdf"})
		return
	}
	invoicebytes := invoicedocs.GetBytes()
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=OrderInvoice.pdf")
	c.Header("Content-Length", strconv.Itoa(len(invoicebytes)))

	c.Data(http.StatusOK, "application/pdf", invoicebytes)

}

func OrderInvoice(Orderitems []models.OrderItem, order models.Order, userData models.Users) (core.Maroto, error) {

	config := config.NewBuilder().
		WithOrientation(orientation.Vertical).
		WithPageSize(pagesize.A4).
		WithRightMargin(15).
		WithLeftMargin(15).
		WithBottomMargin(15).
		WithTopMargin(15).
		Build()

	m := maroto.New(config)

	err := m.RegisterHeader(getInvoiceHeader())
	if err != nil {
		return nil, fmt.Errorf("failed to create header")
	}

	m.AddRow(15,
		col.New(6).Add(
			text.New("Billing to:", props.Text{
				Align:  align.Left,
				Style:  fontstyle.Bold,
				Size:   10,
				Top:    0,
				Bottom: 0, // Reduced bottom spacing
			}),
			text.New(fmt.Sprintf("  %s %s", userData.FirstName, userData.LastName), props.Text{
				Align:  align.Left,
				Style:  fontstyle.Normal,
				Size:   9,
				Top:    5,
				Bottom: 0.5,
			}),
			text.New(fmt.Sprintf("  %s", userData.Email), props.Text{
				Align:  align.Left,
				Style:  fontstyle.Normal,
				Size:   9,
				Top:    10,
				Bottom: 0.5,
			}),
			text.New(fmt.Sprintf("  %s", userData.PhoneNumber), props.Text{
				Align:  align.Left,
				Style:  fontstyle.Normal,
				Size:   9,
				Top:    15,
				Bottom: 0.5,
			}),
		),
		col.New(6).Add(
			text.New(fmt.Sprintf("Date:%v", time.Now().Format("2006-01-02 15:04:05")), props.Text{
				Align:  align.Right,
				Style:  fontstyle.Bold,
				Size:   9,
				Top:    0.5,
				Bottom: 0.5,
			}),
			text.New(fmt.Sprintf("INVOICE:ORD#%d", order.ID), props.Text{
				Align:  align.Right,
				Style:  fontstyle.Bold,
				Size:   9,
				Top:    5,
				Bottom: 0.5,
			}),
		),
	)

	m.AddRow(20)

	var invoiceItems []models.OrderInvoice
	for _, items := range Orderitems {
		invoiceItems = append(invoiceItems, models.OrderInvoice{
			Item:         items.Product.ProductName,
			Descritption: items.Product.Description,
			Quantity:     fmt.Sprintf("%d", items.Quantity),
			Price:        fmt.Sprintf("%.2f", items.Product.Price),
			DiscountRate: fmt.Sprintf("%.2f", items.Product.Discount),
			Total:        fmt.Sprintf("%.2f", items.Product.Price*float64(items.Quantity)),
		})
	}

	m.AddRows(AddContent(invoiceItems)...)

	m.AddRow(20)

	lineCol := line.NewCol(12)
	m.AddRow(10, lineCol)
	m.AddRows(AddTotalRow(order.TotalDiscount, order.SubTotal, order.CouponOffer, order.DeliveryCharge, order.TotalAmount)...)

	return m, nil

}

func getInvoiceHeader() core.Row {
	return row.New(50).Add(
		image.NewFromFileCol(12, "asset/logo2.png", props.Rect{
			Center:  true,
			Percent: 100,
		}),
	)

}

func AddContent(InvoiceItems []models.OrderInvoice) []core.Row {
	rows := []core.Row{
		row.New(10).Add(
			text.NewCol(2, "Item", props.Text{Size: 10, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(3, "Description", props.Text{Size: 10, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(1, "Qty", props.Text{Size: 10, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(2, "Price", props.Text{Size: 10, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(2, "Discount", props.Text{Size: 10, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(2, "Amount", props.Text{Size: 10, Align: align.Center, Style: fontstyle.Bold}),
		),
	}

	var invoicedata []core.Row

	for i, items := range InvoiceItems {
		r := row.New(7).Add(
			text.NewCol(2, items.Item, props.Text{Size: 7, Align: align.Center, Top: 2}),
			text.NewCol(3, items.Descritption, props.Text{Size: 5, Align: align.Center, Top: 2}),
			text.NewCol(1, items.Quantity, props.Text{Size: 7, Align: align.Center, Top: 2}),
			text.NewCol(2, items.Price, props.Text{Size: 7, Align: align.Center, Top: 2}),
			text.NewCol(2, items.DiscountRate, props.Text{Size: 7, Align: align.Center, Top: 2}),
			text.NewCol(2, items.Total, props.Text{Size: 7, Align: align.Center, Top: 2}),
		)
		if i%2 == 0 {
			skyBlue := getSkyBlueColor()
			r.WithStyle(&props.Cell{BackgroundColor: skyBlue})

		}
		invoicedata = append(invoicedata, r)
	}
	rows = append(rows, invoicedata...)

	return rows

}

func AddTotalRow(Discount float64, SubTotal float64, CouponDiscount float64, DeliveryCharge float64, TotalAmount float64) []core.Row {
	var rows []core.Row
	rows = append(rows, row.New(30).Add(
		col.New(11).Add(
			text.New(fmt.Sprintf("Sub Total: %.2f", SubTotal), props.Text{
				Align:  align.Right,
				Style:  fontstyle.Bold,
				Size:   10,
				Top:    2,
				Bottom: 0.5,
			}),
			text.New(fmt.Sprintf("Product Discount: %.2f", Discount), props.Text{
				Align:  align.Right,
				Style:  fontstyle.Normal,
				Size:   10,
				Top:    10,
				Bottom: 0.5,
			}),
			text.New(fmt.Sprintf("Coupon offer: %.2f", CouponDiscount), props.Text{
				Align:  align.Right,
				Style:  fontstyle.Normal,
				Size:   10,
				Top:    16,
				Bottom: 0.5,
			}),
			text.New(fmt.Sprintf("Other Charges: %.2f", DeliveryCharge), props.Text{
				Align:  align.Right,
				Style:  fontstyle.Normal,
				Size:   10,
				Top:    22,
				Bottom: 0.5,
			}),

			text.New(fmt.Sprintf("Total amount: %.2f", TotalAmount), props.Text{
				Align:  align.Right,
				Style:  fontstyle.Bold,
				Size:   11,
				Top:    31,
				Bottom: 5,
			}),
		),
	),
	)
	lineCol := line.NewCol(12)
	rows = append(rows, row.New(10).Add(lineCol))

	rows = append(rows, row.New(30).Add(
		signature.NewCol(9, "Authorized Signatory", props.Signature{FontFamily: fontfamily.Courier}),
		code.NewQrCol(3, "https://outfiro.com", props.Rect{
			Percent: 75,
			Center:  true,
		}),
	))

	return rows
}
