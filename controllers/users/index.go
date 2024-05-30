package users_controller

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/database"
	users_model "main.go/models/public/users"
	users_response "main.go/response"

	"main.go/helper"
)

func GetAllUsers(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	pagesStr := c.DefaultQuery("pages", "1")
	sortKeyStr := c.DefaultQuery("sort_key", "id")
	sortByStr := c.DefaultQuery("sort_by", "asc")
	full_name := c.DefaultQuery("full_name", "")
	email := c.DefaultQuery("email", "")
	status := c.DefaultQuery("status", "")

	if limitStr == "" || pagesStr == "" {
		failedResponse := users_response.FailedResponse{
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

	users := new([]users_response.GetAllUsersResponse)
	query := database.DB.Table("users")

	// Collect filter conditions
	var conditions []string
	var args []interface{}

	if full_name != "" {
		conditions = append(conditions, "full_name ILIKE ?")
		args = append(args, "%"+full_name+"%")
	}
	if email != "" {
		conditions = append(conditions, "email ILIKE ?")
		args = append(args, "%"+email+"%")
	}
	if status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	log.Default().Println("condition", conditions)

	// Apply conditions if there are any
	if len(conditions) > 0 {
		query = query.Where(conditions[0], args...)
		for i := 1; i < len(conditions); i++ {
			query = query.Where(conditions[i], args[i])
		}
	}

	log.Default().Println("query")
	log.Default().Println(query.Statement.Vars...)

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
		result := newQuery.Joins("left join users_privilege on users.id = users_privilege.id_users").Limit(limit).Offset(offset).Order(Sort_key + " " + Sort_by).Find(&users)
		if result.Error != nil {
			err1 = result.Error
		}
		log.Default().Println("result", result)
	}()

	wg.Wait()

	if err1 != nil {
		notfound := users_response.NewFailedResponse(nil, "User not found", http.StatusNotFound, "User not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	} else if err2 != nil {
		existingResponse := users_response.NewFailedResponse(nil, "Count user error", http.StatusConflict, "Count user error")
		c.JSON(http.StatusConflict, existingResponse)
		return
	}

	// Create a success response using the SuccessResponse struct from the response package
	successResponse := users_response.Responses{
		Data:    users,
		Message: "Users retrieved successfully",
		Status:  http.StatusOK,
		Meta_data: users_response.MetaData{
			Limit:    limit,
			Total:    int(count),
			Pages:    pages,
			Sort_by:  Sort_by,
			Sort_key: Sort_key,
		},
	}

	// Return the success response
	c.JSON(http.StatusOK, successResponse)
}

func GetOneUsers(c *gin.Context) {
	users := new(users_response.GetAllUsersResponse)
	result := database.DB.Table("users").Where("id = ?", c.Param("id")).Find(&users).First(users)
	if result.Error != nil {
		log.Printf("Error retrieving user: %s", result.Error.Error())

		failedResponse := users_response.NewFailedResponse(nil, result.Error.Error(), http.StatusBadRequest, result.Error.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, failedResponse)
		return
	}

	if result.RowsAffected == 0 {
		notFoundResponse := users_response.NewFailedResponse(nil, "User not found", http.StatusNotFound, "User not found")

		c.AbortWithStatusJSON(http.StatusNotFound, notFoundResponse)
		return
	}
	successResponse := users_response.OneResponse(users, "Users retrieved successfully", http.StatusOK)

	c.JSON(http.StatusOK, successResponse)
}

func CreateOneUsers(c *gin.Context) {
	// Parse JSON request body into a users_model.UserTable instance
	var newUser users_model.UserTable
	// log.Print(c.ShouldBindJSON(&newUser))
	if err := c.ShouldBindJSON(&newUser); err != nil {
		validation := users_response.NewFailedResponse(nil, "Validation Error", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, validation)
		return
	}

	var existingUser users_model.UserTable
	result := database.DB.Table("users").Where("full_name = ? OR email = ?", newUser.FullName, newUser.Email).First(&existingUser)
	if result.Error == nil {
		existingResponse := users_response.NewFailedResponse(nil, "User already exists", http.StatusConflict, "User already exists")
		c.JSON(http.StatusConflict, existingResponse)
		return
	}

	hashedPassword, err := helper.HashPassword(newUser.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err.Error())
		failedResponse := users_response.NewFailedResponse(nil, "Error hashing password", http.StatusInternalServerError, "Error hashing password")
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}
	newUser.Password = hashedPassword

	createResult := database.DB.Table("users").Create(&newUser)
	if createResult.Error != nil {
		log.Printf("Error creating user: %s", createResult.Error.Error())
		failedResponse := users_response.NewFailedResponse(nil, "Error creating user", http.StatusInternalServerError, "Error creating user")
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// User created successfully, return a success response
	successResponse := users_response.OneResponse(newUser.ID, "User created successfully", http.StatusCreated)
	c.JSON(http.StatusCreated, successResponse)
}

func UpdateOneUsers(c *gin.Context) {
	var newUser users_model.UpdateUserTable
	if err := c.ShouldBindJSON(&newUser); err != nil {
		validation := users_response.NewFailedResponse(nil, "Validation Error", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, validation)
		return
	}

	var wg sync.WaitGroup
	var existingUserByID, existingUserByEmailOrName users_model.UpdateUserTable
	var err1, err2 error

	// Check if the user exists by ID
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("users").Where("id = ?", c.Params.ByName("id")).First(&existingUserByID)
		if result.Error != nil {
			err1 = result.Error
		}
	}()

	// Check if the user exists by full name or email
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("users").Where("(full_name = ? OR email = ?) and id != ?", newUser.FullName, newUser.Email, c.Params.ByName("id")).First(&existingUserByEmailOrName)
		if result.Error != nil {
			err2 = result.Error
		}
	}()

	wg.Wait()

	// If either query found a user, return a conflict response
	log.Default().Println(existingUserByID, existingUserByEmailOrName)
	log.Default().Println(err2)
	if err1 != nil {
		notfound := users_response.NewFailedResponse(nil, "User not found", http.StatusNotFound, "User not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	} else if err2 == nil {
		existingResponse := users_response.NewFailedResponse(nil, "User already exists", http.StatusConflict, "User already exists")
		c.JSON(http.StatusConflict, existingResponse)
		return
	}

	updatedResult := database.DB.Table("users").Where("id = ?", c.Params.ByName("id")).UpdateColumns(&newUser)
	if updatedResult.Error != nil {
		log.Printf("Error updating user: %s", updatedResult.Error.Error())
		failedResponse := users_response.NewFailedResponse(nil, "Error updating user", http.StatusInternalServerError, updatedResult.Error)
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// // User update successfully, return a success response
	successResponse := users_response.OneResponse(c.Params.ByName("id"), "User updated successfully", http.StatusOK)
	c.JSON(http.StatusCreated, successResponse)
}

func DeleteOneUsers(c *gin.Context) {
	var wg sync.WaitGroup
	var existingUserByID users_model.UpdateUserTable
	var err1 error

	// Check if the user exists by ID
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("users").Where("id = ?", c.Params.ByName("id")).First(&existingUserByID)
		if result.Error != nil {
			err1 = result.Error
		}
	}()

	wg.Wait()

	if err1 != nil {
		notfound := users_response.NewFailedResponse(nil, "User not found", http.StatusNotFound, "User not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	}

	deletedResult := database.DB.Table("users").Where("id = ?", c.Params.ByName("id")).Delete(&existingUserByID)
	if deletedResult.Error != nil {
		log.Printf("Error updating user: %s", deletedResult.Error.Error())
		failedResponse := users_response.NewFailedResponse(nil, "Error updating user", http.StatusInternalServerError, deletedResult.Error)
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// // User update successfully, return a success response
	successResponse := users_response.OneResponse(c.Params.ByName("id"), "User deleted successfully", http.StatusOK)
	c.JSON(http.StatusCreated, successResponse)
}
