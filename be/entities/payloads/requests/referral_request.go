package requests

type (
	FindReferralRequest struct {
		ReferralID   uint64 `json:"referral_id"`
		UserID       uint64 `json:"user_id"`
		Code         string `json:"code"`
		IdentifiedBy string `json:"identified_by"`
	}
)
