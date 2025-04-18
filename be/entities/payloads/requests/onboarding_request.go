package requests

import (
	"io"
)

type (
	RegisterUserRequest struct {
		FirstName      string    `form:"first_name"`
		LastName       string    `form:"last_name"`
		Username       string    `form:"username"`
		Email          string    `form:"email"`
		Password       string    `form:"password"`
		ProfilePicture io.Reader `form:"profile_picture"`
		PhoneNumber    string    `form:"phone_number"`
		DoB            string    `form:"dob"`
		ReferralCode   string    `form:"referral_code"`
		Address        string    `form:"address"`
		Longitude      float32   `form:"longitude"`
		Latitude       float32   `form:"latitude"`
	}

	LoginUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)
