package forms

type PassWordLoginForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,len=11,mobile"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}
