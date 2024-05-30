package privilege_controller

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/database"
	privilege_model "main.go/models/public/privilege"
	privilege_response "main.go/response"
)

func GetAllPrivilege(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	pagesStr := c.DefaultQuery("pages", "1")
	sortKeyStr := c.DefaultQuery("sort_key", "id")
	sortByStr := c.DefaultQuery("sort_by", "asc")
	privilegeStr := c.DefaultQuery("privilege_name", "")
	status := c.DefaultQuery("status", "")

	if limitStr == "" || pagesStr == "" {
		failedResponse := privilege_response.FailedResponse{
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

	privilege := new([]privilege_response.GetAllPrivilege)
	query := database.DB.Table("privilege")

	// Collect filter conditions
	var conditions []string
	var args []interface{}

	if privilegeStr != "" {
		conditions = append(conditions, "privilege_name ILIKE ?")
		args = append(args, "%"+privilegeStr+"%")
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
		result := newQuery.Limit(limit).Offset(offset).Order(Sort_key + " " + Sort_by).Find(&privilege)
		if result.Error != nil {
			err1 = result.Error
		}
	}()

	wg.Wait()

	if err1 != nil {
		notfound := privilege_response.NewFailedResponse(nil, "Privilege not found", http.StatusNotFound, "Privilege not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	} else if err2 != nil {
		existingResponse := privilege_response.NewFailedResponse(nil, "Count privilege error", http.StatusConflict, "Count privilege error")
		c.JSON(http.StatusConflict, existingResponse)
		return
	}
	// Create a success response using the SuccessResponse struct from the response package
	successResponse := privilege_response.Responses{
		Data:    privilege,
		Message: "privilege retrieved successfully",
		Status:  http.StatusOK,
		Meta_data: privilege_response.MetaData{
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

func GetOnePrivilege(c *gin.Context) {
	privilege := new(privilege_response.GetAllPrivilege)
	result := database.DB.Table("privilege").Where("id = ?", c.Param("id")).Find(&privilege).First(privilege)
	if result.Error != nil {
		log.Printf("Error retrieving privilege: %s", result.Error.Error())

		failedResponse := privilege_response.NewFailedResponse(nil, result.Error.Error(), http.StatusBadRequest, result.Error.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, failedResponse)
		return
	}

	if result.RowsAffected == 0 {
		notFoundResponse := privilege_response.NewFailedResponse(nil, "Privilege not found", http.StatusNotFound, "Privilege not found")

		c.AbortWithStatusJSON(http.StatusNotFound, notFoundResponse)
		return
	}
	successResponse := privilege_response.OneResponse(privilege, "privilege retrieved successfully", http.StatusOK)

	c.JSON(http.StatusOK, successResponse)
}

func CreateOnePrivilege(c *gin.Context) {
	var newPrivilege privilege_model.PrivilegeTable
	// log.Print(c.ShouldBindJSON(&newPrivilege))
	if err := c.ShouldBindJSON(&newPrivilege); err != nil {
		validation := privilege_response.NewFailedResponse(nil, "Validation Error", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, validation)
		return
	}

	var existingPrivilege privilege_model.PrivilegeTable
	result := database.DB.Table("privilege").Where("privilege_name = ?", newPrivilege.PrivilegeName).First(&existingPrivilege)
	if result.Error == nil {
		existingResponse := privilege_response.NewFailedResponse(nil, "Privilege already exists", http.StatusConflict, "Privilege already exists")
		c.JSON(http.StatusConflict, existingResponse)
		return
	}

	createResult := database.DB.Table("privilege").Create(&newPrivilege)
	if createResult.Error != nil {
		log.Printf("Error creating privilege: %s", createResult.Error.Error())
		failedResponse := privilege_response.NewFailedResponse(nil, "Error creating user", http.StatusInternalServerError, "Error creating user")
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// User created successfully, return a success response
	successResponse := privilege_response.OneResponse(newPrivilege.ID, "Privilege created successfully", http.StatusCreated)
	c.JSON(http.StatusCreated, successResponse)
}

func UpdateOnePrivilege(c *gin.Context) {
	var newPrivilege privilege_model.PrivilegeTable
	if err := c.ShouldBindJSON(&newPrivilege); err != nil {
		validation := privilege_response.NewFailedResponse(nil, "Validation Error", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, validation)
		return
	}

	var wg sync.WaitGroup
	var existingPrivilege, existingPrivilegeByName privilege_model.PrivilegeTable
	var err1, err2 error

	// Check if the user exists by ID
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("privilege").Where("id = ?", c.Params.ByName("id")).First(&existingPrivilege)
		if result.Error != nil {
			err1 = result.Error
		}
	}()

	// Check if the user exists by full name or email
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("privilege").Where("(privilege_name = ? ) and id != ?", newPrivilege.PrivilegeName, c.Params.ByName("id")).First(&existingPrivilegeByName)
		if result.Error != nil {
			err2 = result.Error
		}
	}()

	wg.Wait()

	// If either query found a user, return a conflict response
	log.Default().Println(existingPrivilege, existingPrivilegeByName)
	log.Default().Println(err2)
	if err1 != nil {
		notfound := privilege_response.NewFailedResponse(nil, "Privilege not found", http.StatusNotFound, "Privilege not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	} else if err2 == nil {
		existingResponse := privilege_response.NewFailedResponse(nil, "Privilege already exists", http.StatusConflict, "Privilege already exists")
		c.JSON(http.StatusConflict, existingResponse)
		return
	}

	updatedResult := database.DB.Table("privilege").Where("id = ?", c.Params.ByName("id")).UpdateColumns(&newPrivilege)
	if updatedResult.Error != nil {
		log.Printf("Error updating privilege: %s", updatedResult.Error.Error())
		failedResponse := privilege_response.NewFailedResponse(nil, "Error updating privilege", http.StatusInternalServerError, updatedResult.Error)
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// // User update successfully, return a success response
	successResponse := privilege_response.OneResponse(c.Params.ByName("id"), "Privilege updated successfully", http.StatusOK)
	c.JSON(http.StatusCreated, successResponse)
}

func DeleteOnePrivilege(c *gin.Context) {
	var wg sync.WaitGroup
	var existingPrivilege privilege_model.PrivilegeTable
	var err1 error

	// Check if the user exists by ID
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("privilege").Where("id = ?", c.Params.ByName("id")).First(&existingPrivilege)
		if result.Error != nil {
			err1 = result.Error
		}
	}()

	wg.Wait()

	if err1 != nil {
		notfound := privilege_response.NewFailedResponse(nil, "Privilege not found", http.StatusNotFound, "Privilege not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	}

	deletedResult := database.DB.Table("privilege").Where("id = ?", c.Params.ByName("id")).Delete(&existingPrivilege)
	if deletedResult.Error != nil {
		log.Printf("Error deleting privilege: %s", deletedResult.Error.Error())
		failedResponse := privilege_response.NewFailedResponse(nil, "Error deleting privilege", http.StatusInternalServerError, deletedResult.Error)
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// // User update successfully, return a success response
	successResponse := privilege_response.OneResponse(c.Params.ByName("id"), "Privilege deleted successfully", http.StatusOK)
	c.JSON(http.StatusCreated, successResponse)
}
