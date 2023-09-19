package entity

type Person struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" validate:"required,alpha,startsWithUpperCase" example:"Ivan"`
	Surname     string `json:"surname" validate:"required,alpha,startsWithUpperCase" example:"Ivanov"`
	Patronymic  string `json:"patronymic,omitempty" example:"Sergeevich"`
	Age         int    `json:"age" validate:"gte=0,lte=120" example:"70"`
	Gender      string `json:"gender" validate:"oneof=male female" example:"male"`
	Nationality string `json:"nationality" validate:"alpha" example:"RU"`
}
