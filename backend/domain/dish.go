package domain

type Dish struct {
	Id          int32
	Name        string
	Description string
	Price       int32
	Url         string   `json:",omitempty"`
	Categories  []string `json:",omitempty"`
}

type GetDishesRequest struct {
	Ids    string `query:"ids"`
	Limit  int32  `query:"limit"`
	Offset int32  `query:"offset"`
}

type AddDishRequest struct {
	Name        string `validate:"required,min=1"`
	Description string
	Price       int32   `validate:"gt=0"`
	Categories  []int32 `json:",omitempty"`
	Image       []byte  `json:",omitempty"`
}
