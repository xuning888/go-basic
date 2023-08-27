package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-basic/webbook/internal/domain"
	"go-basic/webbook/internal/service"
	"go-basic/webbook/internal/util"
	"net/http"
	"time"
)

const (
	// 校验邮箱格式的正则表达式
	emailRegexPattern = `^[\w\.-]+@[a-zA-Z\d\.-]+\.[a-zA-Z]{2,}$`
	// 校验密码格式的正则表达式
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	// 校验用户的昵称为长度为1 - 20
	nicknameRegexPattern = `^.{1,20}$`
	// 生日的格式校验
	birthdayRegexPattern = `^\d{4}-\d{2}-\d{2}$`
	// 自我介绍不能超过 500 个字符
	introductionRegexPattern = `^.{1,500}$`
)

// UserHandler 定义所有跟用户有关的路由
// 为了分组方便， 为了依赖注入
type UserHandler struct {
	emailExp        *regexp.Regexp
	passwordExp     *regexp.Regexp
	nicknameExp     *regexp.Regexp
	birthdayExp     *regexp.Regexp
	introductionExp *regexp.Regexp
	svc             *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	// 将正则表达式的预编译放到 UserHandler的初始化中
	return &UserHandler{
		emailExp:        regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp:     regexp.MustCompile(passwordRegexPattern, regexp.None),
		nicknameExp:     regexp.MustCompile(nicknameRegexPattern, regexp.None),
		birthdayExp:     regexp.MustCompile(birthdayRegexPattern, regexp.None),
		introductionExp: regexp.MustCompile(introductionRegexPattern, regexp.None),
		svc:             svc,
	}
}

// SignUp 注册用户信息
func (u *UserHandler) SignUp(ctx *gin.Context) {

	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"Password"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里边
	// 解析错了， 就会直接会写一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不相同")
		return
	}

	isPassword, err := u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含数字、特殊字符，并且长度不能小于 8 位")
		return
	}

	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicate {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
	return
}

// Login 用户登录
func (u *UserHandler) Login(ctx *gin.Context) {

	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	sess := sessions.Default(ctx)
	// Secure 表示启用 HTTP, HttpOnly: 表示
	sess.Options(sessions.Options{
		MaxAge: 60,
	})
	sess.Set("userId", user.Id)
	err = sess.Save()
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	ctx.String(http.StatusOK, "登录成功")

	return
}

// LoginJWT 用户登录
func (u *UserHandler) LoginJWT(ctx *gin.Context) {

	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	sess := sessions.Default(ctx)
	// Secure 表示启用 HTTP, HttpOnly: 表示
	sess.Options(sessions.Options{
		MaxAge: 60,
	})
	sess.Set("userId", user.Id)
	err = sess.Save()
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}

	claims := UserClaims{
		// 设置过期时间
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid: user.Id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("MxQP9pSI6BzUL9XVSZrdSeJm6Jbhw42z"))

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "登录成功")
	return
}

// Edit 用户这是信息
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Id           int64  `json:"id"`
		Nickname     string `json:"nickname"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}

	var req EditReq

	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	sess := sessions.Default(ctx)

	userId, _ := sess.Get("userId").(int64)

	if userId == 0 {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	req.Id = userId

	if util.IsNotBlank(req.Nickname) {
		validNickname, err := u.nicknameExp.MatchString(req.Nickname)
		if err != nil {
			ctx.String(http.StatusOK, "系统异常")
			return
		}
		if !validNickname {
			ctx.String(http.StatusOK, "昵称长度为1-20")
			return
		}
	} else {
		ctx.String(http.StatusOK, "昵称不能为空")
		return
	}

	if util.IsNotBlank(req.Birthday) {
		validBirthday, err := u.birthdayExp.MatchString(req.Birthday)
		if err != nil {
			ctx.String(http.StatusOK, "系统异常")
			return
		}
		if !validBirthday {
			ctx.String(http.StatusOK, "生日格式错误")
			return
		}
	}

	if util.IsNotBlank(req.Introduction) {
		validIntroduction, err := u.introductionExp.MatchString(req.Introduction)
		if err != nil {
			ctx.String(http.StatusOK, "系统异常")
			return
		}
		if !validIntroduction {
			ctx.String(http.StatusOK, "自我介绍不能超过500个字")
			return
		}
	}

	err := u.svc.Edit(ctx, domain.User{
		Id:           req.Id,
		NickName:     req.Nickname,
		Birthday:     req.Birthday,
		Introduction: req.Introduction,
	})
	if err == service.ErrUserNotFound {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "修改成功")
	return
}

// Profile 返回所有用户信息
func (u *UserHandler) Profile(ctx *gin.Context) {

	sess := sessions.Default(ctx)

	userId, ok := sess.Get("userId").(int64)

	if !ok {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	profile, err := u.svc.Profile(ctx, userId)

	if err == service.ErrUserNotFound {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	type ProfileRes struct {
		Id           int64  `json:"id"`
		Email        string `json:"email"`
		Nickname     string `json:"nickname"`
		Introduction string `json:"introduction"`
		Birthday     string `json:"birthday"`
	}

	res := ProfileRes{
		Id:           profile.Id,
		Email:        profile.Email,
		Nickname:     profile.NickName,
		Introduction: profile.Introduction,
		Birthday:     profile.Birthday,
	}

	ctx.JSON(http.StatusOK, res)
	return
}

// ProfileJWT 返回所有用户信息
func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")

	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	userId := claims.Uid
	profile, err := u.svc.Profile(ctx, userId)
	if err == service.ErrUserNotFound {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	type ProfileRes struct {
		Id           int64  `json:"id"`
		Email        string `json:"email"`
		Nickname     string `json:"nickname"`
		Introduction string `json:"introduction"`
		Birthday     string `json:"birthday"`
	}

	res := ProfileRes{
		Id:           profile.Id,
		Email:        profile.Email,
		Nickname:     profile.NickName,
		Introduction: profile.Introduction,
		Birthday:     profile.Birthday,
	}
	ctx.JSON(http.StatusOK, res)
	return
}

// RegisterRoutes 注册路由以及路由对应的处理方法
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.ProfileJWT)
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明我自己的要放到 token中的数据
	Uid int64
	// 自己随便加，但是不能放敏感信息
}
