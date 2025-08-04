package email_event

import (
	"context"
	"erp/internal/app/plugin/email/emailtypes"
	emailservice "erp/internal/app/plugin/email/service"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/config"
	"erp/pkg/db"
	"erp/pkg/logger"
	"fmt"

	"github.com/matcornic/hermes/v2"
)

type UserEventHandler struct {
	bus          bus.Bus
	emailService emailservice.EmailService
	emitLog      logger.EmitLog
	appConfig    *config.AppConfig
	locale       helpers.Locale
	conn         db.Connection
}

func NewUserEventHandler(
	conn db.Connection,
	helpers *helpers.Helpers,
	appConfig *config.AppConfig,
	bus bus.Bus,
	logger logger.Logger,
	emailService emailservice.EmailService,
) {
	handler := UserEventHandler{
		conn:         conn,
		bus:          bus,
		emitLog:      logger.EmitLog("user-vent-email"),
		emailService: emailService,
		locale:       helpers.Locale,
		appConfig:    appConfig,
	}
	bus.RegisterHandler(domain.UserCreatedEvent, handler.OnUserCreated())
}

func (h *UserEventHandler) OnUserCreated() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			go func ()  {
				payload, ok := e.Data.(event.UserCreatedEventData)
				if !ok {
					return 
				}
				payloadData := payload.UseRelation
				languageCode := string(payload.LanguageCode)
				
				userPassword, err := h.GetUserPassword(payloadData.User.ID)
				if err != nil {
					h.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserPassword"))
					return 
				}
				
				clientConfig := h.appConfig.Client
				b := h.emailService.EmailBuilder(&emailtypes.EmailBuilderOpts{
					LanguageCode: string(payload.LanguageCode),
					Company:      &payloadData.Company,
				})
				email := hermes.Email{
					Body: hermes.Body{
						Name: h.locale.MustLocalize(
							helpers.OptionsLocale.WithID("Email.Dear"),
							helpers.OptionsLocale.WithLang(languageCode),
							helpers.OptionsLocale.WithTemplate(map[string]string{
								"FirstName": payloadData.Profile.GivenName,
								"LastName":  payloadData.Profile.FamilyName,
							}),
						),
						Intros: []string{
							h.locale.MustLocalize(
								helpers.OptionsLocale.WithID("EmailIntro.Credentials"),
								helpers.OptionsLocale.WithLang(languageCode),
							),
						},
						Dictionary: []hermes.Entry{
							{Key: h.locale.MustLocalize(
								helpers.OptionsLocale.WithID("Email.Base"),
								helpers.OptionsLocale.WithLang(languageCode),
							),
							Value: payloadData.User.Identifier,
						},
						
						{Key: h.locale.MustLocalize(
							helpers.OptionsLocale.WithID("Base.Password"),
							helpers.OptionsLocale.WithLang(languageCode),
						),
						Value: userPassword,
					},
				},
				Actions: []hermes.Action{
					{
						Instructions: h.locale.MustLocalize(
							helpers.OptionsLocale.WithID("EmailAction.StartedIntro"),
							helpers.OptionsLocale.WithLang(languageCode),
						),
						Button: hermes.Button{
							Color: "#22BC66", // Optional action button color
							Text: h.locale.MustLocalize(
								helpers.OptionsLocale.WithID("EmailAction.SignIn"),
								helpers.OptionsLocale.WithLang(languageCode),
							),
							Link: fmt.Sprintf("%s/signin?uuid=%s", clientConfig.Url, payloadData.Company.UUID),
						},
					},
				},
				Outros: []string{
					h.locale.MustLocalize(
						helpers.OptionsLocale.WithID("Email.DefaultOutor"),
						helpers.OptionsLocale.WithLang(languageCode),
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
		helpers.OptionsLocale.WithLang(languageCode)))
		b.Send()
		fmt.Println("CLIENT EMAIL", payloadData)
		}()
		return nil
	},
		Matcher: domain.UserCreatedEvent,
	}
}

func (s *UserEventHandler) GetUserPassword(id int64) (string, error) {
	dbConfig := s.appConfig.PG
	var password string
	err := s.conn.GetDB().Raw("select pgp_sym_decrypt(password_hash::bytea, ?) as password_hash from users where id = ?", dbConfig.CryptoPass, id).
	Scan(&password).Error
	return password, err
}
