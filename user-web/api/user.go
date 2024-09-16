package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	middlewares "mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"strconv"
	"time"

	"mxshop-api/user-web/proto"
	"net/http"
)

func GetUserList(ctx *gin.Context) {
	// 获取偏移量和分页数
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)

	pSize := ctx.DefaultQuery("pSize", "25")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := global.UserClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		global.HandleGrpcErrToHttp(err, ctx)
		return
	}
	data := mapUserData(rsp)
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// 封装返回对象
func mapUserData(rsp *proto.UserListResponse) []interface{} {
	u := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := make(map[string]interface{}, 0)
		user["id"] = value.Id
		user["mobile"] = value.Mobile
		user["nickname"] = value.Nickname
		user["gender"] = value.Gender
		user["role"] = value.Role
		u = append(u, user)
	}
	return u
}

// 校验用户登录
func PassWordLoginForms(c *gin.Context) {
	// 实例化表单
	passWordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBindJSON(&passWordLoginForm); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": global.RemoveTopStruct(errs.Translate(global.Translator)),
		})
		return
	}
	// 校验验证码
	if global.ServerConf.CaptchaInfo.EnableCaptcha {
		if !global.RedisStore.Verify(passWordLoginForm.CaptchaId, passWordLoginForm.Captcha, true) {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "验证码错误",
			})
			return
		}
	}
	// 查询手机号是否存在
	if rsp, err := global.UserClient.GetUserByMobile(c, &proto.MobileRequest{
		Mobile: passWordLoginForm.Mobile,
	}); err != nil {
		global.HandleGrpcErrToHttp(err, c)
		return

	} else {
		// 校验密码
		verify, err := global.UserClient.CheckUserPasswd(c, &proto.PasswordCheckInfo{
			Password:          passWordLoginForm.Password,
			EncryptedPassword: rsp.Password,
		})
		if err != nil {
			global.HandleGrpcErrToHttp(err, c)
			return
		}
		if !verify.Success {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "密码错误"})
			return
		}
		data, err := CreateUserToken(c, rsp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"msg": "登陆成功", "data": data})
	}

}

// 创建用户token
func CreateUserToken(c *gin.Context, rsp *proto.UserInfoResponse) (map[string]interface{}, error) {
	// 实例化jwt对象
	j := middlewares.NewJWT()
	// 实例化claims
	claims := models.CustomClaims{
		ID:          uint(rsp.Id),
		NickName:    rsp.Nickname,
		AuthorityId: uint(rsp.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    "test",
		},
	}
	// 创建token
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("创建token失败")
	}
	data := map[string]interface{}{
		"Token":     token,
		"ExpiresAt": time.Now().Unix() + 60*60,
	}
	return data, nil
}
