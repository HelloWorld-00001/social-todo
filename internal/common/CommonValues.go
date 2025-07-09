package common

import (
	models2 "github.com/coderconquerer/go-login-app/internal/TodoItem/models"
	"github.com/coderconquerer/go-login-app/internal/user/models"
)

var (
	InvalidID = -1
)

const (
	AccountContextKey = "currentAccount"
)

type Role int
type Entity int

const (
	AdminRole   Role = iota // 0
	DbAdminRole             // 1
	UserRole                // 2
)

const (
	UserEntity Entity = iota
	TodoEntity
	InvalidEntity
)

func (r Role) ToString() string {
	switch r {
	case AdminRole:
		return "Admin"
	case DbAdminRole:
		return "DbAdmin"
	default:
		return "User"
	}
}

func (e Entity) ToString() string {
	switch e {
	case UserEntity:
		return models.User{}.TableName()
	case TodoEntity:
		return models2.Todo{}.TableName()
	default:
		return "InvalidEntity"
	}
}

func EntityFromString(s string) Entity {
	switch s {
	case models.User{}.TableName():
		return UserEntity
	case models2.Todo{}.TableName():
		return TodoEntity
	default:
		return InvalidEntity
	}
}
