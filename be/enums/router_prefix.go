package enums

type RouterPrefix string

const (
	DevPrefix        RouterPrefix = `/dev`
	OnboardingPrefix RouterPrefix = `/onboarding`
	UserPrefix       RouterPrefix = `/users`
	CafePrefix       RouterPrefix = `/cafes`
	ReviewPrefix     RouterPrefix = `/reviews`
	ProductPrefix    RouterPrefix = `/products`
	CartPrefix       RouterPrefix = `/carts`
)

func (r RouterPrefix) ToString() string {
	return string(r)
}
