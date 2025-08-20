package ginadp

type UsersQuery struct {
	Search  string `form:"search"`
	Page    int    `form:"page,default=1"`
	PerPage int    `form:"per_page,default=20"`
}

type UsersListItem struct {
	ID        int64  `json:"id"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type UsersListResponse struct {
	Data []UsersListItem `json:"data"`
	Meta PaginationMeta  `json:"meta"`
}
