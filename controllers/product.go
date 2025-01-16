package controllers

import (
	
	"outfiro/database"
	"outfiro/models"
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)




func GetProducts(c *gin.Context){
var products []models.Products
  if err:=database.DB.Find(&products).Error; err!=nil{
	 c.JSON(500,gin.H{
		"error":"fail to fetch the product",
	 })
	 return
  }
  c.JSON(200,gin.H{
	"status":"success",
	 "message":"product succesfully fetched",
	 "data":gin.H{
		"product":products,
	 },
  })  	
}
// type Products struct {
// 	gorm.Model
// 	ProductName string  `gorm:"not null" json:"product name"`
// 	Description string  `gorm:"not null" json:"description"`
// 	CategoryId  uint    `gorm:"not null" json:"categoryid"`
// 	Price       float64 `gorm:"not null" json:"price"`
// 	Status      string  `gorm:"not null type:enum('in stock','out of stock');" json:"status"`
// 	Size        string  `gorm:"not null"`
// 	Quntity     string   `gorm:"not null" json:"quantity"`
// 	Discount    string  `json:"discount"`
// }
func GetProduct(c *gin.Context){
	product_id:=c.Param("id")
	id,err:=strconv.Atoi(product_id)
	if err !=nil{
		c.JSON(400,gin.H{
		  "error":"invalid product id",
		})
		return
	}

	var product models.Products 
    result:=database.DB.First(&product,id)
	 if result.Error !=nil{
		if result.Error==gorm.ErrRecordNotFound{
			 c.JSON(404,gin.H{
				"status":"error",
				"error":"product not found",
			 })
		}else{
			c.JSON(500,gin.H{
				 "error":"failed to fetch the product",
			})	
		}
		return
	 }
    c.JSON(200,gin.H{
		"status":"success",
		"product":product,

	})
}

func DeleteProduct(c *gin.Context){
//admin check
    ProductId:=("id")
	var product models.Products
    id,err:=strconv.Atoi(ProductId)
	if err !=nil{
		c.JSON(400,gin.H{
			"error":"invalid product id",
		})
		return
	}
    result:=database.DB.First(&product,id)
	 if result.Error !=nil{
		if result.Error==gorm.ErrRecordNotFound{
			 c.JSON(404,gin.H{
				"status":"error",
				"error":"product not found",
			 })
		}else{
			c.JSON(500,gin.H{
				 "error":"failed to found the product",
			})	
		}
		return
	}

    if err:=database.DB.Delete(&product,id).Error;err !=nil{
		c.JSON(500,gin.H{
			"status":"error",
			"error":"Internal server error",
		})
		return
	}
    c.JSON(200,gin.H{
		"status":"success",
		"message":"product deleted succesfully",
	})

}


func EditProduct(c *gin.Context){
	//otherise the admin 
   ProductId:=c.Param("id")
   id,err:=strconv.Atoi(ProductId)
   if err!=nil{
	c.JSON(400,gin.H{
		"status":"error",
		"message":"Invalid product_id",
	})
	return
   }
   var UpadatProduct models.Products 
    if err:=c.BindJSON(&UpadatProduct);err !=nil{
		 c.JSON(400,gin.H{
			"error":"Invalid formate",
		 })
		 return

	}
   var product models.Products
   result:=database.DB.First(&product,id) 
   if result.Error !=nil{
	  if result.Error==gorm.ErrRecordNotFound{
		 c.JSON(404,gin.H{"error":"product not found"})
	  }else{
		c.JSON(500,gin.H{"error":"Internal server error"})
	  }
	  return
   }
   product.ProductName=UpadatProduct.ProductName
   product.Status=UpadatProduct.Status
   product.Quntity=UpadatProduct.Quntity
   product.Discount=UpadatProduct.Discount
   
   result=database.DB.Save(&product)
   if result.Error !=nil{
      c.JSON(500,gin.H{
		"error":"Internal server error",
	  })
	  return
   }

    c.JSON(200,gin.H{
		"status":"success",
		"message":"product updated success",
		"Data":product,
	  })


}



func AddProduct(c *gin.Context){
	var req models.ProductRequest
	if err:=c.BindJSON(&req);err !=nil {
       c.JSON(400,gin.H{"error":"Invalid code"})
	   return
	}
    var category models.Categories
	if err:=database.DB.Where("category_name = ?", req.CategoriesName).First(&category).Error;err !=nil{
        if err ==gorm.ErrRecordNotFound{
			 c.JSON(400,gin.H{"error":"category not availible"})
		}else{
			c.JSON(500,gin.H{"error":"fail to fetch the category details of the product"})
		}
	   return
	}
	//crete product 
	NewProduct:=models.Products{
		ProductName: req.ProductName,
		Description: req.Description,
		CategoryId: category.ID,
		Price: req.Price,
		Status: req.Status,
        Size: req.Size,
		Quntity: req.Quntity,
		Discount: req.Discount,

	}
	if err:=database.DB.Create(&NewProduct).Error;err !=nil{
		c.JSON(500,gin.H{"error":"Internal server error"})
		return
	}

   c.JSON(200,gin.H{
	 "status":"success",
	 "message":"new product added",
	 "data":NewProduct,
   })
   

}

func SearchProduct(c *gin.Context){
	query:=c.Query("Query")
	if query==""{
		c.JSON(400,gin.H{"error":"queri paramters required"})
		return
	}
  var products []models.Products
  result := database.DB.Where("product_name LIKE ? OR description LIKE ? OR categories.category_name LIKE ?", 
         "%"+query+"%", "%"+query+"%", "%"+query+"%").
         Joins("JOIN categories ON products.CategoryID = categories.ID").
         Find(&products)

  if result.Error != nil {
     c.JSON(500, gin.H{"error": "Error searching for products"})
     return
 }

 c.JSON(200,gin.H{
	"status":"success",
	"message":"fetched all product related to the queri",
	 "data":gin.H{
		"products":products,
	 },
 })
}



