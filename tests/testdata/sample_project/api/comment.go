package api

// GetComments 获取所有评论
// @Router /api/comments [GET]
// @Summary 获取所有评论
// @Description 从数据库获取所有评论
// @Tags Comment
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param postId query int false "文章 ID"
// @Success 200 {array} Comment "成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "服务器错误"
func GetComments() {}

// GetCommentByID 根据 ID 获取评论
// @Router /api/comments/{id} [GET]
// @Summary 根据 ID 获取评论
// @Description 根据评论 ID 获取评论详情
// @Tags Comment
// @Param id path int true "评论 ID"
// @Success 200 {object} Comment "成功"
// @Failure 404 {object} ErrorResponse "评论不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
func GetCommentByID(id int) {}

// CreateComment 创建评论
// @Router /api/comments [POST]
// @Summary 创建评论
// @Description 创建新评论
// @Tags Comment
// @Param body body CreateCommentRequest true "评论信息"
// @Success 201 {object} Comment "创建成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "服务器错误"
func CreateComment(req CreateCommentRequest) {}

// DeleteComment 删除评论
// @Router /api/comments/{id} [DELETE]
// @Summary 删除评论
// @Description 删除评论
// @Tags Comment
// @Param id path int true "评论 ID"
// @Success 204 "删除成功"
// @Failure 404 {object} ErrorResponse "评论不存在"
// @Failure 500 {object} ErrorResponse "服务器错误"
func DeleteComment(id int) {}

// Comment 评论模型
type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID  int    `json:"post_id"`
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}
