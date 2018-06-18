package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	config := &oauth2.Config{
		ClientID:     "1055300065102-6jbjc6hc8inlnpme9bt1emnesta5b337.apps.googleusercontent.com ",
		ClientSecret: "Fk12e1-dycx0JXFK_wtax9HZ",
		Endpoint:     google.Endpoint,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "",
			TokenURL: "",
		},
	}
	//svc, err := drive.New(httpClient)

}
