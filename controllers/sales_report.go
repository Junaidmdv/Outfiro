package controllers

import (
	"fmt"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func SalesReport(c *gin.Context) {
	var request models.SalesReportRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	}

	var err error
	var StartDate, EndDate time.Time

	layout := "2006-01-02 15:04:05"

	switch request.Limit {
	case "day":
		today := time.Now()
		StartDate = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
		EndDate = time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 999999999, today.Location())

	case "week":
		StartDate = time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
		EndDate = time.Now()

	case "month":
		StartDate = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
		EndDate = time.Now()

	case "year":
		StartDate = time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Now().Location())
		EndDate = time.Now()

	case "custom":
		StartDate, err = time.Parse(layout, request.StartDate)
		if err != nil {
			fmt.Println("Invalid start date:", err)
			c.JSON(400, gin.H{"error": "Invalid start date format"})
			return
		}
		EndDate, err = time.Parse(layout, request.EndDate)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid end date format"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid method"})
		return

	}
	var SalesReport []models.SalesReport
	if err := database.DB.Model(&models.Order{}).
		Select("orders.user_id, orders.id, users.email,orders.total_discount, orders.created_at, orders.product_quantity, orders.coupon_offer, orders.total_amount, payments.payment_method, payments.payment_status, orders.order_status").
		Joins("JOIN payments ON orders.id = payments.order_id").
		Joins("JOIN users ON orders.user_id = users.id").
		Where("orders.created_at BETWEEN ? AND ?", StartDate, EndDate).
		Find(&SalesReport).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var TotalOrderAmount float64
	var TotalDiscount float64
	var TotalCouponOffer float64
	var TotalProductSold uint

	for _, item := range SalesReport {
		TotalOrderAmount += item.TotalAmount
		TotalDiscount += item.Discount
		TotalCouponOffer += item.CouponOffer
		TotalProductSold += item.ProductQuantity
	}
	c.JSON(200, gin.H{
		"status": "success",
		"Data": gin.H{
			"Total sales":        TotalOrderAmount,
			"Total coupon offer": TotalCouponOffer,
			"Total discount":     TotalDiscount,
			"Total product sold": TotalProductSold,
			"Orders":             SalesReport,
		},
	})
}

func SalesReportPDF(c *gin.Context) {
	var request models.SalesReportRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	}

	var err error
	var StartDate, EndDate time.Time
	SalesReportHead := ""

	layout := "2006-01-02 15:04:05"

	switch request.Limit {
	case "day":
		today := time.Now()
		StartDate = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
		EndDate = time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 999999999, today.Location())
		SalesReportHead = "Daily"
	case "week":
		StartDate = time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
		EndDate = time.Now()
		SalesReportHead = "Weekly"
	case "month":
		StartDate = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
		EndDate = time.Now()
		SalesReportHead = "Montly"
	case "year":
		StartDate = time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Now().Location())
		EndDate = time.Now()
		SalesReportHead = "Annual"
	case "custom":
		StartDate, err = time.Parse(layout, request.StartDate)
		if err != nil {
			
			c.JSON(400, gin.H{"error": "Invalid start date format"})
			return
		}
		EndDate, err = time.Parse(layout, request.EndDate)

		if err != nil {
			
			c.JSON(400, gin.H{"error": "Invalid end date format"})
			return
		}
	    lyout:="02 Jan 06"
		st := StartDate.Format(lyout)
		end := EndDate.Format(lyout)
		SalesReportHead = fmt.Sprintf("%v-%v ", st, end)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filtering method"})
		return
	}
	var SalesReport []models.SalesReport
	if err := database.DB.Model(&models.Order{}).
		Select("orders.user_id, orders.id, users.email,orders.total_discount, orders.created_at, orders.product_quantity, orders.coupon_offer, orders.total_amount, payments.payment_method, payments.payment_status, orders.order_status").
		Joins("JOIN payments ON orders.id = payments.order_id").
		Joins("JOIN users ON orders.user_id = users.id").
		Where("orders.created_at BETWEEN ? AND ?", StartDate, EndDate).
		Find(&SalesReport).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	var TotalOrderAmount float64
	var TotalDiscount float64
	var TotalCouponOffer float64
	var TotalProductSold uint

	for _, item := range SalesReport {
		TotalOrderAmount += item.TotalAmount
		TotalDiscount += item.Discount
		TotalCouponOffer += item.CouponOffer
		TotalProductSold += item.ProductQuantity

	}
	// Generate PDF
	m, err := CreateSalesReport(SalesReportHead, SalesReport, TotalProductSold, TotalDiscount, TotalCouponOffer, TotalOrderAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfdocs, err := m.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed generate the pdf"})
		return
	}

	pdfbytes := pdfdocs.GetBytes()

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=document.pdf")
	c.Header("Content-Length", strconv.Itoa(len(pdfbytes)))

	// Write the PDF to the response
	c.Data(http.StatusOK, "application/pdf", pdfbytes)
}

func CreateSalesReport(salesReportHead string, SalesReport []models.SalesReport, TotalProductSold uint, TotalDiscount float64, TotalCouponOffer float64, TotalOrderAmount float64) (core.Maroto, error) {
	//moroto configration
	configrate := config.NewBuilder().
		WithOrientation(orientation.Vertical).
		WithPageSize(pagesize.A4).
		WithRightMargin(10).
		WithLeftMargin(10).
		WithBottomMargin(10).
		WithTopMargin(10).Build()

	m := maroto.New(configrate)

	err := m.RegisterHeader(getPageHeader(salesReportHead))
	if err != nil {
		return nil, fmt.Errorf("failed to create pdf page header")

	}
	err = m.RegisterFooter(getPageFooter())
	if err != nil {
		return nil, fmt.Errorf("failed to create pdf footer")
	}
	m.AddRows(text.NewRow(12, "Orders", props.Text{
		Top:   4,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	var salesData []models.SalesReportPDf

	for _, val := range SalesReport {
		salesData = append(salesData, models.SalesReportPDf{
			OrderDate:     val.OrderDate.Format("2006-01-02 15:04:05"),
			Email:         val.Email,
			Quantity:      fmt.Sprintf("%d", val.ProductQuantity),
			Discount:      fmt.Sprintf("%.2f", val.Discount),
			TotalAmount:   fmt.Sprintf("%.2f", val.TotalAmount),
			PaymentMethod: val.PaymentMethod,
		})
	}

	m.AddRows(getContent(salesData, TotalProductSold, TotalDiscount, TotalCouponOffer, TotalOrderAmount)...)

	return m, nil

}

func getPageHeader(SalesReportHead string) core.Row {
	return row.New(20).Add(
		image.NewFromFileCol(4, "asset/logo.png", props.Rect{
			// Center: true,
			Percent: 100,
		}),
		col.New(10).Add(
			text.New(fmt.Sprintf("%s Sales Report", SalesReportHead), props.Text{
				Top:   6,
				Right: 4,
				Align: align.Left,
				Style: fontstyle.Bold,
				Size:  13,
			}),
		),
	)
}

func getPageFooter() core.Row {
	now := time.Now()
	return row.New(30).Add(
		col.New(12).Add(
			text.New("Report Generated On: "+now.Format("January 2, 2006"), props.Text{
				Style: fontstyle.BoldItalic,
				Size:  8,
			}),
			text.New("Contact Us: outfiro@gmail.com", props.Text{
				Style: fontstyle.Bold,
				Size:  8,
				Top:   5,
			}),
		),
	)
}
func getContent(salesData []models.SalesReportPDf, TotalProductSold uint, offerDiscount float64, coupon float64, TotalAmount float64) []core.Row {
	rows := []core.Row{
		row.New(10).Add(
			// col.New(5),
			text.NewCol(3, "Date", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(2, "User Email", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(1, "Quantity", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(1, "Discount", props.Text{Size: 8, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(2, "Total", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(3, "Payment Method", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
		),
	}
	var contentRow []core.Row

	for i, data := range salesData {
		r := row.New(4).Add(
			// col.New(5),
			text.NewCol(3, data.OrderDate, props.Text{Size: 8, Align: align.Center}),
			text.NewCol(2, data.Email, props.Text{Size: 8, Align: align.Center}),
			text.NewCol(1, data.Quantity, props.Text{Size: 8, Align: align.Center}),
			text.NewCol(1, data.Discount, props.Text{Size: 7, Align: align.Center}),
			text.NewCol(2, data.TotalAmount, props.Text{Size: 8, Align: align.Center}),
			text.NewCol(3, data.PaymentMethod, props.Text{Size: 8, Align: align.Center}),
		)
		if i%2 == 0 {
			skyBlue := getSkyBlueColor()
			r.WithStyle(&props.Cell{BackgroundColor: skyBlue})
		}
		contentRow = append(contentRow, r)
	}
	rows = append(rows, contentRow...)

	rows = append(rows, row.New(5).Add(
		text.NewCol(4, fmt.Sprintf("Total Product sold:%d", TotalProductSold), props.Text{
			Top:   3,
			Style: fontstyle.BoldItalic,
			Size:  8,
			Align: align.Left,
		}),
	),
	)
	rows = append(rows, row.New(5).Add(
		text.NewCol(4, fmt.Sprintf("Total Discount offer:%.2f", offerDiscount), props.Text{
			Top:   2,
			Style: fontstyle.BoldItalic,
			Size:  8,
			Align: align.Left,
		}),
	),
	)
	rows = append(rows, row.New(5).Add(
		text.NewCol(4, fmt.Sprintf("Coupon offer:%.2f", coupon), props.Text{
			Top:   2,
			Style: fontstyle.BoldItalic,
			Size:  8,
			Align: align.Left,
		}),
	),
	)
	rows = append(rows, row.New(5).Add(
		text.NewCol(4, fmt.Sprintf("Total Amount:%.2f", TotalAmount), props.Text{
			Top:   2,
			Style: fontstyle.BoldItalic,
			Size:  8,
			Align: align.Left,
		}),
	),
	)

	return rows
}

func getSkyBlueColor() *props.Color {
	return &props.Color{
		//r: 135, g: 206, b: 235
		//(160,217,239)
		Red:   160,
		Green: 217,
		Blue:  239,
	}
}
