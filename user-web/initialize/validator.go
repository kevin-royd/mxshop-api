package initialize

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	"mxshop-api/user-web/global"
	myValidator "mxshop-api/user-web/validator"
)

func InitValidator(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok { //更换gin的验证引擎
		//将struct字段转为tag的json字段
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() //新建中文翻译
		enT := en.New()
		uni := ut.New(enT, zhT, enT)                      //存储在容器中
		global.Translator, ok = uni.GetTranslator(locale) //获取一个翻译器
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}
		switch locale {
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, global.Translator) //注册中文翻译器
			break
		case "en":
			en_translations.RegisterDefaultTranslations(v, global.Translator)
			break
		default:
			zh_translations.RegisterDefaultTranslations(v, global.Translator) //注册中文翻译器
			break
		}
		//注册自定义验证方法
		initCustomizeValidator(v)
		return nil
	}
	return nil
}

//自定义验证器初始化
func initCustomizeValidator(v *validator.Validate) {
	checkMobile(v)
}

//注册mobile验证器
func checkMobile(v *validator.Validate) {
	v.RegisterValidation("mobile", myValidator.ValidateMobile)
	v.RegisterTranslation("mobile", global.Translator, func(ut ut.Translator) error {
		return ut.Add("mobile", "{0}手机号格式不正确!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile", fe.Field())

		return t
	})
}
