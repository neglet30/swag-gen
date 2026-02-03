package api

// GetUsers 获取所有用户
// @Router /api/users [GET]
// @Summary 获取所有用户
// @Description 从数据库获取所有用户
// @Tags User
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {array} User "成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "服务器错误"
func GetUsers() {}

// GetUserByID 根据 ID 获取用户
// @Router /api/users/{id} [GET]
// @Summary 根据 ID 获取用户
// @Description 根据用户 ID 获取用户详情
// @Tags User
// @Param id path int true "用户 ID"
// @Success 200 {object} User "成功"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
func GetUserByID(id int) {}

// CreateUser 创建用户
// @Router /api/users [POST]
// @Summary 创建用户
// @Description 创建新用户
// @Tags User
// @Param body body CreateUserRequest true "用户信息"
// @Success 201 {object} User "创建成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "服务器错误"
func CreateUser(req CreateUserRequest) {}

// UpdateUser 更新用户
// @Router /api/users/{id} [PUT]
// @Summary 更新用户
// @Description 更新用户信息
// @Tags User
// @Param id path int true "用户 ID"
// @Param body body UpdateUserRequest true "用户信息"
// @Success 200 {object} User "更新成功"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
func UpdateUser(id int, req UpdateUserRequest) {}

// DeleteUser 删除用户
// @Router /api/users/{id} [DELETE]
// @Summary 删除用户
// @Description 删除用户
// @Tags User
// @Param id path int true "用户 ID"
// @Success 204 "删除成功"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
func DeleteUser(id int) {}

// User 用户模型
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
