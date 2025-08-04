package email_event

import (
	"context"
	"erp/api/common"
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

type AccountEventHandler struct {
	bus          bus.Bus
	emailService emailservice.EmailService
	emitLog      logger.EmitLog
	appConfig    *config.AppConfig
	locale       helpers.Locale
	conn         db.Connection
	session      helpers.SessionHelper
	jwt          helpers.JwtHelper
}

func NewAccountEventHandler(
	conn db.Connection,
	helpers *helpers.Helpers,
	appConfig *config.AppConfig,
	bus bus.Bus,
	logger logger.Logger,
	emailService emailservice.EmailService,
) {
	handler := AccountEventHandler{
		conn:         conn,
		bus:          bus,
		emitLog:      logger.EmitLog("account-event-email"),
		emailService: emailService,
		locale:       helpers.Locale,
		appConfig:    appConfig,
		session:      helpers.Session,
		jwt:          helpers.Jwt,
	}
	bus.RegisterHandler(domain.PasswordResetEvent, handler.OnPasswordReset())
}

func (h *AccountEventHandler) OnPasswordReset() bus.Handler {
	return bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) error {
			go func() {
				defer func() {
					fmt.Println("DEFER ON PASSWORD RESET...")
				}()
				payload, ok := e.Data.(event.PasswordResetEventData)
				if !ok {
					return
				}
				payloadData := payload
				languageCode := h.session.ParseAcceptLanguage(payload.LanguageCode)

				clientConfig := h.appConfig.Client
				token,err := h.jwt.GenerateToken(common.Claims{
					ID:   payloadData.User.ID,
					Uuid: payloadData.User.UUID,
				})
				if err != nil {
					return
				}

				b := h.emailService.EmailBuilder(&emailtypes.EmailBuilderOpts{
					LanguageCode: string(payload.LanguageCode),
					// Company:      payloadData.Company,
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
								helpers.OptionsLocale.WithID("EmailIntro.PasswordReset"),
								helpers.OptionsLocale.WithLang(languageCode),
								helpers.OptionsLocale.WithTemplate(map[string]string{
									"Email": payloadData.User.Identifier,
								}),
							),
						},

						Actions: []hermes.Action{
							{
								Button: hermes.Button{
									Color: "#22BC66", // Optional action button color
									Text: h.locale.MustLocalize(
										helpers.OptionsLocale.WithID("EmailAction.ResetYouPassword"),
										helpers.OptionsLocale.WithLang(languageCode),
									),
									Link: fmt.Sprintf("%s/change-password?c=%s", clientConfig.Url, token),
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
					h.emitLog.Err(err) // Tip: Handle error with something else than a panic ;)
				}
				b.SetTo(payloadData.User.Identifier)
				b.SetBody(emailBody)
				b.SetSubject(
					h.locale.MustLocalize(
						helpers.OptionsLocale.WithID("EmailAction.ResetYouPassword"),
						helpers.OptionsLocale.WithLang(languageCode),
					),
				)
				b.Send()
				fmt.Println("CLIENT EMAIL", payloadData)
			}()
			return nil
		},
		Matcher: domain.PasswordResetEvent,
	}
}
