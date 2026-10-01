package validators

type BookQuery struct {
	Page   int
	Limit  int
	Search string
	ISBN   string
	Sort   string
	Order  string
}
