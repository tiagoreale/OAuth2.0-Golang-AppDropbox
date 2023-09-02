package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
)

func main() {
	config := oauth2.Config{
		ClientID:     "&&&&&&&",
		ClientSecret: "$$$$$$$",
		Scopes:       []string{"files.metadata.write", "files.content.write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
		RedirectURL: "http://localhost:8080/dbupload",
	}

	refreshToken := "$$$$$$$$$&&&&&&&&&&&&"

	tokenSource := config.TokenSource(oauth2.NoContext, &oauth2.Token{RefreshToken: refreshToken})
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatal("Erro ao obter token de acesso: ", err)
	}

	fmt.Println("Token de Acesso:", newToken.AccessToken)

	// Salvar o novo token de acesso temporário em um arquivo
	err = ioutil.WriteFile("/home/tiagoreale/go/src/DropBox/aplicacao/temp.txt", []byte(newToken.AccessToken), 0644)
	if err != nil {
		log.Fatal("Erro ao salvar o token em arquivo: ", err)
	}
}
