package form

type SignForm struct {
	Keyword  string `json:"keyword" binding:"required"`
	SignType int    `json:"sign_type" binding:"required,min=1,max=5"`
}
