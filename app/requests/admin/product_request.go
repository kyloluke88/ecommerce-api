package admin

type GetProductFrom1688 struct {
	ProductUrl string `json:"product_url" validate:"required"`
}

type CreateProduct struct {
	LoginId  string `json:"login_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateProduct struct {
	LoginId  string `json:"login_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}
