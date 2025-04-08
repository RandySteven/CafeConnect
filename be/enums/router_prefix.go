package enums

type RouterPrefix string

const (
	DevPrefix        RouterPrefix = `/dev`
	OnboardingPrefix RouterPrefix = `/onboarding`
	UserPrefix       RouterPrefix = `/users`
)

func (r RouterPrefix) ToString() string {
	return string(r)
}
