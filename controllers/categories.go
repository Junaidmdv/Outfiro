package controllers

import (
	"fmt"
	"outfiro/database"
	"outfiro/models"
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCategories(c *gin.Context) {
	var categories models.Categories
	result := database.DB.Find(&categories)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{"error": "categories not found"})
		} else {
			c.JSON(500, gin.H{"error": fmt.Sprintf("faild to fetch categories error %v", result.Error)})
		}
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "categories fetched succesfully",
		"data": gin.H{
			"categories": categories,
		},
	})
}

func AddCategories(c *gin.Context) {
	var categories models.Categories
	if err := c.BindJSON(&categories); err != nil {
		c.JSON(400, gin.H{"error": "Invalide code"})
	}
	var verifyCategories models.Categories
	var count int64
	if err := database.DB.Model(&verifyCategories).Where("verify_categories=?", categories.CategoryName).Count(&count); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("internal server error %v", err)})
	}
	if count < 0 {
		c.JSON(409, gin.H{"error": "category already exist "})
	}

	if err := database.DB.Create(&categories); err != nil {
		c.JSON(500, gin.H{"error": "failed to add new category"})
	}

	c.JSON(201, gin.H{
		"status":  "success",
		"message": "new category added",
	})

}

func DeleteCategory(c *gin.Context) {
	category_id := c.Param("id")
	categoryId, err := strconv.Atoi(category_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid input formate"})
	}
	var category models.Categories
	result := database.DB.Model(&category).First("id", categoryId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"error": "category is not found",
			})
		} else {
			c.JSON(500, gin.H{"error": "failed fetch the category"})
		}
		return
	}
	if err := database.DB.Delete(&category, categoryId).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "error to fetch the category",
		})
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "product deleted succesfully",
	})
}


func EditCategory(c *gin.Context) {
    categoryIdstr :=c.Param("id")
	categoryId,err:=strconv.Atoi(categoryIdstr)
	if err !=nil{
		c.JSON(200,gin.H{"error":"Invalid code"})
	}
	var category models.Categories
    result:=database.DB.First(&category,categoryId)
	 if result.Error !=nil{
		 if result.Error==gorm.ErrRecordNotFound{
			 c.JSON(404,gin.H{"error":fmt.Sprintf("%v",result.Error)})
		 }else{
			 c.JSON(500,gin.H{"error":ErrInternalServer})
		 }
	 }
	var EditCategory models.UpdateCategoryRequest 
     if err:=c.BindJSON(&EditCategory);err !=nil{
		c.JSON(400,gin.H{"error":"invalid input formate"})
	 }
	if EditCategory.CategoriesName !=nil{
        category.CategoryName=*EditCategory.CategoriesName
	}
	if EditCategory.Description !=nil{
		category.CategoryName=*EditCategory.Description
	}
    
	if err :=database.DB.Save(&category).Error;err !=nil{
		c.JSON(500,gin.H{"error":"failed to update the resourse"})
	}

	c.JSON(200,gin.H{
		"status":"success",
		 "message":"categories details updated",
		 "category":category,
	})

}
