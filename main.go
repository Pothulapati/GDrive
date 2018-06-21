package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	drive "google.golang.org/api/drive/v3"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	conf := &oauth2.Config{
		ClientID:     "1055300065102-6jbjc6hc8inlnpme9bt1emnesta5b337.apps.googleusercontent.com",
		ClientSecret: "Fk12e1-dycx0JXFK_wtax9HZ",
		Scopes:       []string{drive.DriveScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
		RedirectURL: "http://127.0.0.1:9004",
	}
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	browser.OpenURL(url)
	authcode := GetAuthorizationCode(ctx, cancel)
	print(authcode)
	//err := drive.New(httpClient)
}

//Starts a server listens to the code and returns it
func GetAuthorizationCode(ctx context.Context, cancel context.CancelFunc) string {
	var x string
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		x = fmt.Sprint(r.URL.Query().Get("code"))
		io.WriteString(w, "Autherization Successful. You may Close the window now.")
		cancel()
	})
	srv := &http.Server{Addr: ":9004"}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			print("Not able to listen to 9004")
		}
	}()
	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil && err != context.Canceled {
		fmt.Print(err)
	}
	return x
}
