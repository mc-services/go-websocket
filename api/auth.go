package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gp-websoket/database"
	"gp-websoket/middleware"
	"gp-websoket/model"
	"log"
	"net/http"
	"time"
)

// Login 登录
func Login(c *gin.Context) {
	var form model.User
	var user model.User

	if c.BindJSON(&form) == nil {
		res := database.DB.Where("name = ?", form.Name).First(&user)

		if res.RecordNotFound() {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"exception": "json 解析失败",
			})
			return
		}

		fmt.Println(23333)
		fmt.Println(form)
		if form.Password != user.Password {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"exception": "验证失败",
			})
			return
		}

		generateToken(c, user)
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"exception": "json 解析失败",
		})
	}
}

// 生成令牌
func generateToken(c *gin.Context, user model.User) {
	j := &middleware.JWT{
		[]byte("sign"),
	}

	claims := middleware.CustomClaims{
		user.Id,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
			Issuer:    "sign",                       //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}

	log.Println(token)

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "登录成功！",
		"token":   token,
	})
	return
}
