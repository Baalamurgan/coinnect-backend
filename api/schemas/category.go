package schemas

type GetUserRequest struct {
	UserID string `json:"user_id"`
}

type CreateCategoryRequest struct {
	Name             string `gorm:"not null" json:"name"`
	Description      string `json:"description"`
	ParentCategoryID string `gorm:"uuid; default: null" json:"parent_category_id"`
}

type Item struct {
	Name        string   `json:"name"`
	Year        int      `json:"year"`
	ImageURL    string   `json:"image_url"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Details     []Detail `json:"details"`
}

type Detail struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}
