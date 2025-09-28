package storage

import (
"github.com/coderconquerer/social-todo/module/user/entity"
"github.com/gin-gonic/gin"
)

func (db *MySQLConnection) CreateUser(c *gin.Context, acc *entity.User) error {
// filter deleted first
dbc := db.conn.Table(entity.User{}.TableName())

if err := dbc.Create(acc).Error; err != nil {
return err
}

return nil
}
