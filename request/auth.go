package request

type Auth struct {
	Account     *string `json:"account"`
	Password    *string `json:"password"`
	MobilePhone *string `json:"mobile_phone"`
	VerifyCode  *string `json:"verify_code"`
	VerifyBy    string  `json:"verify_by"`
}
