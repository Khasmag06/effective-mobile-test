package entity

type Person struct {
	ID          int    `json:"id"`
	Name        string `json:"name" validate:"required,alpha,startsWithUpperCase"`
	Surname     string `json:"surname" validate:"required,alpha,startsWithUpperCase"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age" validate:"gte=0,lte=120"`
	Gender      string `json:"gender" validate:"oneof=male female"`
	Nationality string `json:"nationality" validate:"alpha"`
}
