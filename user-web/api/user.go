package api

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"strconv"

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
	ctx.JSON(http.StatusOK, gin.H{"data": rsp.Data})
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
	// 查询手机号是否存在
	if _, err := global.UserClient.GetUserByMobile(c, &proto.MobileRequest{
		Mobile: passWordLoginForm.Mobile,
	}); err != nil {
		global.HandleGrpcErrToHttp(err, c)
		return
	} else {
		// 校验密码
		options := &password.Options{16, 100, 32, sha512.New}
		salt, encodedPwd := password.Encode(passWordLoginForm.Password, options)
		pwd := fmt.Sprintf("$sha512$%s$%s", salt, encodedPwd)
		fmt.Printf("pwd %s\n", pwd)
		verify, _ := global.UserClient.CheckUserPassword(c, &proto.PasswordCheckInfo{
			Password:          passWordLoginForm.Password,
			EncryptedPassword: pwd,
		})

		fmt.Printf("%+v\n", verify)
	}

}
