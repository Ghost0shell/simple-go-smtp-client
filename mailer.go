package main

import (
  "fmt"
  "log"
  "net/smtp"
  "strings"
  "bytes"
  "os"
  "gopkg.in/yaml.v2"
)


type conf struct {
  From string `yaml:"from"`
  To string `yaml:"to"`
  Host string `yaml:"host"`
  Port int64 `yaml:"port"`
  Subject string `yaml:"subject"`
  Body string `yaml:"body"`
}

// function to get mail parameters from file
func (c *conf) getConf() *conf {
  configFile, err := os.ReadFile("mail_settings.yaml")
  if err != nil {
    log.Printf("configFile.Get err #%v ", err)
  }
  err = yaml.Unmarshal(configFile, c)
  if err != nil {
    log.Fatalf("Unmarshal: %s", err)
  }

  return c
}


func main() {
  var c conf
  c.getConf()

  toHeader := strings.Join(strings.Split(c.To, ","), ",")

  fmt.Println("will send email from " + c.From + ", to " + toHeader)

  client, err := smtp.Dial(fmt.Sprintf("%s:%d", c.Host, c.Port))

  if err != nil {
    log.Fatal(err)
  }

  defer client.Close()

  if err := client.Mail(c.From); err != nil {
    log.Fatal(err)
  }
  for _, s := range strings.Split(c.To, ",") {
    if err := client.Rcpt(string(s)); err != nil {
    log.Fatal(err)
    }
  }

  wclient, err := client.Data()
  if err != nil {
    log.Fatal(err)
  }
  defer wclient.Close()

  buf := bytes.NewBufferString("Subject: " + c.Subject + c.Body)

  if _, err = buf.WriteTo(wclient); err != nil {
    log.Fatal(err)
  }

  fmt.Println("Email sent!")

}
