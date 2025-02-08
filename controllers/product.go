package controllers

import (
	"fmt"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"outfiro/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func GetProducts(c *gin.Context) {
	var products models.Products
	var productsRes []models.ProductResponce
	if err := database.DB.Model(&products).Find(&productsRes).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "fail to fetch the product",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "product succesfully fetched",
		"data": gin.H{
			"product": productsRes,
		},
	})
}
func GetProduct(c *gin.Context) {
	fmt.Println("Product Page")
	product_id := c.Param("id")
	id, err := strconv.Atoi(product_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"code":   "StatusBadRequest(400)",
			"error":  "Invalid product id.Provide valid product id",
		})
		return
	}
	var product models.Products
	var productRes models.ProductResponce
	result := database.DB.Model(product).Where("id=?",id).First(&productRes)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status": "error",
				"error":  "product not found",
			})
			return
		} else {
			c.JSON(500, gin.H{
				"error": "failed to fetch the product",
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"product":productRes,
	})
}

func DeleteProduct(c *gin.Context) {
	ProductId := c.Param("id")
	id, err := strconv.Atoi(ProductId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid product id",
		})
		return
	}
	var product models.Products
	result := database.DB.First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status": "error",
				"error":  "product not found",
			})
		} else {
			c.JSON(500, gin.H{
				"error": "failed to found the product",
			})
		}
		return
	}

	if err := database.DB.Delete(&product, id).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"error":  "Internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "product deleted succesfully",
	})

}

func EditProduct(c *gin.Context) {
	product_id := c.Param("id")
	id, err := strconv.Atoi(product_id)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"code":    "StatusBadRequest(400)",
			"message": "missing product id",
		})
		return
	}
	var product models.Products
	result := database.DB.First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"code":    "StatusNotFound(404)",
				"message": "Product is not availible",
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "error",
				"code":    "StatusInternalServerError",
				"message": "Database error",
			})
		}
		return
	}

	var update models.UpadatProduct
	if err := c.BindJSON(&update); err != nil {
		c.JSON(400, gin.H{
			"status": "error",
			"error":  "invalid inpute formate",
		})
	}
	validate := validator.New()
	validate.RegisterValidation("alpha_space", utils.ValidateAlphaNumSpace)

	if err := validate.Struct(&update); err != nil {
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		c.JSON(400, gin.H{
			"status":  "error",
			"code":    "StatusBadRequest",
			"message": errors,
		})
		return
	}
	fmt.Println(update)

	updateProduct := make(map[string]interface{})

	if update.ProductName != "" {
		updateProduct["product_name"] = update.ProductName
	}
	if update.Price != 0 {
		updateProduct["price"] = update.Price

	}
	if update.Quantity != 0 {
		updateProduct["quantity"] = update.Quantity
	}
	if update.Status != "" {
		updateProduct["status"] = update.Status
	}
	fmt.Println(len(updateProduct))
	fmt.Println(updateProduct)
	if len(updateProduct) == 0 {
		c.JSON(400, gin.H{
			"status":  "error",
			"code":    "StatusBadRequest(400)",
			"message": "No valid fields to update",
		})
		return
	}

	if err := database.DB.Model(&product).Updates(updateProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "StatusInternalServerError(500)",
			"message": "Failed to update product",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"code":    "statusOk(200)",
		"message": "Product details updated",
		"data":    product,
	})
}

func AddProduct(c *gin.Context) {
	var req models.ProductRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Please provide valid input",
			"error":   err.Error(),
		})
		return
	}
	validate := validator.New()
	validate.RegisterValidation("alpha_space", utils.ValidateAlphaNumSpace)
	if err := validate.Struct(&req); err != nil {
		fmt.Println(err.Error())
		errors := utils.UserFormateError(err.(validator.ValidationErrors))
		fmt.Println(errors)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "http.StatusBadRequest(400)",
			"message": errors,
		})
		return
	}

	var NewProduct models.Products
	var category models.Categories
	var count int64
	if err := database.DB.Model(NewProduct).Where("product_name=? and size=?", req.ProductName, req.Size).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "error in the database",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{

			"status":  "error",
			"code":    "StatusConflict(409)",
			"message": "product already exist",
		})
		return
	}

	fmt.Println(req.CategoriesName)
	if err := database.DB.Where("category_name=?", req.CategoriesName).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{"error": "category not availible"})
		} else {
			c.JSON(500, gin.H{"error": "fail to fetch the category details of the product"})
		}
		return
	}
	NewProduct = models.Products{
		ProductName: req.ProductName,
		Description: req.Description,
		CategoryId:  category.ID,
		Price:       req.Price,
		Status:      req.Status,
		Size:        req.Size,
		Quntity:     req.Quntity,
		ImageUrl:    req.ImageUrl,
	}
	if err := database.DB.Create(&NewProduct).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "new product added",
		"data":    NewProduct,
	})

}

func SearchProduct(c *gin.Context) {
	query := c.Query("product")
	if query == "" {
		c.JSON(400, gin.H{"error": "queri paramters required"})
		return
	}
	fmt.Println(query)
	var products []models.Products
	result := database.DB.Where("products.product_name LIKE ? OR products.description LIKE ? OR categories.category_name LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Joins("JOIN categories ON products.category_id = categories.id").
		Find(&products)
	fmt.Println(result)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"status":  "error",
				"code":    "StatusNotFound",
				"message": "product not found",
			})
			return
		} else {
			c.JSON(500, gin.H{"error": "Error searching for products"})
			return
		}

	}
	if len(products) == 0 {
		c.JSON(404, gin.H{
			"status":  "error",
			"code":    "StatusNotFound(404)",
			"message": "Product is found",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "fetched all product related to the queri",
		"data": gin.H{
			"products": products,
		},
	})
}
