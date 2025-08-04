package emailuser

import (
	"erp/gen/db/model"
	"erp/internal/app/config"
	"erp/internal/app/connection"
	"erp/internal/app/event-bus/event"
	"erp/internal/app/plugin/email/emailtypes"
	emailservice "erp/internal/app/plugin/email/service"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services"
	userservice "erp/internal/app/service/services/user_service"
	"erp/pkg/logger"
	"fmt"

	"github.com/matcornic/hermes/v2"
)

type UserEmailHandler struct {
	configService *config.ConfigService
	conn          *connection.Connection
	emailService  emailservice.EmailService
	locale        helpers.Locale
	logger        helpers.EmitLog
	userService *userservice.UserService
}

func NewEmailUserHandler(
	configService *config.ConfigService,
	conn *connection.Connection,
	helpers *helpers.Helpers,
	emailService emailservice.EmailService,
	services *services.Services,
) *UserEmailHandler {
	return &UserEmailHandler{
		configService: configService,
		conn:          conn,
		logger:        helpers.Logger.EmitLog("email-user"),
		emailService:  emailService,
		locale:        helpers.Locale,
		userService: services.UserService,
	}
}

func (h *UserEmailHandler) SendUserCredentials(payload *event.NotificationData) {
	fmt.Println("PAYLOAD EMAIL",payload)
	payloadData, ok := payload.Data.Payload.(model.UserRelation)
	if !ok {
		h.logger.Err(config.ErrTypeAssertion)
	}

	userPassword,err := h.userService.GetUserPassword(payloadData.User.ID)
	if err != nil {
		h.logger.Err(err,logger.OptionsLog.WithMethod("GetUserPassword"))
		return
	}

	clientConfig := h.configService.GetClientConfig()
	reqCtx := payload.Data.RequestContext
	b := h.emailService.EmailBuilder(&emailtypes.EmailBuilderOpts{
		// RequestContext: &reqCtx,
	})
	email := hermes.Email{
		Body: hermes.Body{
			Name: h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Email.Dear"),
				helpers.OptionsLocale.WithLang(string(reqCtx.LanguageCode)),
				helpers.OptionsLocale.WithTemplate(map[string]string{
					"FirstName": payloadData.Profile.GivenName,
					"LastName":  payloadData.Profile.FamilyName,
				}),
			),
			Intros: []string{
				h.locale.MustLocalize(
					helpers.OptionsLocale.WithID("EmailIntro.Credentials"),
					helpers.OptionsLocale.WithLang(string(reqCtx.LanguageCode)),
					helpers.OptionsLocale.WithTemplate(map[string]string{
						"CompanyName": reqCtx.ActiveCompany.Name,
					}),
				),
			},
			Dictionary: []hermes.Entry{
				{Key: h.locale.MustLocalize(
					helpers.OptionsLocale.WithID("Email.Base"),
					helpers.OptionsLocale.WithLang(string(reqCtx.LanguageCode)),
				),
				Value: payloadData.User.Identifier,
			},
			
			{Key: h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("Base.Password"),
				helpers.OptionsLocale.WithLang(string(reqCtx.LanguageCode)),
			),
			Value: userPassword,
		},
	},
	Actions: []hermes.Action{
		{
			Instructions: h.locale.MustLocalize(
				helpers.OptionsLocale.WithID("EmailAction.StartedIntro"),
				helpers.OptionsLocale.WithLang(string(reqCtx.LanguageCode)),
			),
			Button: hermes.Button{
				Color: "#22BC66", // Optional action button color
				Text:  h.locale.MustLocalize(
					helpers.OptionsLocale.WithID("EmailAction.SignIn"),
					helpers.OptionsLocale.WithLang(string(reqCtx.LanguageCode)),
				),
				Link:  fmt.Sprintf("%s/signin?uuid=%s",clientConfig.Url,payloadData.Company.UUID),
			},
		},
	},
	Outros: []string{
				h.locale.MustLocalize(
					helpers.OptionsLocale.WithID("Email.DefaultOutor"),
					helpers.OptionsLocale.WithLang(string(reqCtx.LanguageCode)),
				),
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := b.Hermes.GenerateHTML(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}
	b.SetTo(payloadData.User.Identifier)
	b.SetBody(emailBody)
	b.SetSubject(h.locale.MustLocalize(helpers.OptionsLocale.WithID("EmailSubject.Credentials"),
		helpers.OptionsLocale.WithLang(string(reqCtx.LanguageCode))))
	b.Send()
	fmt.Println("CLIENT EMAIL", payloadData)

}
