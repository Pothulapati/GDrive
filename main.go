package main

import (
	"context"
	"fmt"
	"io"
	"log"
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
	tok, err := conf.Exchange(context.Background(), authcode)
	if err != nil {
		print(err.Error())
	}
	srv, err := drive.New(conf.Client(context.Background(), tok))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	r, err := srv.Files.List().PageSize(20).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}

}

func GetAuthorizationCode(ctx context.Context, cancel context.CancelFunc) string {
	var x string
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		x = fmt.Sprint(r.URL.Query().Get("code"))
		io.WriteString(w, "Autherization Successful. You may Close the window now.")
		cancel()
	})
	srv := &http.Server{Addr: ":9004"}
	go func() {
		srv.ListenAndServe()
	}()
	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil && err != context.Canceled {
		fmt.Print("Error at shutdown")
		fmt.Print(err)
	}
	return x
}
