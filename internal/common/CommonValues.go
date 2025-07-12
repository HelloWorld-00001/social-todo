package common

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
	DbMainName = "mysql"
)

const (
	User = "User"
	Todo = "Todo"
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
		return User
	case TodoEntity:
		return Todo
	default:
		return "InvalidEntity"
	}
}

func EntityFromString(s string) Entity {
	switch s {
	case User:
		return UserEntity
	case Todo:
		return TodoEntity
	default:
		return InvalidEntity
	}
}
