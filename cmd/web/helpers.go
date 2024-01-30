package main

func (app *Configuration) sendEmail(msg Message) {
	app.Wait.Add(1)
	app.Mailer.MailerChan <-msg
}