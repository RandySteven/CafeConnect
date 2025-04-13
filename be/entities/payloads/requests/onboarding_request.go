package requests

import "time"

type (
	RegisterUserRequest struct {
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		Username       string    `json:"username"`
		Email          string    `json:"email"`
		Password       string    `json:"password"`
		ProfilePicture string    `json:"profile_picture"`
		PhoneNumber    string    `json:"phone_number"`
		DoB            time.Time `json:"dob"`
		Referral       *struct {
			Code string `json:"code"`
		} `json:"referral"`
		Address *struct {
			Address   string  `json:"address"`
			Longitude float32 `json:"longitude"`
			Latitude  float32 `json:"latitude"`
		} `json:"address"`
	}

	LoginUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)
