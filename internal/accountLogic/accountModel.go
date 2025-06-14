package accountLogic

type Account struct {
	ID           int    `json:"id" example:"1"`
	Username     string `json:"username" example:"johndoe"`
	HashPassword string `json:"hash_password" example:"$2a$10$..."`
	Salt         string `json:"salt" example:"a1b2c3d4"`
}
