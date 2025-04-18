package schema

type InternalCreatePostRequest struct {
	Title    string `json:"title" validate:"required,max=200"`
	Content  string `json:"content" validate:"required"`
	Category string `json:"category" validate:"required,max=100"`
	Status   string `json:"status" validate:"omitempty,oneof=Draft Publish Trash"`
}

type InternalUpdatePostRequest struct {
	ID       string `json:"id" validate:"required,uuid"`
	Title    string `json:"title" validate:"required,max=200"`
	Content  string `json:"content" validate:"required"`
	Category string `json:"category" validate:"required,max=100"`
	Status   string `json:"status" validate:"omitempty,oneof=Draft Publish Trash"`
}
