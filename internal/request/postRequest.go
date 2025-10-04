package request

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// form query参数或form表单，get请求
// json json正文body，post、put请求
type PostListRequest struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
}

// 初始化时设置默认值
func (r *PostListRequest) SetDefault() {
	if r.PageNum <= 0 {
		r.PageNum = 1 // 没传或传了 <=0 的值，用默认 1
	}
	if r.PageSize <= 0 || r.PageSize > 100 {
		r.PageSize = 10 // 没传或超出范围，用默认 10
	}
}
