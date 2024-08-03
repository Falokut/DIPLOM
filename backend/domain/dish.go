package domain

type Dish struct {
	Id          int32
	Name        string
	Description string
	Price       int32
	Url         string   `json:",omitempty"`
	Categories  []string `json:",omitempty"`
}

type AddDishRequest struct {
	Name        string `validate:"min=0"`
	Description string
	Price       int32   `validate:"gt=0"`
	Categories  []int32 `json:",omitempty"`
	Image       []byte  `json:",omitempty"`
}
