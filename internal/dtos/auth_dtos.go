package dtos

type RegisterUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"firstName" validate:"required,min=2,max=50"`
	LastName    string `json:"lastName" validate:"required,min=2,max=50"`
	CompanyName string `json:"companyName" validate:"required,min=2,max=100"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
}

type UserResponse struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	CompanyName  string `json:"companyName"`
	IsActive     bool   `json:"isActive"`
	TwoFAEnabled bool   `json:"twoFaEnabled"`
}

type RegisterUserResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse struct {
	User                 UserResponse `json:"user"`
	AccessToken          string       `json:"accessToken"`
	RefreshToken         string       `json:"refreshToken"`
	Message              string       `json:"message"`
	TwoFARequired        bool         `json:"twoFaRequired"`
	AccessTokenExpiresAt int64        `json:"accessTokenExpiresAt"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken          string `json:"accessToken"`
	RefreshToken         string `json:"refreshToken"`
	AccessTokenExpiresAt int64  `json:"accessTokenExpiresAt"`
	Message              string `json:"message"`
}

type RequestOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestOTPResponse struct {
	Message string `json:"message"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6,numeric"`
}

type VerifyOTPResponse struct {
	User                 UserResponse `json:"user"`
	AccessToken          string       `json:"accessToken"`
	RefreshToken         string       `json:"refreshToken"`
	Message              string       `json:"message"`
	AccessTokenExpiresAt int64        `json:"accessTokenExpiresAt"`
}

type Enable2FARequest struct {
	Enable bool `json:"enable"`
}

type Enable2FAResponse struct {
	Message      string `json:"message"`
	TwoFAEnabled bool   `json:"twoFaEnabled"`
}
