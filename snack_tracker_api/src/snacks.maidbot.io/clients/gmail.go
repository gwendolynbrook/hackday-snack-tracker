package clients

import (
  "os"
  "log"
  "net/smtp"
  "path/filepath"
  "github.com/jordan-wright/email"

  "app/src/snacks.maidbot.io/data"
)

func EmailInventorySummary(toAddress string, cleanUp bool) {
  dataDir := os.Getenv(data.DATA_DIR_ENV)
  pw := os.Getenv(GMAIL_PASSWORD)
  csvFiles, err := filepath.Glob(dataDir + "/*.csv")

  e := email.NewEmail()
  e.From = "Maidbot Snacktracker <mbsnacktracker@gmail.com>"
  e.To = []string{toAddress}
  e.Cc = []string{"gwen@maidbot.com"}
  e.Subject = "Snacktracker Summary!"
  e.Text = []byte("tadaa!")

  if err != nil {
    log.Print("Failed to glob csv files")
    e.Text = []byte("Snacktracker had an error. Ask Gwen for help.")
  }

  for _, fn := range csvFiles {
      e.AttachFile(fn)
  }

  e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "gwen@maidbot.com", pw, "smtp.gmail.com"))
  if cleanUp {
    data.CleanupCsvs()
  }
}
