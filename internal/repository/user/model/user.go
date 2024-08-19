package userrepositorymodel

type User struct {
	ID           uint32
	Role         string
	Email        string
	HashPassword string
}
