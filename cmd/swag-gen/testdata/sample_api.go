package api

// User 用户模型
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUsers 获取所有用户
// @Router /api/users [GET]
// @Summary 获取所有用户
// @Description 获取系统中的所有用户
// @Tags User
// @Success 200 {array} User
// @Failure 500 {object} ErrorResponse
func GetUsers() {
	// 实现
}

// GetUser 获取单个用户
// @Router /api/users/{id} [GET]
// @Summary 获取单个用户
// @Description 根据 ID 获取用户信息
// @Tags User
// @Param id path int true "用户 ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func GetUser(id int) {
	// 实现
}

// CreateUser 创建用户
// @Router /api/users [POST]
// @Summary 创建用户
// @Description 创建新用户
// @Tags User
// @Param body body User true "用户信息"
// @Success 201 {object} User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func CreateUser(user User) {
	// 实现
}

// UpdateUser 更新用户
// @Router /api/users/{id} [PUT]
// @Summary 更新用户
// @Description 更新用户信息
// @Tags User
// @Param id path int true "用户 ID"
// @Param body body User true "用户信息"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func UpdateUser(id int, user User) {
	// 实现
}

// DeleteUser 删除用户
// @Router /api/users/{id} [DELETE]
// @Summary 删除用户
// @Description 删除用户
// @Tags User
// @Param id path int true "用户 ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func DeleteUser(id int) {
	// 实现
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
