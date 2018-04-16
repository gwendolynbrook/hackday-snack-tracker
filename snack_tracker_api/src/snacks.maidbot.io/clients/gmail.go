package clients

import (
  "os"
  "net/smtp"
  "github.com/jordan-wright/email"
)

func emailInventorySummary(toAddress string) {
  // data_dir := os.Getenv(DATA_DIR_ENV)
  pw := os.Getenv(PASSWORD)

  e := email.NewEmail()
  e.From = "Maidbot Snacktracker <gwen@maidbot.com>"
  e.To = []string{toAddress}
  e.Cc = []string{"gwen@maidbot.com"}
  e.Subject = "Snacktracker Summary!"
  e.Text = []byte("tadaa!")
  e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "gwen@maidbot.com", pw, "smtp.gmail.com"))
}
