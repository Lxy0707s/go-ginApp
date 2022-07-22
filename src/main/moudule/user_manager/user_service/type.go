package user_service

type QueryArgs struct {
	UserName string `json:"user_name,omitempty"`
	Email    string `json:"email,omitempty"`
}

type UpdateArgs struct {
	UserName string `json:"book_name,omitempty"`
	Email    string `json:"email,omitempty"`
}

type DeleteArgs struct {
	UserName string `json:"book_name,omitempty"`
	Email    string `json:"email,omitempty"`
}
