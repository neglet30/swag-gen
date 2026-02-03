package parser

import (
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// CommentParser 代表注释解析器
type CommentParser struct {
	logger *zap.Logger
}

// NewCommentParser 创建一个新的注释解析器
func NewCommentParser(logger *zap.Logger) *CommentParser {
	return &CommentParser{
		logger: logger,
	}
}

// ParseEndpoint 从注释中解析端点信息
func (cp *CommentParser) ParseEndpoint(comments []string, filePath string, line int) *Endpoint {
	if len(comments) == 0 {
		return nil
	}

	endpoint := &Endpoint{
		File:       filePath,
		Line:       line,
		Tags:       make([]string, 0),
		Parameters: make([]Parameter, 0),
		Responses:  make(map[string]Response),
	}

	hasRouter := false

	for _, comment := range comments {
		// 移除注释前缀
		text := strings.TrimPrefix(comment, "//")
		text = strings.TrimSpace(text)

		// 解析 @Router 标签
		if router := cp.parseRouter(text); router != nil {
			endpoint.Method = router.Method
			endpoint.Path = router.Path
			hasRouter = true
		}

		// 解析 @Summary 标签
		if summary := cp.parseSimpleTag(text, "@Summary"); summary != "" {
			endpoint.Summary = summary
		}

		// 解析 @Description 标签
		if desc := cp.parseSimpleTag(text, "@Description"); desc != "" {
			endpoint.Description = desc
		}

		// 解析 @Tags 标签
		if tags := cp.parseSimpleTag(text, "@Tags"); tags != "" {
			endpoint.Tags = append(endpoint.Tags, tags)
		}

		// 解析 @Param 标签
		if param := cp.parseParam(text); param != nil {
			endpoint.Parameters = append(endpoint.Parameters, *param)
		}

		// 解析 @Success 标签
		if resp := cp.parseResponse(text, "@Success"); resp != nil {
			endpoint.Responses[resp.StatusCode] = *resp
		}

		// 解析 @Failure 标签
		if resp := cp.parseResponse(text, "@Failure"); resp != nil {
			endpoint.Responses[resp.StatusCode] = *resp
		}

		// 解析 @Deprecated 标签
		if strings.Contains(text, "@Deprecated") {
			endpoint.Deprecated = true
		}
	}

	// 只返回有 @Router 标签的端点
	if !hasRouter {
		return nil
	}

	return endpoint
}

// parseRouter 解析 @Router 标签
// 格式: @Router /api/users [GET]
func (cp *CommentParser) parseRouter(text string) *struct {
	Method string
	Path   string
} {
	if !strings.Contains(text, "@Router") {
		return nil
	}

	// 使用正则表达式解析
	re := regexp.MustCompile(`@Router\s+(\S+)\s+\[(\w+)\]`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 3 {
		cp.logger.Warn("@Router 标签格式错误", zap.String("text", text))
		return nil
	}

	return &struct {
		Method string
		Path   string
	}{
		Method: matches[2],
		Path:   matches[1],
	}
}

// parseSimpleTag 解析简单标签
func (cp *CommentParser) parseSimpleTag(text, tag string) string {
	if !strings.Contains(text, tag) {
		return ""
	}

	// 提取标签后的内容
	idx := strings.Index(text, tag)
	if idx == -1 {
		return ""
	}

	content := text[idx+len(tag):]
	// 移除前导空格
	content = strings.TrimSpace(content)

	return content
}

// parseParam 解析 @Param 标签
// 格式: @Param page query int false "Page number"
func (cp *CommentParser) parseParam(text string) *Parameter {
	if !strings.Contains(text, "@Param") {
		return nil
	}

	// 使用正则表达式解析
	re := regexp.MustCompile(`@Param\s+(\w+)\s+(\w+)\s+(\w+)\s+(true|false)\s+"([^"]*)"`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 6 {
		cp.logger.Debug("@Param 标签格式不完整", zap.String("text", text))
		return nil
	}

	return &Parameter{
		Name:        matches[1],
		In:          matches[2],
		Type:        matches[3],
		Required:    matches[4] == "true",
		Description: matches[5],
	}
}

// parseResponse 解析响应标签
// 格式: @Success 200 {object} User
func (cp *CommentParser) parseResponse(text, tag string) *Response {
	if !strings.Contains(text, tag) {
		return nil
	}

	// 使用正则表达式解析
	re := regexp.MustCompile(tag + `\s+(\d+)\s+{(\w+)}\s+(\S+)`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 4 {
		cp.logger.Debug("响应标签格式不完整", zap.String("text", text))
		return nil
	}

	return &Response{
		StatusCode:  matches[1],
		Description: fmt.Sprintf("%s response", tag),
		Schema: &Schema{
			Type: matches[2],
		},
	}
}

// ValidateTag 验证标签格式
func (cp *CommentParser) ValidateTag(text string) error {
	if !strings.HasPrefix(strings.TrimSpace(text), "@") {
		return fmt.Errorf("标签必须以 @ 开头")
	}

	return nil
}

// ExtractAllTags 提取所有标签
func (cp *CommentParser) ExtractAllTags(comments []string) map[string][]string {
	tags := make(map[string][]string)

	for _, comment := range comments {
		// 移除注释前缀
		text := strings.TrimPrefix(comment, "//")
		text = strings.TrimSpace(text)

		// 查找标签
		if strings.HasPrefix(text, "@") {
			parts := strings.SplitN(text, " ", 2)
			if len(parts) >= 1 {
				tagName := parts[0]
				tagValue := ""
				if len(parts) > 1 {
					tagValue = parts[1]
				}
				tags[tagName] = append(tags[tagName], tagValue)
			}
		}
	}

	return tags
}

// ParseMultilineTag 解析多行标签
func (cp *CommentParser) ParseMultilineTag(comments []string, tagName string) []string {
	var values []string

	for _, comment := range comments {
		text := strings.TrimPrefix(comment, "//")
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, tagName) {
			content := strings.TrimPrefix(text, tagName)
			content = strings.TrimSpace(content)
			if content != "" {
				values = append(values, content)
			}
		}
	}

	return values
}

// SupportedTags 返回所有支持的标签
func (cp *CommentParser) SupportedTags() []string {
	return []string{
		"@Router",
		"@Summary",
		"@Description",
		"@Tags",
		"@Param",
		"@Success",
		"@Failure",
		"@Deprecated",
	}
}
