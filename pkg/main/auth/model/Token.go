package model

type Token struct {
	AccessToken          string `json:"accessToken"`
	RefreshToken         string `json:"refreshToken"`
	Auth_Id              string `json:"auth_id"`
	Auth_Context         string `json:"auth_context"`
	Auth_Account_Profile string `json:"auth_account_profile"`
}
