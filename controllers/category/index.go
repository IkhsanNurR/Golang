package category_controller

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/database"
	category_model "main.go/models/public/category"
	category_response "main.go/response"
)

func GetAllCategory(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	pagesStr := c.DefaultQuery("pages", "1")
	sortKeyStr := c.DefaultQuery("sort_key", "id")
	sortByStr := c.DefaultQuery("sort_by", "asc")
	CategoryStr := c.DefaultQuery("category_name", "")
	status := c.DefaultQuery("status", "")

	if limitStr == "" || pagesStr == "" {
		failedResponse := category_response.FailedResponse{
			Data:    nil,
			Message: "Missing limit or pages query parameter",
			Status:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, failedResponse)
		return
	}

	// Convert limit and offset to integers
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	pages, err := strconv.Atoi(pagesStr)
	if err != nil || pages < 0 {
		pages = 0
	}
	var offset = 0
	if pages > 1 {
		offset = (pages - 1) * limit
	}

	var Sort_by string
	if sortByStr != "asc" && sortByStr != "desc" {
		Sort_by = "asc"
	} else {
		Sort_by = sortByStr
	}

	var Sort_key string
	if sortByStr != "" {
		Sort_key = sortKeyStr
	} else {
		Sort_key = "id"
	}

	category := new([]category_response.GetAllCategory)
	query := database.DB.Table("category")

	// Collect filter conditions
	var conditions []string
	var args []interface{}

	if CategoryStr != "" {
		conditions = append(conditions, "category_name ILIKE ?")
		args = append(args, "%"+CategoryStr+"%")
	}
	if status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	// Apply conditions if there are any
	if len(conditions) > 0 {
		query = query.Where(conditions[0], args...)
		for i := 1; i < len(conditions); i++ {
			query = query.Where(conditions[i], args[i])
		}
	}

	var wg sync.WaitGroup
	var count int64
	var err1, err2 error

	wg.Add(1)
	go func() {
		defer wg.Done()
		countQuery := query.Session(&gorm.Session{})
		err2 = countQuery.Count(&count).Error
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		newQuery := query.Session(&gorm.Session{})
		result := newQuery.Limit(limit).Offset(offset).Order(Sort_key + " " + Sort_by).Find(&category)
		if result.Error != nil {
			err1 = result.Error
		}
	}()

	wg.Wait()

	if err1 != nil {
		notfound := category_response.NewFailedResponse(nil, "Privilege not found", http.StatusNotFound, "Privilege not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	} else if err2 != nil {
		existingResponse := category_response.NewFailedResponse(nil, "Count privilege error", http.StatusConflict, "Count privilege error")
		c.JSON(http.StatusConflict, existingResponse)
		return
	}
	// Create a success response using the SuccessResponse struct from the response package
	successResponse := category_response.Responses{
		Data:    category,
		Message: "category retrieved successfully",
		Status:  http.StatusOK,
		Meta_data: category_response.MetaData{
			Limit:    limit,
			Pages:    pages,
			Total:    int(count),
			Sort_by:  Sort_by,
			Sort_key: Sort_key,
		},
	}

	// Return the success response
	c.JSON(http.StatusOK, successResponse)
}

func GetOneCategory(c *gin.Context) {
	category := new(category_response.GetAllCategory)
	result := database.DB.Table("category").Where("id = ?", c.Param("id")).Find(&category).First(category)
	if result.Error != nil {
		log.Printf("Error retrieving category: %s", result.Error.Error())

		failedResponse := category_response.NewFailedResponse(nil, result.Error.Error(), http.StatusBadRequest, result.Error.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, failedResponse)
		return
	}

	if result.RowsAffected == 0 {
		notFoundResponse := category_response.NewFailedResponse(nil, "category not found", http.StatusNotFound, "category not found")

		c.AbortWithStatusJSON(http.StatusNotFound, notFoundResponse)
		return
	}
	successResponse := category_response.OneResponse(category, "category retrieved successfully", http.StatusOK)

	c.JSON(http.StatusOK, successResponse)
}

func CreateOneCategory(c *gin.Context) {
	var newCategory category_model.CategoryTable
	// log.Print(c.ShouldBindJSON(&newCategory))
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		validation := category_response.NewFailedResponse(nil, "Validation Error", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, validation)
		return
	}

	var existingCategory category_model.CategoryTable
	result := database.DB.Table("category").Where("category_name = ?", newCategory.CategoryName).First(&existingCategory)
	if result.Error == nil {
		existingResponse := category_response.NewFailedResponse(nil, "Category already exists", http.StatusConflict, "Category already exists")
		c.JSON(http.StatusConflict, existingResponse)
		return
	}

	createResult := database.DB.Table("category").Create(&newCategory)
	if createResult.Error != nil {
		log.Printf("Error creating category: %s", createResult.Error.Error())
		failedResponse := category_response.NewFailedResponse(nil, "Error creating category", http.StatusInternalServerError, "Error creating user")
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// User created successfully, return a success response
	successResponse := category_response.OneResponse(newCategory.ID, "Category created successfully", http.StatusCreated)
	c.JSON(http.StatusCreated, successResponse)
}

func UpdateOneCategory(c *gin.Context) {
	var newCategory category_model.CategoryTable
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		validation := category_response.NewFailedResponse(nil, "Validation Error", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, validation)
		return
	}

	var wg sync.WaitGroup
	var existingCategory, existingCategoryByName category_model.CategoryTable
	var err1, err2 error

	// Check if the user exists by ID
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("category").Where("id = ?", c.Params.ByName("id")).First(&existingCategory)
		if result.Error != nil {
			err1 = result.Error
		}
	}()

	// Check if the user exists by full name or email
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("category").Where("(category_name = ? ) and id != ?", newCategory.CategoryName, c.Params.ByName("id")).First(&existingCategoryByName)
		if result.Error != nil {
			err2 = result.Error
		}
	}()

	wg.Wait()

	// If either query found a user, return a conflict response
	log.Default().Println(existingCategory, existingCategoryByName)
	log.Default().Println(err2)
	if err1 != nil {
		notfound := category_response.NewFailedResponse(nil, "Category not found", http.StatusNotFound, "Category not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	} else if err2 == nil {
		existingResponse := category_response.NewFailedResponse(nil, "Category already exists", http.StatusConflict, "Category already exists")
		c.JSON(http.StatusConflict, existingResponse)
		return
	}

	updatedResult := database.DB.Table("category").Where("id = ?", c.Params.ByName("id")).UpdateColumns(&newCategory)
	if updatedResult.Error != nil {
		log.Printf("Error updating category: %s", updatedResult.Error.Error())
		failedResponse := category_response.NewFailedResponse(nil, "Error updating category", http.StatusInternalServerError, updatedResult.Error)
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// // User update successfully, return a success response
	successResponse := category_response.OneResponse(c.Params.ByName("id"), "Category updated successfully", http.StatusOK)
	c.JSON(http.StatusCreated, successResponse)
}

func DeleteOneCategory(c *gin.Context) {
	var wg sync.WaitGroup
	var existingCategory category_model.CategoryTable
	var err1 error

	// Check if the user exists by ID
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("category").Where("id = ?", c.Params.ByName("id")).First(&existingCategory)
		if result.Error != nil {
			err1 = result.Error
		}
	}()

	wg.Wait()

	if err1 != nil {
		notfound := category_response.NewFailedResponse(nil, "Category not found", http.StatusNotFound, "Category not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	}

	deletedResult := database.DB.Table("category").Where("id = ?", c.Params.ByName("id")).Delete(&existingCategory)
	if deletedResult.Error != nil {
		log.Printf("Error deleting category: %s", deletedResult.Error.Error())
		failedResponse := category_response.NewFailedResponse(nil, "Error deleting category", http.StatusInternalServerError, deletedResult.Error)
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// // User update successfully, return a success response
	successResponse := category_response.OneResponse(c.Params.ByName("id"), "Category deleted successfully", http.StatusOK)
	c.JSON(http.StatusCreated, successResponse)
}
