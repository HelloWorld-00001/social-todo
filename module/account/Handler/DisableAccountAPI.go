package Handler

import (
	"errors"
	common2 "github.com/coderconquerer/go-login-app/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (ah *AccountHandler) DisableAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract query parameters
		id := c.Query("id")           // returns string (empty if missing)
		disable := c.Query("disable") // returns string

		// Optional: convert to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, common2.NewInvalidInputError(err))
			return
		}
		disableInt, err := strconv.Atoi(disable)
		if err != nil {
			c.JSON(http.StatusBadRequest, common2.NewInvalidInputError(err))
			return
		}
		if disableInt != 0 && disableInt != 1 {
			c.JSON(http.StatusBadRequest, common2.NewInvalidInputError(errors.New("disable must be 0 or 1")))
			return
		}

		if err := ah.DisableLogic.DisableAccount(c, idInt, disableInt == 1); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, common2.SimpleResponse(true))
	}
}
