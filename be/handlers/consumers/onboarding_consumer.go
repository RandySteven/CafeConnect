package consumers

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/messages"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	email_client "github.com/RandySteven/CafeConnect/be/pkg/email"
	"github.com/RandySteven/CafeConnect/be/utils"
	"log"
	"time"
)

const verifyTokenHTMLPath = `files/html/email/template/verify.token.html`

type OnboardingConsumer struct {
	onboardingTopic topics_interfaces.OnboardingTopic
	pointTopic      topics_interfaces.PointTopic
	email           email_client.Email
	userRepo        repository_interfaces.UserRepository
	verifyTokenRepo repository_interfaces.VerifyTokenRepository
	pointRepo       repository_interfaces.PointRepository
}

func (o *OnboardingConsumer) VerifyOnboardingToken(ctx context.Context) error {
	//consume(ctx, func(ctx context.Context) {
	//verifyTokenMessageStr, err := o.onboardingTopic.ReadMessage(ctx)
	//if err != nil {
	//	log.Println(`failed to consumer verify token message`, err)
	//	return
	//}
	//
	//tmpl, err := template.ParseFiles(verifyTokenHTMLPath)
	//if err != nil {
	//	log.Println(`failed to read template`, err)
	//	return
	//}
	//
	//verifyToken := utils.ReadJSONObject[messages.VerifyTokenMessage](verifyTokenMessageStr)
	//user, err := o.userRepo.FindByID(ctx, verifyToken.UserID)
	//if err != nil {
	//	log.Println(`failed to get user`, err)
	//	return
	//}
	//
	//_, err = o.verifyTokenRepo.Save(ctx, &models.VerifyToken{
	//	Token:       verifyToken.Token,
	//	UserID:      verifyToken.UserID,
	//	ExpiredTime: time.Now().Add(2 * time.Hour),
	//	IsClicked:   false,
	//})
	//if err != nil {
	//	log.Println(`failed to save verify token`, err)
	//	return
	//}
	//
	//contentMap := make(map[string]string)
	//contentMap[`token_url`] = verifyToken.Token
	//contentMap[`full_name`] = user.Name
	//
	//err = tmpl.Execute(nil, contentMap)
	//if err != nil {
	//	log.Println(`failed to read html`)
	//	return
	//}
	//})
	return nil
}

func (o *OnboardingConsumer) UserPointUpdate(ctx context.Context) error {
	return o.pointTopic.RegisterConsumer(func(message string) {
		userPointMessage := utils.ReadJSONObject[messages.TransactionPointMessage](message)

		user, err := o.userRepo.FindByID(ctx, userPointMessage.UserID)
		if err != nil {
			log.Println(`failed to get user`, err)
			return
		}

		point, err := o.pointRepo.FindByUserID(ctx, user.ID)
		if err != nil {
			log.Println(`failed to get point`, err)
			return
		}

		point.Point += userPointMessage.Point
		point.UpdatedAt = time.Now()
		_, err = o.pointRepo.Update(ctx, point)
		if err != nil {
			log.Println(`failed to update point`, err)
			return
		}
	})
}

var _ consumer_interfaces.OnboardingConsumer = &OnboardingConsumer{}

func newOnboardingConsumer(
	onboardingTopic topics_interfaces.OnboardingTopic,
	pointTopic topics_interfaces.PointTopic,
	email email_client.Email,
	userRepo repository_interfaces.UserRepository,
	verifyTokenRepo repository_interfaces.VerifyTokenRepository,
	pointRepo repository_interfaces.PointRepository,
) *OnboardingConsumer {
	return &OnboardingConsumer{
		onboardingTopic: onboardingTopic,
		pointTopic:      pointTopic,
		email:           email,
		userRepo:        userRepo,
		verifyTokenRepo: verifyTokenRepo,
		pointRepo:       pointRepo,
	}
}
