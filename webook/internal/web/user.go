package web

import (
	"fmt"
	"net/http"
	"time"
	"xws/webook/internal/domain"
	"xws/webook/internal/service"

	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

// UserHandler 我准备在它上面定义跟用户有关的路由
type UserHandler struct {
	svc         *service.UserService
	codeSvc     *service.CodeService
	emailExp    *regexp2.Regexp // 编译好的正则表达式
	passwordExp *regexp2.Regexp
	birthdayExp *regexp2.Regexp
}

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
		birthdayRegexPattern = `^\d{4}-\d{2}-\d{2}$`
	)

	emailExp := regexp2.MustCompile(emailRegexPattern, regexp2.None)
	passwordExp := regexp2.MustCompile(passwordRegexPattern, regexp2.None)
	birthdayExp := regexp2.MustCompile(birthdayRegexPattern, regexp2.None)

	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		birthdayExp: birthdayExp,
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
		// 前端的问题,前端传过来的应该是json格式
		ctx.String(http.StatusOK, "解析错误")
		return
	}

	// 检验邮箱格式
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "正则匹配超时")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}

	// 检查两次密码是否一致
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}

	// 密码强度是否符合要求
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "正则匹配超时")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8为,包含数字,特殊字符")
		return
	}

	// 调用 scv方法, 尝试注册新用户
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
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "解析错误")
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil { //数据库错误
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 从当前请求上下文ctx中获取默认的会话对象
	// Gin 框架中每个ctx上下文都有一份会化数据.
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		//Secure: true,
		//Path: "/users/edit",
		MaxAge: 60 * 5,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "解析错误")
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil { //数据未知错误
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 用户账户密码正确
	// 在这里用JWT 设置登录态
	// 生成一个token
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}

	fmt.Println(user)
	//token := jwt.New(jwt.SigningMethodHS512)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		fmt.Println("jwt 系统错误")
		fmt.Println(err)
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	//fmt.Println(tokenStr)
	ctx.Header("x-jwt-token", tokenStr)
	//ctx.String(http.StatusOK, "登录成功:%s", tokenStr)
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "退出登录")
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	//{nickname: "qwj", birthday: "2024-06-12", aboutMe: "NB"}
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "解析错误")
		return
	}
	fmt.Println(req)
	ok, err := u.birthdayExp.MatchString(req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "正则匹配超时")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}
	if len(req.Nickname) > 16 {
		ctx.String(http.StatusOK, "昵称长度不得超过16")
		return
	}
	if len(req.AboutMe) > 50 {
		ctx.String(http.StatusOK, "简介长度不得超过50")
		return
	}
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	userIdInt, ok := userId.(int64)
	if !ok {
		ctx.String(http.StatusOK, "未登录")
		return
	}
	userInfo := domain.User{
		Id:       userIdInt,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe,
	}
	err = u.svc.EditProfile(ctx, userInfo)
	if err != nil {
		ctx.String(http.StatusOK, "更新资料时出错")
		return
	}
	ctx.String(http.StatusOK, "更新profile成功")
	return
}

func (u *UserHandler) EditJWT(ctx *gin.Context) {
	//{nickname: "qwj", birthday: "2024-06-12", aboutMe: "NB"}
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "解析错误")
		return
	}
	fmt.Println(req)
	ok, err := u.birthdayExp.MatchString(req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "正则匹配超时")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}
	if len(req.Nickname) > 16 {
		ctx.String(http.StatusOK, "昵称长度不得超过16")
		return
	}
	if len(req.AboutMe) > 50 {
		ctx.String(http.StatusOK, "简介长度不得超过50")
		return
	}
	c, exists := ctx.Get("claims")
	if !exists {
		//	预期之外,添加监控
		ctx.String(http.StatusOK, "%s", "系统错误")
		return
	}
	claims, ok := c.(*UserClaims)
	if !ok {
		//	预期之外,添加监控
		ctx.String(http.StatusOK, "%s", "系统错误")
		return
	}
	userInfo := domain.User{
		Id:       claims.Uid,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe,
	}
	err = u.svc.EditProfile(ctx, userInfo)
	if err != nil {
		ctx.String(http.StatusOK, "更新资料时出错")
		return
	}
	ctx.String(http.StatusOK, "更新profile成功")
	return
}
func (u *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	userIdInt, ok := userId.(int64)
	if !ok {
		ctx.String(http.StatusOK, "未登录")
		return
	}
	userInfo, err := u.svc.ShowProfile(ctx, userIdInt)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	ctx.String(http.StatusOK, "Nickname: %s", userInfo)
	return
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	/*
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &UserClaims{}
		jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		claims := claims.Uid
	*/
	c, exists := ctx.Get("claims")
	if !exists {
		//	预期之外,添加监控
		ctx.String(http.StatusOK, "%s", "系统错误")
		return
	}
	claims, ok := c.(*UserClaims)
	if !ok {
		//	预期之外,添加监控
		ctx.String(http.StatusOK, "%s", "系统错误")
		return
	}
	userInfo, err := u.svc.ShowProfile(ctx, claims.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "mysql/redis错误")
		return
	}
	ctx.String(http.StatusOK, "%s", userInfo)
	return
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.EditJWT)
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	const biz = "login"
	if err := ctx.Bind(&req); err != nil {
		// 前端的问题,前端传过来的应该是json格式
		ctx.String(http.StatusOK, "解析错误")
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码发送成功",
		})

	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码请求太频繁,请稍后再试",
		})
	default:
		// cache.ErrEnryWithoutExpire
		// cache.ErrUnknowForLuaScript
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})

	}
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	const biz = "login"
	if err := ctx.Bind(&req); err != nil {
		// 前端的问题,前端传过来的应该是json格式
		ctx.String(http.StatusOK, "解析错误")
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入正确的手机号",
		})
	}
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err == service.ErrCodeExpired {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证失败",
		})
		return
	}
	if err != nil {
		//cache.ErrUnknowForLuaScript
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证成功",
		})
	} else {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误",
		})

	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己要放进token里的数据
	Uid int64
	// 自己随便加
	// 敏感信息不要加
	UserAgent string
}
