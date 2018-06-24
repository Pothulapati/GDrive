package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	drive "google.golang.org/api/drive/v3"
)

func getClient(config *oauth2.Config) *http.Client {

	token, err := tokenFromFile("token.json")
	if err != nil {
		ctx, cancel := context.WithCancel(context.Background())
		url := config.AuthCodeURL("", oauth2.AccessTypeOffline)
		browser.OpenURL(url)
		authcode := GetAuthorizationCode(ctx, cancel)
		token, err = config.Exchange(context.Background(), authcode)
		if err != nil {
			print(err.Error())
		}
		saveToken(token)
	}
	return config.Client(context.Background(), token)
}

func tokenFromFile(path string) (*oauth2.Token, error) {

	f, err := os.Open(path)
	defer f.Close()
	tok := &oauth2.Token{}
	json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(token *oauth2.Token) {

	fmt.Printf("Saving Credential file")
	f, err := os.OpenFile("token.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		fmt.Print("Can't Save Token")
	}
	json.NewEncoder(f).Encode(token)
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

func main() {

	config := &oauth2.Config{
		ClientID:     "1055300065102-6jbjc6hc8inlnpme9bt1emnesta5b337.apps.googleusercontent.com",
		ClientSecret: "Fk12e1-dycx0JXFK_wtax9HZ",
		Scopes:       []string{drive.DriveScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
		RedirectURL: "http://127.0.0.1:9004",
	}
	client := getClient(config)
	srv, err := drive.New(client)
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
