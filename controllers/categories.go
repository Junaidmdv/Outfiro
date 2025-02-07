package controllers

import (
	"fmt"
	"net/http"
	"outfiro/database"
	"outfiro/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCategories(c *gin.Context) {
	var categories []models.Categories
	result := database.DB.Find(&categories)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{"error": "categories not found"})
		} else {
			c.JSON(500, gin.H{"error": fmt.Sprintf("faild to fetch categories error %v", result.Error)})
		}
		return
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
		return
	}
	var verifyCategories models.Categories
	var count int64
	if err := database.DB.Model(&verifyCategories).Where("category_name=?", categories.CategoryName).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("internal server error %v", err)})
		return
	}
	if count > 0 {
		c.JSON(409, gin.H{"error": "category already exist "})
		return
	}
	fmt.Println(verifyCategories.CategoryName)
	fmt.Println(categories.CategoryName)
	if err := database.DB.Create(&categories).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to add new category",
			"errors": err.Error()})
		return
	}
	c.JSON(201, gin.H{
		"status":  "success",
		"message": "new category added",
		"data": gin.H{
			"category": categories,
		},
	})

}

func DeleteCategory(c *gin.Context) {
	category_id := c.Param("id")
	categoryId, err := strconv.Atoi(category_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid input formate"})
		return
	}
	fmt.Println(categoryId)
	var category models.Categories
	var count int64
	if err := database.DB.Model(&category).Where("id", categoryId).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed fetch category"})
		return
	}
	fmt.Println(count)
	if count < 0 {
		c.JSON(400, gin.H{"error": "category is not exist"})
		return
	}

	if err := database.DB.Delete(&category, categoryId).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "error to fetch the category",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "product category deleted succesfully",
	})
}

func EditCategory(c *gin.Context) {
	categoryIdstr := c.Param("id")
	fmt.Println(categoryIdstr)
	categoryId, err := strconv.Atoi(categoryIdstr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid code"})
		return
	}
	var EditCategory models.UpdateCategoryRequest
	var category models.Categories
	var count int64
	if err := database.DB.Model(&category).Where("id", categoryId).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed fetch category"})
		return
	}
	fmt.Println(count)
	if count == 0 {
		c.JSON(400, gin.H{"error": "category is not exist"})
		return
	}

	if err := c.BindJSON(&EditCategory); err != nil {
		c.JSON(400, gin.H{"error": "invalid input formate"})
		return
	}
	if err := database.DB.First(&category, categoryId).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch category"})
		return
	}
	if EditCategory.CategoriesName != nil {
		category.CategoryName = *EditCategory.CategoriesName
	}
	if EditCategory.Description != nil {
		category.Description = *EditCategory.Description
	}
	var exist bool
	if err := database.DB.Model(&models.Categories{}).
		Where("category_name = ?", category.CategoryName).
		Select("count(*) > 0").
		Scan(&exist).Error; err != nil {
		c.JSON(500, gin.H{"error": "database error occurred"})
		return
	}

	if exist {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"code":    "StatusConflict(409)",
			"details": "category already exists",
		})
		return
	}

	if err := database.DB.Where("id=?", categoryId).Updates(&category).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to update the resourse"})
		return
	}
	fmt.Println(category)

	c.JSON(200, gin.H{
		"status":   "success",
		"message":  "categories details updated",
		"category": category,
	})

}
