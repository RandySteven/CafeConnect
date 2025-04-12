package requests

type (
	RegisterUserRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Referral  *struct {
			Code string `json:"code"`
		} `json:"referral"`
	}

	LoginUserRequest struct {
	}
)
