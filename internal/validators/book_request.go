package validators

type CreateBookRequest struct {
	Title       string `json:"title" binding:"required"`
	ISBN        string `json:"isbn" binding:"required"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" binding:"gte=0"`
	Available   int    `json:"available" binding:"gte=0"`
	AuthorIDs   []uint `json:"author_ids"`
}

type UpdateBookRequest struct {
	Title       string `json:"title" binding:"required"`
	ISBN        string `json:"isbn" binding:"required"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" binding:"gte=0"`
	Available   int    `json:"available" binding:"gte=0"`
	AuthorIDs   []uint `json:"author_ids"`
}
