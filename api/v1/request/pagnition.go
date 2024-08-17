package request

type Pagination struct {
	Start *int   `form:"start" binding:"required,min=0"`
	Limit *int   `form:"limit" binding:"required,min=1,max=100"`
	Query string `form:"query" binding:"max=100"`
}
