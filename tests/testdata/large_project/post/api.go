package post

// GetPosts 获取文章列表
// @Router /api/posts [GET]
// @Summary 获取文章列表
// @Tags Post
// @Param page query int false "页码"
// @Success 200 {array} Post
func GetPosts() {}

// GetPostByID 获取文章详情
// @Router /api/posts/{id} [GET]
// @Summary 获取文章详情
// @Tags Post
// @Param id path int true "文章 ID"
// @Success 200 {object} Post
func GetPostByID(id int) {}

// CreatePost 创建文章
// @Router /api/posts [POST]
// @Summary 创建文章
// @Tags Post
// @Success 201 {object} Post
func CreatePost() {}

// UpdatePost 更新文章
// @Router /api/posts/{id} [PUT]
// @Summary 更新文章
// @Tags Post
// @Param id path int true "文章 ID"
// @Success 200 {object} Post
func UpdatePost(id int) {}

// DeletePost 删除文章
// @Router /api/posts/{id} [DELETE]
// @Summary 删除文章
// @Tags Post
// @Param id path int true "文章 ID"
// @Success 204
func DeletePost(id int) {}

// Post 文章模型
type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
