package user

import (
	"net/http"

	"github.com/linuxxiaoyu/blog/pkg/auth"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/pkg/setting"
)

// Login a user, return a token
// GET /user
// form: name password
func Login(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")

	valid := validation.Validation{}

	valid.Required(name, "name").Message("name can't empty")
	valid.MaxSize(name, 16, "name").Message("name's max length is 16")

	valid.Required(password, "password").Message("password can't empty")
	valid.MinSize(password, 6, "password").Message("password's min length is 6")
	valid.MaxSize(password, 16, "password").Message("password's max length is 16")

	if valid.HasErrors() {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	user := User{
		Name: name,
	}
	db := setting.DB
	result := db.Where("name = ? and password = ?", name, addSalt(password)).First(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, nil)
		return
	}
	token, err := auth.GenerateToken(user.ID, name)
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}