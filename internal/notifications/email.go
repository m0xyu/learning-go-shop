package notifications

import (
	"fmt"
	"net"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SimpleEmail struct {
	To      string
	Subject string
	Body    string
}

type EmailNotifier struct {
	config *SMTPConfig
}

func NewEmailNotifier(config *SMTPConfig) *EmailNotifier {
	return &EmailNotifier{
		config: config,
	}
}

// SendSimpleEmail は、SMTPサーバーを使用してシンプルなメールを送信します。
func (e *EmailNotifier) SendSimpleEmail(email *SimpleEmail) error {
	addr := net.JoinHostPort(e.config.Host, fmt.Sprintf("%d", e.config.Port))

	// 開発用にTLSなしで直接接続する
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.config.Host)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Quit()
	}()

	// 認証情報がある場合は認証する
	if e.config.Username != "" || e.config.Password != "" {
		auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	// メールの送信
	if err := client.Mail(e.config.From); err != nil {
		return err
	}

	// 送信先を追加
	if err := client.Rcpt(email.To); err != nil {
		return err
	}

	// データの書き込み
	w, err := client.Data()
	if err != nil {
		return err
	}

	// メールの内容をフォーマットして送信
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.config.From, email.To, email.Subject, email.Body)

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	return w.Close()
}

// SendLoginNotification は、ユーザーがログインした際に通知メールを送信します。
func (e *EmailNotifier) SendLoginNotification(userEmail, userName string) error {
	email := &SimpleEmail{
		To:      userEmail,
		Subject: "Login Notification",
		Body: fmt.Sprintf(`Hello %s,

You have successfully logged into your account.

If this wasn't you, please contact support immediately.

Best regards,
The Shop Team`, userName),
	}

	return e.SendSimpleEmail(email)
}
