package controller

import (
	"aion/model"
	"aion/service"
	"aion/zlog"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"strings"
)

type userController struct {
	Controller
	userService service.UserService
}

var UserController = &userController{
	userService: service.NewUser(),
}

// Register 用户注册
func (uc userController) Register(ctx *gin.Context) {
	name := ctx.PostForm("name")
	if !govalidator.StringLength(name, "1", "10") {
		uc.Failed(ctx, ParamError, "名称长度不正确1-10")
		return
	}
	email := ctx.PostForm("email")
	if !govalidator.IsEmail(email) {
		uc.Failed(ctx, ParamError, "邮箱不正确")
		return
	}
	password := ctx.PostForm("password")
	if !govalidator.StringLength(password, "6", "16") {
		uc.Failed(ctx, ParamError, "密码长度不正确6-16")
		return
	}
	token, err := uc.userService.Register(name, email, password)
	if err != nil {
		zlog.WithContext(ctx).Sugar().Errorf("uc register Failed, error: %s", err.Error())
		if _, ok := err.(service.Error); ok {
			uc.Failed(ctx, ParamError, err.Error())
		} else {
			uc.Failed(ctx, Failed, "注册失败")
		}
	} else {
		zlog.WithContext(ctx).Sugar().Infof("register uc Success, email: %s", email)
		uc.Success(ctx, "ok", gin.H{"token": token})
	}
	return
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login 用户登录
func (uc userController) Login(ctx *gin.Context) {
	var request LoginRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		uc.Failed(ctx, ParamError, err.Error())
		return
	}
	if !govalidator.IsEmail(request.Email) {
		uc.Failed(ctx, ParamError, "邮箱不正确")
		return
	}
	if !govalidator.StringLength(request.Password, "6", "16") {
		uc.Failed(ctx, ParamError, "密码长度不正确6-16")
		return
	}
	token, err := uc.userService.Login(request.Email, request.Password)
	if err != nil {
		zlog.WithContext(ctx).Sugar().Errorf("register failed, error: %s", err.Error())
		if errors.As(err, &service.Error{}) {
			uc.Failed(ctx, Failed, err.Error())
		}
	} else {
		zlog.WithContext(ctx).Sugar().Infof("login success, email: %s", request.Email)
		uc.Success(ctx, "ok", gin.H{"token": "Bearer " + token, "currentAuthority": "admin"})
	}
	return
}

// Current 用户登录
func (uc userController) Current(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		uc.Success(ctx, "ok", gin.H{"name": "游客"})
		return
	}
	userId, err := uc.userService.ParseToken(strings.TrimSpace(strings.Trim(token, "Bearer")))
	if err != nil {
		uc.Success(ctx, "ok", gin.H{"name": "游客"})
	} else {
		user, _ := model.User{}.GetUser(int(userId))
		uc.Success(ctx, "ok", gin.H{"name": user.Name, "currentAuthority": "admin"})
	}
	return
}

// Logout 用户退出
func (uc userController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	zlog.WithContext(ctx).Sugar().Debugf("add token into blacklist, token: %s", token)

	uc.Success(ctx, "ok", "")
}
