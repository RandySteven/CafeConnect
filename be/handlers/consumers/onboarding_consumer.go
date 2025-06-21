package consumers

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/messages"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	email_client "github.com/RandySteven/CafeConnect/be/pkg/email"
	"github.com/RandySteven/CafeConnect/be/utils"
	"html/template"
	"log"
)

const verifyTokenHTMLPath = `files/html/email/template/verify.token.html`

type OnboardingConsumer struct {
	onboardingTopic topics_interfaces.OnboardingTopic
	email           email_client.Email
	userRepo        repository_interfaces.UserRepository
}

func (o *OnboardingConsumer) VerifyOnboardingToken(ctx context.Context) {
	consume(ctx, func(ctx context.Context) {
		verifyTokenMessageStr, err := o.onboardingTopic.ReadMessage(ctx, ``)
		if err != nil {
			log.Println(`failed to consumer verify token message`, err)
			return
		}

		tmpl, err := template.ParseFiles(verifyTokenHTMLPath)
		if err != nil {
			log.Println(`failed to read template`, err)
			return
		}

		verifyToken := utils.ReadJSONObject[messages.VerifyTokenMessage](verifyTokenMessageStr)
		user, err := o.userRepo.FindByID(ctx, verifyToken.UserID)

		contentMap := make(map[string]string)
		contentMap[`token_url`] = verifyToken.Token
		contentMap[`full_name`] = user.Name

		err = tmpl.Execute(nil, contentMap)
		if err != nil {
			log.Println(`failed to read html`)
			return
		}
	})
}

var _ consumer_interfaces.OnboardingConsumer = &OnboardingConsumer{}

func newOnboardingConsumer(
	onboardingTopic topics_interfaces.OnboardingTopic,
	email email_client.Email,
	userRepo repository_interfaces.UserRepository,
) *OnboardingConsumer {
	return &OnboardingConsumer{
		onboardingTopic: onboardingTopic,
		email:           email,
		userRepo:        userRepo,
	}
}
