package entity

type Dish struct {
	Id             int32
	Name           string
	Description    string
	Price          int32
	ImageId        string
	Categories     string
	RestaurantName string
}

type AddDishRequest struct {
	Name         string
	Description  string
	ImageId      string
	Price        int32
	Categories   []int32
	RestaurantId int32
}
type EditDishRequest struct {
	Id           int32
	Name         string
	Description  string
	ImageId      string
	Price        int32
	Categories   []int32
	RestaurantId int32
}
