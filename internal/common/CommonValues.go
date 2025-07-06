package common

var (
	InvalidID = -1
)

const (
	AccountContextKey = "currentAccount"
)

type Role int

const (
	AdminRole   Role = iota // 0
	DbAdminRole             // 1
	UserRole                // 2
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
