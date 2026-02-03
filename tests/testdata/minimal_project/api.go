package main

// GetHealth 健康检查
// @Router /health [GET]
// @Summary 健康检查
// @Tags Health
// @Success 200 {object} HealthResponse
func GetHealth() {}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status string `json:"status"`
}
