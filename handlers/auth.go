package handlers

import (
	"fiber_curd_api/app"
	"fiber_curd_api/models"
	"fiber_curd_api/tools/fmtres"
	"fiber_curd_api/tools/passwd"

	"github.com/gofiber/fiber/v3"
)

// 注册请求数据
type RegisterDTO struct {
	Email    string `json:"email"     validate:"required,email"`  // 用户邮箱
	Password string `json:"password"  validate:"required,len=32"` // 用户密码md5(避免明文传输)
	Username string `json:"username"`                             // 用户名(可选)
}

// 登录请求数据
type LoginDTO struct {
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码md5
}

// 登录响应数据
type LoginRes struct {
	User         models.User `json:"user"`          // 用户信息
	AccessToken  string      `json:"access_token"`  // 访问token
	RefreshToken string      `json:"refresh_token"` // 刷新token
}

// 续签 AccessToken 响应数据
type RenewAccessTokenRes struct {
	AccessToken  string `json:"access_token"`  // 访问token
	RefreshToken string `json:"refresh_token"` // 刷新token
}

// @Tags    auth
// @Summary 注册
// @Accept  json
// @Produce json
// @Param   data body RegisterDTO true "Register user account"
// @Success 200 {object} fmtres.FormatResponse{results=models.User}
// @Router /api/register [post]
func Register(c fiber.Ctx) error {
	// 1.获取数据 & 数据验证
	regData := new(RegisterDTO)
	if err := c.Bind().JSON(regData); err != nil {
		return err
	}

	// 2.构建数据
	password, err := passwd.PasswdEncrypt(regData.Password)
	if err != nil {
		return c.JSON(fmtres.ErrorFmt("failed to encrypt password", err))
	}
	newUser := models.User{
		Email:    regData.Email,
		Username: regData.Username,
		Password: password,
	}

	// 3.将数据插入到数据库并响应
	result := app.ConnectDB().Create(&newUser)
	if result.Error != nil {
		return c.JSON(fmtres.Error(result.Error))
	}
	return c.JSON(fmtres.OKWithResults(newUser))
}

// @Tags    auth
// @Summary 登录
// @Accept  json
// @Produce json
// @Param   data body LoginDTO true "Register user account"
// @Success 200 {object} fmtres.FormatResponse{results=LoginRes}
// @Router /api/login [post]
func Login(c fiber.Ctx) error {
	// 1.validate
	loginData := new(LoginDTO)
	if err := c.Bind().JSON(loginData); err != nil {
		return err
	}

	// 2.find record by email in database
	user := models.User{
		Email: loginData.Email,
	}
	result := app.ConnectDB().Where(user).Find(&user)
	if result.Error != nil {
		return c.JSON(fmtres.Error(result.Error))
	}

	// 3.verify password
	if !passwd.PasswdVerify(loginData.Password, user.Password) {
		return c.JSON(fmtres.ErrorStr("invalid email or password"))
	}

	// 4.generate tokens
	accessToken, err := app.JWTService.GenAccessToken(user.Email)
	if err != nil {
		return c.JSON(fmtres.ErrorFmt("failed to generate accessToken", err))
	}
	refreshToken, err2 := app.JWTService.GenRefreshToken(user.Email)
	if err2 != nil {
		return c.JSON(fmtres.ErrorFmt("failed to generate refreshToken", err2))
	}

	// 5.build response data
	loginRes := LoginRes{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return c.JSON(fmtres.OKWithResults(loginRes))
}

// @Tags    auth
// @Summary 刷新 AccessToken
// @Accept  json
// @Produce json
// @Param   refresh_token query string true "token-string"
// @Success 200 {object} fmtres.FormatResponse{results=RenewAccessTokenRes}
// @Router /api/refresh_acccess_token [get]
func RenewAccessToken(c fiber.Ctx) error {
	refreshToken := c.Query("refresh_token") // 刷新Token参数
	if refreshToken == "" {
		return c.JSON(fmtres.ErrorStr("'refresh_token' paramter not exists"))
	}

	accessToken, refreshToken, err := app.JWTService.RenewAccessToken(refreshToken)
	if err != nil {
		return c.JSON(fmtres.ErrorFmt("failed to renew accessToken", err))
	}

	res := RenewAccessTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(fmtres.OKWithResults(res))
}
