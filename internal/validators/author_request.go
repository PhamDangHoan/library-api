package validators

type CreateAuthorRequest struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}

type UpdateAuthorRequest struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}
