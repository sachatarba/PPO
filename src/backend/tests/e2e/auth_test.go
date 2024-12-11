package e2e

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/bool64/godogx/allure"
	"github.com/cucumber/godog"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

type authFeature struct {
	expect   *httpexpect.Expect
	code     string
	status   string
	password string
	client   map[string]interface{}
}

func (a *authFeature) setup() {
	a.password = os.Getenv("PASSWORD")
	host := os.Getenv("AUTH_SERVER_HOST")
	if host == "" {
		host = "localhost"
	}

	a.expect = httpexpect.WithConfig(httpexpect.Config{
		Client:   &http.Client{},
		BaseURL:  fmt.Sprintf("http://%s:8099/api/v2/", host),
		Reporter: httpexpect.NewAssertReporter(nil),
	})
}

func ExtractCode(input string) (string, error) {
	re := regexp.MustCompile(`Для подтверждения авториазции введите код:\s(\S+)\r\n`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return "", fmt.Errorf("код авторизации не найден")
	}

	return matches[1], nil
}

func fetchCodeFromEmail(email, password string) (string, error) {
	c, err := connectToIMAPServer()
	if err != nil {
		return "", fmt.Errorf("Ошибка подключения к серверу: %v", err)
	}
	defer c.Logout()

	if err := authenticateIMAP(c, email, password); err != nil {
		return "", fmt.Errorf("Ошибка аутентификации: %v", err)
	}

	mbox, err := selectMailbox(c, "INBOX")
	if err != nil {
		return "", err
	}

	msg, err := fetchLatestMessage(c, mbox)
	if err != nil {
		return "", err
	}

	body, err := extractMessageBody(msg)
	if err != nil {
		return "", err
	}

	return ExtractCode(body)
}

func connectToIMAPServer() (*client.Client, error) {
	return client.DialTLS("imap.rambler.ru:993", &tls.Config{})
}

func authenticateIMAP(c *client.Client, email, password string) error {
	return c.Login(email, password)
}

func selectMailbox(c *client.Client, mailbox string) (*imap.MailboxStatus, error) {
	mbox, err := c.Select(mailbox, false)
	if err != nil {
		return nil, fmt.Errorf("Ошибка выбора папки: %v", err)
	}
	if mbox.Messages == 0 {
		return nil, fmt.Errorf("Папка пуста")
	}
	return mbox, nil
}

func fetchLatestMessage(c *client.Client, mbox *imap.MailboxStatus) (*imap.Message, error) {
	seqset := new(imap.SeqSet)
	seqset.AddNum(mbox.Messages)

	messages := make(chan *imap.Message, 1)
	section := &imap.BodySectionName{}
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages)
	}()

	msg := <-messages
	if msg == nil {
		return nil, fmt.Errorf("Сообщение не найдено")
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("Ошибка загрузки сообщения: %v", err)
	}

	return msg, nil
}

func extractMessageBody(msg *imap.Message) (string, error) {
	section := &imap.BodySectionName{}
	r := msg.GetBody(section)
	if r == nil {
		return "", fmt.Errorf("Ошибка чтения тела сообщения")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		return "", fmt.Errorf("Ошибка парсинга письма: %v", err)
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("Ошибка чтения части письма: %v", err)
		}

		body, err := io.ReadAll(part.Body)
		if err != nil {
			continue
		}
		return string(body), nil
	}

	return "", fmt.Errorf("Сообщение не содержит текста")
}

func (a *authFeature) aUserWithDetails(table *godog.Table) error {
	a.client = make(map[string]interface{})
	for _, row := range table.Rows[1:] {
		for i, cell := range row.Cells {
			column := table.Rows[0].Cells[i].Value
			a.client[column] = cell.Value
		}
	}

	a.client["ID"] = uuid.MustParse(a.client["ID"].(string))
	return nil
}

func (a *authFeature) theUserRegisters() error {
	response := a.expect.POST("/register").
		WithJSON(&a.client).
		Expect().
		Raw()

	a.status = response.Status
	return nil
}

func (a *authFeature) theUserLogsIn() error {
	login := map[string]interface{}{
		"Login":    a.client["Login"],
		"Password": a.client["Password"],
	}

	response := a.expect.POST("/login").
		WithJSON(&login).
		Expect().
		Raw()

	a.status = response.Status
	return nil
}

func (a *authFeature) theConfirmationCodeIsFetchedFromEmail(email string) error {
	code, err := fetchCodeFromEmail(email, a.password)
	if err != nil {
		return err
	}

	a.code = code
	return nil
}

func (a *authFeature) theUserConfirmsTheAccount() error {
	confirm := map[string]interface{}{
		"client_id": a.client["ID"],
		"code":      a.code,
	}
	response := a.expect.POST("/confirm").
		WithJSON(&confirm).
		Expect().
		Raw()

	a.status = response.Status
	time.Sleep(time.Second * 2)
	return nil
}

func (a *authFeature) theUserChangesThePasswordTo(newPassword string) error {
	changePassword := map[string]interface{}{
		"login":        a.client["Login"],
		"code":         a.code,
		"new_password": newPassword,
	}
	response := a.expect.POST("/change_password").
		WithJSON(&changePassword).
		Expect().
		Raw()
	a.client["Password"] = newPassword

	a.status = response.Status
	return nil
}

func (a *authFeature) theStatusCodeShouldBe(status string) error {
	if a.status != status {
		return fmt.Errorf("invalid status code: %s, but should be: %s", a.status, status)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	feature := &authFeature{}
	feature.setup()

	ctx.Step(`^a user with the following details:$`, feature.aUserWithDetails)
	ctx.Step(`^the user registers$`, feature.theUserRegisters)
	ctx.Step(`^the user logs in$`, feature.theUserLogsIn)
	ctx.Step(`^the confirmation code is fetched from the email "([^"]*)"$`, feature.theConfirmationCodeIsFetchedFromEmail)
	ctx.Step(`^the status code should be "([^"]*)"$`, feature.theStatusCodeShouldBe)
	ctx.Step(`^the user confirms the account$`, feature.theUserConfirmsTheAccount)
	ctx.Step(`^the user changes the password to "([^"]*)"$`, feature.theUserChangesThePasswordTo)
}

func TestFeatures(t *testing.T) {
	allure.RegisterFormatter()
	options := godog.Options{
		Format: "pretty",
		Paths:  []string{"features"},
	}

	suite := godog.TestSuite{
		Name:                "auth",
		ScenarioInitializer: InitializeScenario,
		Options:             &options,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
