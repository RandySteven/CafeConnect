package enums

type RouterPrefix string

const (
	DevPrefix        RouterPrefix = `/dev`
	OnboardingPrefix RouterPrefix = `/onboarding`
	UserPrefix       RouterPrefix = `/users`
	CafePrefix       RouterPrefix = `/cafes`
	ReviewPrefix     RouterPrefix = `/reviews`
)

func (r RouterPrefix) ToString() string {
	return string(r)
}
