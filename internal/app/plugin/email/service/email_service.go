package emailservice

import (
	"erp/internal/app/plugin/email/emailtypes"
	"erp/internal/app/service/helpers"
	"erp/pkg/config"
	"fmt"
	"time"

	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	EmailBuilder(options *emailtypes.EmailBuilderOpts) *BuilderEmail
}

type emailService struct {
	appConfig *config.AppConfig
	locale helpers.Locale
}

func NewEmailService(
	appConfig *config.AppConfig,
	helpers *helpers.Helpers,
) *emailService {
	return &emailService{
		appConfig:  appConfig,
		locale:         helpers.Locale,
	}
}

func (s *emailService) getTransportLayer() *gomail.Dialer {
	transportOpts := s.appConfig.Email
	fmt.Println("Transport Options",transportOpts)
	return gomail.NewDialer(transportOpts.Host, transportOpts.Port, transportOpts.User, transportOpts.Password)
}

func (s *emailService) EmailBuilder(options *emailtypes.EmailBuilderOpts) *BuilderEmail {
	clientConfig := s.appConfig.Client
	transportOpts := s.appConfig.Email
	year := time.Now().Year()
	herm := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Teclu",
			Link: clientConfig.Url,
			// Optional product logo
			Logo: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRrQYvpiIiQiQm598N84WpLQcUgW-RgcwhfvQ&s",
		},

	}
	if options.Company != nil {
		herm.Product.Name = options.Company.Name
		herm.Product.Copyright = s.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Email.Copyright"),
			helpers.OptionsLocale.WithLang(options.LanguageCode),
			helpers.OptionsLocale.WithTemplate(map[string]string{
				"Year":fmt.Sprintf("%d",year),
				"CompanyName":options.Company.Name,
			}),
		)
	}
	return &BuilderEmail{
		dial:   s.getTransportLayer(),
		from:   transportOpts.User,
		Hermes: &herm,
	}
}

type BuilderEmail struct {
	from string
	to   []string
	cc   struct {
		address string
		name    string
	}
	subject string
	body    string
	dial    *gomail.Dialer
	Hermes  *hermes.Hermes
}

// func (b *BuilderEmail) SetFrom(v string) {
// 	b.from = v
// }

func (b *BuilderEmail) SetTo(v ...string) { b.to = v }

func (b *BuilderEmail) SetCc(address string, name string) {
	b.cc.address = address
	b.cc.name = name
}

func (b *BuilderEmail) SetSubject(v string) { b.subject = v }

func (b *BuilderEmail) SetBody(v string) { b.body = v }

func (b *BuilderEmail) Send() {
	m := gomail.NewMessage()
	m.SetHeader("From", b.from)
	m.SetHeader("To", b.to...)
	// m.SetAddressHeader("Cc", b.cc.address,b.cc.name)
	m.SetHeader("Subject", b.subject)
	m.SetBody("text/html", b.body)

	// m.Attach("/home/Alex/lolcat.jpg")
	// Send the email to Bob, Cora and Dan.
	if err := b.dial.DialAndSend(m); err != nil {
		fmt.Println(err)
	}

}
