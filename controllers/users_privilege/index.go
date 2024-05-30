package users_controller

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/database"
	privilege "main.go/models/public/privilege"
	users "main.go/models/public/users"
	users_privilege "main.go/models/public/users_privilege"
	users_privilege_response "main.go/response"
)

func GetAllUsersPrivilege(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	pagesStr := c.DefaultQuery("pages", "1")
	sortKeyStr := c.DefaultQuery("sort_key", "id")
	sortByStr := c.DefaultQuery("sort_by", "asc")
	id_usersStr := c.DefaultQuery("id_users", "")
	id_privilegeStr := c.DefaultQuery("id_privilege", "")
	status := c.DefaultQuery("status", "")

	if limitStr == "" || pagesStr == "" {
		failedResponse := users_privilege_response.FailedResponse{
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

	var id_users int
	if id_usersStr != "" {
		id_users_, err := strconv.Atoi(id_usersStr)
		if err != nil || id_users_ < 0 {
			id_users = 0
		} else {
			id_users = id_users_
		}
	}

	var id_privilege int
	if id_privilegeStr != "" {
		id_privilege_, err := strconv.Atoi(id_privilegeStr)
		if err != nil || id_privilege_ < 0 {
			id_privilege = 0
		} else {
			id_privilege = id_privilege_
		}
	}

	users_privilege := new([]users_privilege_response.GetAllUsersPrivilege)
	query := database.DB.Table("users_privilege")

	// Collect filter conditions
	var conditions []string
	var args []interface{}

	log.Default().Println("id_users", id_users)

	if id_users != 0 {
		conditions = append(conditions, "id_users = ?")
		args = append(args, id_users)
	}

	if id_privilege != 0 {
		conditions = append(conditions, "id_privilege = ?")
		args = append(args, id_privilege)
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

	fields := `
	users.full_name as full_name,
	users.email as email,
	users.status as status_users,
	privilege.privilege_name,
	users_privilege.status as status_users_privilege,
	users_privilege.id as id,
	users_privilege.id_users as id_users,
	users_privilege.id_privilege as id_privilege,
	users_privilege."createdAt" as created_at,
	users_privilege."updatedAt" as updated_at,
	privilege.status as status_privilege
	`

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
		result := newQuery.Limit(limit).
			Offset(offset).
			Order(Sort_key + " " + Sort_by).
			Select(fields).
			Joins(`left join users on users.id = users_privilege.id_users`).
			Joins(`left join privilege on privilege.id = users_privilege.id_privilege`).
			Find(&users_privilege)
		if result.Error != nil {
			err1 = result.Error
		}
		log.Default().Println("result", result)
	}()

	wg.Wait()

	if err1 != nil {
		notfound := users_privilege_response.NewFailedResponse(nil, "User not found", http.StatusNotFound, "User not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	} else if err2 != nil {
		existingResponse := users_privilege_response.NewFailedResponse(nil, "Count user error", http.StatusConflict, "Count user error")
		c.JSON(http.StatusConflict, existingResponse)
		return
	}

	// Create a success response using the SuccessResponse struct from the response package
	successResponse := users_privilege_response.Responses{
		Data:    users_privilege,
		Message: "Users privilege retrieved successfully",
		Status:  http.StatusOK,
		Meta_data: users_privilege_response.MetaData{
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

func GetOneUsersPrivilege(c *gin.Context) {
	users_privilege := new(users_privilege_response.GetAllUsersPrivilege)

	fields := `
	users.full_name as full_name,
	users.email as email,
	users.status as status_users,
	privilege.privilege_name,
	users_privilege.status as status_users_privilege,
	users_privilege.id as id,
	users_privilege.id_users as id_users,
	users_privilege.id_privilege as id_privilege,
	users_privilege."createdAt" as created_at,
	users_privilege."updatedAt" as updated_at,
	privilege.status as status_privilege
	`

	result := database.DB.Table("users_privilege").Select(fields).Joins(`left join users on users.id = users_privilege.id_users`).
		Joins(`left join privilege on privilege.id = users_privilege.id_privilege`).Where("users_privilege.id = ?", c.Param("id")).Find(&users_privilege).First(users_privilege)
	if result.Error != nil {
		log.Printf("Error retrieving user: %s", result.Error.Error())

		failedResponse := users_privilege_response.NewFailedResponse(nil, result.Error.Error(), http.StatusBadRequest, result.Error.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, failedResponse)
		return
	}

	if result.RowsAffected == 0 {
		notFoundResponse := users_privilege_response.NewFailedResponse(nil, "User not found", http.StatusNotFound, "User not found")

		c.AbortWithStatusJSON(http.StatusNotFound, notFoundResponse)
		return
	}
	successResponse := users_privilege_response.OneResponse(users_privilege, "Users privilege retrieved successfully", http.StatusOK)

	c.JSON(http.StatusOK, successResponse)
}

func GetUsersPrivilegeByUser(c *gin.Context) {
	users_privilege := new([]users_privilege_response.GetAllUsersPrivilege)

	fields := `
	users.full_name as full_name,
	users.email as email,
	users.status as status_users,
	privilege.privilege_name,
	users_privilege.status as status_users_privilege,
	users_privilege.id as id,
	users_privilege.id_users as id_users,
	users_privilege.id_privilege as id_privilege,
	users_privilege."createdAt" as created_at,
	users_privilege."updatedAt" as updated_at,
	privilege.status as status_privilege
	`

	result := database.DB.Table("users_privilege").Select(fields).Joins(`left join users on users.id = users_privilege.id_users`).
		Joins(`left join privilege on privilege.id = users_privilege.id_privilege`).Where("users_privilege.id_users = ?", c.Param("id")).Find(&users_privilege)
	if result.Error != nil {
		log.Printf("Error retrieving user privilege: %s", result.Error.Error())

		failedResponse := users_privilege_response.NewFailedResponse(nil, result.Error.Error(), http.StatusBadRequest, result.Error.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, failedResponse)
		return
	}

	if result.RowsAffected == 0 {
		notFoundResponse := users_privilege_response.NewFailedResponse(nil, "User not found", http.StatusNotFound, "User not found")

		c.AbortWithStatusJSON(http.StatusNotFound, notFoundResponse)
		return
	}
	successResponse := users_privilege_response.OneResponse(users_privilege, "Users privilege retrieved successfully", http.StatusOK)

	c.JSON(http.StatusOK, successResponse)
}

func CreateOneUsers(c *gin.Context) {
	// Parse JSON request body into a users_model.UserTable instance
	var newUser users_privilege.UsersPrivilegeTable
	// log.Print(c.ShouldBindJSON(&newUser))
	if err := c.ShouldBindJSON(&newUser); err != nil {
		validation := users_privilege_response.NewFailedResponse(nil, "Validation Error", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, validation)
		return
	}

	var wg sync.WaitGroup
	var users users.UserTable
	var privilege privilege.PrivilegeTable
	var existing users_privilege.UsersPrivilegeTable
	var err1, err2, err3 error

	//find users
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("users").Where("id = ?", newUser.IdUsers).Find(&users).First(&users)
		if result.Error != nil {
			err1 = result.Error
		}
	}()

	//find users_privilege
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := database.DB.Table("users_privilege").Where("id_users = ? AND id_privilege = ?", newUser.IdUsers, newUser.IdPrivilege).
			Find(&existing).First(&existing).Error
		if err != nil {
			err2 = err
		}
		log.Default().Println("result", existing)
	}()

	//find privilege
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := database.DB.Table("privilege").Where("id = ? AND status is true", newUser.IdPrivilege).
			Find(&privilege).First(&privilege).Error
		if err != nil {
			err3 = err
		}
		log.Default().Println("result", privilege)
	}()

	wg.Wait()

	if err1 != nil {
		notfound := users_privilege_response.NewFailedResponse(nil, "User not found", http.StatusNotFound, "User not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	} else if err2 == nil {
		existingResponse := users_privilege_response.NewFailedResponse(nil, "duplicate data", http.StatusConflict, "duplicate data")
		c.JSON(http.StatusConflict, existingResponse)
		return
	} else if err3 != nil {
		notfound := users_privilege_response.NewFailedResponse(nil, "Privilege not found", http.StatusNotFound, err3.Error())
		c.JSON(http.StatusNotFound, notfound)
		return
	}

	createResult := database.DB.Table("users_privilege").Create(&newUser)
	if createResult.Error != nil {
		log.Printf("Error creating users_privilege: %s", createResult.Error.Error())
		failedResponse := users_privilege_response.NewFailedResponse(nil, "Error creating users_privilege", http.StatusInternalServerError, "Error creating user")
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// User created successfully, return a success response
	successResponse := users_privilege_response.OneResponse(newUser, "User created successfully", http.StatusCreated)
	c.JSON(http.StatusCreated, successResponse)
}

func UpdateOneUsers(c *gin.Context) {

	// Parse JSON request body into a users_model.UserTable instance
	var newUser users_privilege.UpdateUsersPrivilegeTable
	idParam := c.Param("id")
	id, errv := strconv.Atoi(idParam)
	if errv != nil {
		validation := users_privilege_response.NewFailedResponse(nil, "Validation Error", http.StatusBadRequest, "Invalid ID")
		c.JSON(http.StatusBadRequest, validation)
		return
	}
	// log.Print(c.ShouldBindJSON(&newUser))
	if err := c.ShouldBindJSON(&newUser); err != nil {
		validation := users_privilege_response.NewFailedResponse(nil, "Validation Error", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, validation)
		return
	}

	var existing_ users_privilege.UsersPrivilegeTable

	err := database.DB.Table("users_privilege").Where("id = ?", id).Find(&existing_).First(&existing_).Error
	if err != nil {
		notfound := users_privilege_response.NewFailedResponse(nil, "Users privilege not found", http.StatusNotFound, "Users privilege not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	}

	var wg sync.WaitGroup
	var privilege privilege.PrivilegeTable
	var existing users_privilege.UsersPrivilegeTable
	var err1, err2 error

	//find users_privilege
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := database.DB.Table("users_privilege").Where("id_users = ? AND id_privilege = ? AND id != ?", existing_.IdUsers, newUser.IdPrivilege, id).
			Find(&existing).First(&existing).Error
		if err != nil {
			err1 = err
		}
		log.Default().Println("result", existing)
	}()

	//find privilege
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := database.DB.Table("privilege").Where("id = ? AND status is true", newUser.IdPrivilege).
			Find(&privilege).First(&privilege).Error
		if err != nil {
			err2 = err
		}
		log.Default().Println("result", privilege)
	}()

	wg.Wait()

	if err1 == nil {
		existingResponse := users_privilege_response.NewFailedResponse(nil, "duplicate data", http.StatusConflict, "duplicate data")
		c.JSON(http.StatusConflict, existingResponse)
		return
	} else if err2 != nil {
		notfound := users_privilege_response.NewFailedResponse(nil, "Privilege not found", http.StatusNotFound, err2.Error())
		c.JSON(http.StatusNotFound, notfound)
		return
	}

	updatedResult := database.DB.Table("users_privilege").Where("id = ?", id).UpdateColumns(&newUser)
	if updatedResult.Error != nil {
		log.Printf("Error updating users_privilege: %s", updatedResult.Error.Error())
		failedResponse := users_privilege_response.NewFailedResponse(nil, "Error updating users_privilege", http.StatusInternalServerError, updatedResult.Error.Error())
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// User created successfully, return a success response
	successResponse := users_privilege_response.OneResponse(newUser, "User privilege updated successfully", http.StatusCreated)
	c.JSON(http.StatusCreated, successResponse)
}

func DeleteOneUsers(c *gin.Context) {
	var wg sync.WaitGroup
	var existingUserByID users_privilege.UsersPrivilegeTable
	var err1 error

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if the user exists by ID
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := database.DB.Table("users_privilege").Where("id = ?", id).First(&existingUserByID)
		if result.Error != nil {
			err1 = result.Error
		}
	}()

	wg.Wait()

	if err1 != nil {
		notfound := users_privilege_response.NewFailedResponse(nil, "User privilege not found", http.StatusNotFound, "User privilege not found")
		c.JSON(http.StatusNotFound, notfound)
		return
	}

	deletedResult := database.DB.Table("users_privilege").Where("id = ?", id).Delete(&existingUserByID)
	if deletedResult.Error != nil {
		log.Printf("Error deleting user privilege: %s", deletedResult.Error.Error())
		failedResponse := users_privilege_response.NewFailedResponse(nil, "Error deleting users privilege", http.StatusInternalServerError, deletedResult.Error)
		c.JSON(http.StatusInternalServerError, failedResponse)
		return
	}

	// // User update successfully, return a success response
	successResponse := users_privilege_response.OneResponse(id, "Users privilege deleted successfully", http.StatusOK)
	c.JSON(http.StatusCreated, successResponse)
}
