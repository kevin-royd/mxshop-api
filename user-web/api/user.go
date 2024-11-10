package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	middlewares "mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"strconv"
	"time"

	"mxshop-api/user-web/proto"
	"net/http"
)

func GetUserList(c *gin.Context) {
	// 获取偏移量和分页数量
	pn := c.DefaultQuery("pn", "0")
	// 字符串转换uint32
	pnInt, _ := strconv.Atoi(pn)

	pSize := c.DefaultQuery("pSize", "25")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := global.UserClient.GetUserList(c, &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		global.HandleGrpcErrToHttp(err, c)
		return
	}
	data := mapUserData(rsp)
	c.JSON(http.StatusOK, data)
}

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

// 通用表单校验函数
func ValidateAndCheckCaptcha(c *gin.Context, form interface{}, captchaEnabled bool) error {
	if err := c.ShouldBindJSON(form); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return fmt.Errorf(err.Error())
		}
		translatedErrs := global.RemoveTopStruct(errs.Translate(global.Translator))
		jsonErrStr, jsonErr := global.MapToJSONString(translatedErrs)
		if jsonErr != nil {
			return fmt.Errorf("转换错误信息时出错: %v", jsonErr)
		}
		return fmt.Errorf(jsonErrStr)
	}

	// 如果启用了验证码，则校验验证码
	if captchaEnabled {
		if f, ok := form.(*forms.PassWordLoginForm); ok {
			if !global.RedisStore.Verify(f.CaptchaId, f.Captcha, true) {
				return fmt.Errorf("验证码错误")
			}
		}
	}

	return nil
}

// 用户注册
func RegisterUser(c *gin.Context) {
	passWordLoginForm := forms.PassWordLoginForm{}
	err := ValidateAndCheckCaptcha(c, &passWordLoginForm, global.ServerConf.CaptchaInfo.EnableCaptcha)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// 注册用户
	rsp, err := global.UserClient.CreateUser(c, &proto.CreateUserInfo{
		Mobile:   passWordLoginForm.Mobile,
		Password: passWordLoginForm.Password,
	})
	if err != nil {
		global.HandleGrpcErrToHttp(err, c)
		return
	}

	data, err := CreateUserToken(c, rsp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "注册成功",
		"data": data,
	})
}

// 校验用户登录
func PassWordLoginForms(c *gin.Context) {
	passWordLoginForm := forms.PassWordLoginForm{}
	err := ValidateAndCheckCaptcha(c, &passWordLoginForm, global.ServerConf.CaptchaInfo.EnableCaptcha)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// 查询手机号是否存在
	rsp, err := global.UserClient.GetUserByMobile(c, &proto.MobileRequest{
		Mobile: passWordLoginForm.Mobile,
	})
	if err != nil {
		global.HandleGrpcErrToHttp(err, c)
		return
	}

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
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "登陆成功", "data": data})
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
