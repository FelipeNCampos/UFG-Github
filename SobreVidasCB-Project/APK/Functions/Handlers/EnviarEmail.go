package Handlers
import (
    "fmt"
    "net/smtp"
)

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
   }
   // Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}
func enviarEmail(dest string,mensagem string) {
    // Sender data.
    from := "sobrevidascb@gmail.com"
    password := "P@ss12345"
    // Receiver email address.
    to := []string{dest}
    // smtp server configuration.
    smtpServer := smtpServer{host:"smtp.gmail.com", port:"587"}
    // Message.
    message := []byte(mensagem)
    // Authentication.
    auth := smtp.PlainAuth("", from, password, smtpServer.host)
    // Sending email.
    err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Email Sent!")
}