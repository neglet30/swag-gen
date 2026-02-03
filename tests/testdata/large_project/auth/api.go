package auth

// Login 用户登录
// @Router /api/auth/login [POST]
// @Summary 用户登录
// @Tags Auth
// @Success 200 {object} LoginResponse
func Login() {}

// Logout 用户登出
// @Router /api/auth/logout [POST]
// @Summary 用户登出
// @Tags Auth
// @Success 204
func Logout() {}

// RefreshToken 刷新令牌
// @Router /api/auth/refresh [POST]
// @Summary 刷新令牌
// @Tags Auth
// @Success 200 {object} TokenResponse
func RefreshToken() {}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
}

// TokenResponse 令牌响应
type TokenResponse struct {
	Token string `json:"token"`
}
