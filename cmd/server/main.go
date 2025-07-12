package main

import (
	"github.com/coderconquerer/go-login-app/cmd"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
