package book_service

type QueryArgs struct {
	BookName string `json:"book_name,omitempty"`
	Id       string `json:"id,omitempty"`
}
