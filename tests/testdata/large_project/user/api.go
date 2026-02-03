package user

// GetUsers 获取用户列表
// @Router /api/users [GET]
// @Summary 获取用户列表
// @Tags User
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {array} User
func GetUsers() {}

// GetUserByID 获取用户详情
// @Router /api/users/{id} [GET]
// @Summary 获取用户详情
// @Tags User
// @Param id path int true "用户 ID"
// @Success 200 {object} User
func GetUserByID(id int) {}

// CreateUser 创建用户
// @Router /api/users [POST]
// @Summary 创建用户
// @Tags User
// @Success 201 {object} User
func CreateUser() {}

// UpdateUser 更新用户
// @Router /api/users/{id} [PUT]
// @Summary 更新用户
// @Tags User
// @Param id path int true "用户 ID"
// @Success 200 {object} User
func UpdateUser(id int) {}

// DeleteUser 删除用户
// @Router /api/users/{id} [DELETE]
// @Summary 删除用户
// @Tags User
// @Param id path int true "用户 ID"
// @Success 204
func DeleteUser(id int) {}

// User 用户模型
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
