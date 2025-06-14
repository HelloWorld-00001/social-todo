package accountLogic

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
)

const (
	loginCount    = "login"
	registerCount = "register"
	pingCount     = "ping"
)

var calledApi = make(map[string]int)

type AccountController struct {
	service *AccountService
}

func NewController(service *AccountService) *AccountController {
	return &AccountController{service}
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Register a new user with a username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body  AccountDto true  "User credentials"
// @Success      200   {object} map[string]string
// @Failure      400   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /register [post]
func (ctl *AccountController) RegisterUser(c *gin.Context) {
	var req AccountDto
	calledApi[registerCount] += 1
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctl.service.Register(req.Username, req.Password)
	if err != nil {
		var errMessage = err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed because: " + errMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

// LoginUser godoc
// @Summary      Login a user
// @Description  Authenticate a user with a username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body  AccountDto  true  "User credentials"
// @Success      200   {object} map[string]string
// @Failure      400   {object} map[string]string
// @Failure      401   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /login [post]
func (ctl *AccountController) LoginUser(c *gin.Context) {
	var req AccountDto
	calledApi[loginCount] += 1

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctl.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if user != nil {
		data, _ := json.Marshal(user)
		session := sessions.Default(c)
		session.Set("user", string(data))
		err = session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create session"})
		}
		c.JSON(http.StatusOK, gin.H{"message": "login successful"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
	}
}

// GetAPICount godoc
// @Summary      Get api count
// @Description  Only for authorized users
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200   {object} map[string]string
// @Failure      400   {object} map[string]string
// @Failure      401   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /api_count [GET]
func (ctl *AccountController) GetAPICount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "API usage stats",
		"calledApi": calledApi})
}

// PingMe godoc
// @Summary      Ping API
// @Description  Only for authorized users
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200   {object} map[string]string
// @Failure      400   {object} map[string]string
// @Failure      401   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /ping [GET]
func (ctl *AccountController) PingMe(c *gin.Context) {

	calledApi[pingCount] += 1
	c.JSON(http.StatusOK, gin.H{"message": "Hello, you pinged me"})
}
