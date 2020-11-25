package user

import (
	"html"
	"net/http"

	"github.com/linuxxiaoyu/blog/pkg/setting"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	name := html.EscapeString(c.PostForm("name"))
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

	var user User
	db := setting.DB
	result := db.Where("name = ?", name).First(&user)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	user.Name = name
	user.Password = addSalt(password)
	result = db.Create(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotImplemented, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.ID,
	})
}
