package common

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
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

type React int

const (
	ReactionLike React = iota
	ReactionDislike
	ReactionUnreact
	ReactionLove
	ReactionWow
	ReactionAngry
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

var reactToString = map[React]string{
	ReactionLike:    "like",
	ReactionDislike: "dislike",
	ReactionUnreact: "unreact",
	ReactionLove:    "love",
	ReactionWow:     "wow",
	ReactionAngry:   "angry",
}

var stringToReact = map[string]React{
	"like":    ReactionLike,
	"dislike": ReactionDislike,
	"unreact": ReactionUnreact,
	"love":    ReactionLove,
	"wow":     ReactionWow,
	"angry":   ReactionWow,
}

// String returns the string representation of the React enum
func (r React) String() string {
	if val, ok := reactToString[r]; ok {
		return val
	}
	return ""
}

// GetUidFromString parses a string to its corresponding React enum
func GetReactionFromString(s string) (React, error) {
	if val, ok := stringToReact[strings.ToLower(s)]; ok {
		return val, nil
	}
	return React(0), errors.New("invalid react string: " + s)
}

// Scanner: convert from DB (string) to React
func (r *React) Scan(value interface{}) error {
	var strVal string

	switch v := value.(type) {
	case string:
		strVal = v
	case []byte:
		strVal = string(v)
	default:
		return fmt.Errorf("React.Scan: expected string or []byte, got %T", value)
	}

	strVal = strings.ToLower(strVal)
	if val, ok := stringToReact[strVal]; ok {
		*r = val
		return nil
	}

	return fmt.Errorf("React.Scan: invalid value %s", strVal)
}

// Valuer: convert from React to DB (string)
func (r React) Value() (driver.Value, error) {
	str := r.String()
	if str == "" {
		return nil, errors.New("React.Value: invalid enum value")
	}
	return str, nil
}
