package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	drive "google.golang.org/api/drive/v3"
)

type ClientSecret struct {
	ClientID     string `json:"client_id"`
	ProjectID    string `json:"project_id"`
	AuthURI      string `json:"auth_uri"`
	TokenURI     string `json:"token_uri"`
	ClientSecret string `json:"client_secret"`
}

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

func getOrCreateFolder(d *drive.Service, folderName string) string {
	//folderId := ""
	if folderName == "" {
		return ""
	}
	q := fmt.Sprintf("name=\"%s\" and mimeType=\"application/vnd.google-apps.folder\"", folderName)
	print(q)
	r, err := d.Files.List().Do()
	if err != nil {
		fmt.Printf("Unable to retrieve Dlfer name")
	}
	for _, file := range r.Files {
		fmt.Println(file.Name)
	}
	return "done"
}

func readSecretsFromFile(path string) *oauth2.Config {
	var data map[string]ClientSecret
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File Not found.")
		return nil
	}
	json.Unmarshal(b, &data)
	fmt.Println(data)
	fmt.Println(data["installed"])
	client := data["installed"]
	return &oauth2.Config{
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
		Scopes:       []string{drive.DriveScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  client.AuthURI,
			TokenURL: client.TokenURI,
		},
		RedirectURL: "http://127.0.0.1:9004",
	}

}

func main() {
	config := readSecretsFromFile("client_secret.json")
	client := getClient(config)
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	getOrCreateFolder(srv, "unit2_assignment_02.py")
}
