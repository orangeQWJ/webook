package web

import (
	"net/http"
	"xws/webook/internal/domain"
	"xws/webook/internal/service"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

// UserHandler 我准备在它上面定义跟用户有关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp2.Regexp
	passwordExp *regexp2.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
	)

	emailExp := regexp2.MustCompile(emailRegexPattern, regexp2.None)
	passwordExp := regexp2.MustCompile(passwordRegexPattern, regexp2.None)

	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	// 解析错了，就会直接写回一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "解析错误")
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8为,包含数字,特殊字符")
		return
	}

	// 调用 scv方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	//fmt.Printf("%v", req)
	ctx.String(http.StatusOK, "注册成功")

}
func (u *UserHandler) Login(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hahah")
}
func (u *UserHandler) Edit(ctx *gin.Context)    {}
func (u *UserHandler) Profile(ctx *gin.Context) {}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}
