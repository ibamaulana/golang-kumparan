package news

type CreateRequest struct {
	Author string `json:"author" binding:"required"`
	Body   string `json:"body"  binding:"required"`
}

type FindRequest struct {
	ID int `json:"id" binding:"required"`
}
