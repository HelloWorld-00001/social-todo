package main

import (
	"github.com/coderconquerer/social-todo/cmd"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
