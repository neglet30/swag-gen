package comment

// GetComments 获取评论列表
// @Router /api/comments [GET]
// @Summary 获取评论列表
// @Tags Comment
// @Param page query int false "页码"
// @Success 200 {array} Comment
func GetComments() {}

// GetCommentByID 获取评论详情
// @Router /api/comments/{id} [GET]
// @Summary 获取评论详情
// @Tags Comment
// @Param id path int true "评论 ID"
// @Success 200 {object} Comment
func GetCommentByID(id int) {}

// CreateComment 创建评论
// @Router /api/comments [POST]
// @Summary 创建评论
// @Tags Comment
// @Success 201 {object} Comment
func CreateComment() {}

// DeleteComment 删除评论
// @Router /api/comments/{id} [DELETE]
// @Summary 删除评论
// @Tags Comment
// @Param id path int true "评论 ID"
// @Success 204
func DeleteComment(id int) {}

// Comment 评论模型
type Comment struct {
	ID      int    `json:"id"`
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
}
