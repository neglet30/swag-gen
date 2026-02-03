package api

// GetPosts 获取所有文章
// @Router /api/posts [GET]
// @Summary 获取所有文章
// @Description 从数据库获取所有文章
// @Tags Post
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param userId query int false "用户 ID"
// @Success 200 {array} Post "成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "服务器错误"
func GetPosts() {}

// GetPostByID 根据 ID 获取文章
// @Router /api/posts/{id} [GET]
// @Summary 根据 ID 获取文章
// @Description 根据文章 ID 获取文章详情
// @Tags Post
// @Param id path int true "文章 ID"
// @Success 200 {object} Post "成功"
// @Failure 404 {object} ErrorResponse "文章不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
func GetPostByID(id int) {}

// CreatePost 创建文章
// @Router /api/posts [POST]
// @Summary 创建文章
// @Description 创建新文章
// @Tags Post
// @Param body body CreatePostRequest true "文章信息"
// @Success 201 {object} Post "创建成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "服务器错误"
func CreatePost(req CreatePostRequest) {}

// UpdatePost 更新文章
// @Router /api/posts/{id} [PUT]
// @Summary 更新文章
// @Description 更新文章信息
// @Tags Post
// @Param id path int true "文章 ID"
// @Param body body UpdatePostRequest true "文章信息"
// @Success 200 {object} Post "更新成功"
// @Failure 404 {object} ErrorResponse "文章不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
func UpdatePost(id int, req UpdatePostRequest) {}

// DeletePost 删除文章
// @Router /api/posts/{id} [DELETE]
// @Summary 删除文章
// @Description 删除文章
// @Tags Post
// @Param id path int true "文章 ID"
// @Success 204 "删除成功"
// @Failure 404 {object} ErrorResponse "文章不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
func DeletePost(id int) {}

// Post 文章模型
type Post struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
