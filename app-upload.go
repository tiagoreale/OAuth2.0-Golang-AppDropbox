package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func uploadChunk(dbx files.Client, file *os.File, sessionID string, offset uint64, bar *pb.ProgressBar) (uint64, error) {
	chunk := make([]byte, 1<<22)
	bytesRead, err := file.Read(chunk)
	if err != nil && err != io.EOF {
		return offset, err
	}
	if bytesRead == 0 {
		return offset, io.EOF
	}
	chunkReader := bytes.NewReader(chunk[:bytesRead])
	err = dbx.UploadSessionAppendV2(files.NewUploadSessionAppendArg(
		&files.UploadSessionCursor{
			SessionId: sessionID,
			Offset:    offset,
		}), chunkReader)
	bar.Add(bytesRead)
	return offset + uint64(bytesRead), err
}

func formatFileSize(size int64) string {
	units := []string{"B", "KB", "MB", "GB"}
	i := 0
	for size >= 1024 && i < len(units)-1 {
		size /= 1024
		i++
	}
	return fmt.Sprintf("%d %s", size, units[i])
}

func sendTelegramMessage(title, message string) error {
	botToken := "&&&&&&&&&&&&&&&&&"
	chatID := "&&&&&&&&&&&&"
	apiUrl := "https://api.telegram.org/bot" + botToken + "/sendMessage"

	params := url.Values{}
	params.Set("chat_id", chatID)
	params.Set("text", title+"\n"+message)

	_, err := http.PostForm(apiUrl, params)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	logFile, err := os.OpenFile("/home/tiagoreale/go/src/DropBox/log/upload.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Erro ao abrir o arquivo de log: ", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("BrWiki - Upload Dropbox - Tiago Reale dos Reis")

	compressaokCmd := exec.Command("/home/tiagoreale/go/src/DropBox/aplicacao/appk")
	compressaokCmd.Stdout = os.Stdout
	compressaokCmd.Stderr = os.Stderr
	if err := compressaokCmd.Run(); err != nil {
		log.Fatal("Erro ao executar o binário de compressão: ", err)
	}

	tokenBytes, err := ioutil.ReadFile("/home/tiagoreale/go/src/DropBox/aplicacao/temp.txt")
	if err != nil {
		log.Fatal("Erro ao ler o token de acesso do arquivo: ", err)
	}

	accessToken := string(tokenBytes)

	dbx := files.New(dropbox.Config{Token: accessToken})

	filesInDir, err := ioutil.ReadDir("/home/tiagoreale/go/src/DropBox/coletor")
	if err != nil {
		log.Fatal("Erro ao listar arquivos na pasta: ", err)
	}

	var mostRecentFile os.FileInfo
	for _, fileInfo := range filesInDir {

		if fileInfo.IsDir() || filepath.Ext(fileInfo.Name()) != ".zip" {
			continue
		}
		if mostRecentFile == nil || fileInfo.ModTime().After(mostRecentFile.ModTime()) {
			mostRecentFile = fileInfo
		}
	}

	if mostRecentFile == nil {
		log.Fatal("Nenhum arquivo vma.lzo encontrado na pasta.")
	}

	filePath := filepath.Join("/home/tiagoreale/go/src/DropBox/coletor", mostRecentFile.Name())

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Erro ao abrir o arquivo para upload: ", err)
	}
	defer file.Close()

	commitInfo := files.CommitInfo{Path: "/DB/" + mostRecentFile.Name()}
	commitInfo.Mode = &files.WriteMode{Tagged: dropbox.Tagged{Tag: "add"}}

	res, err := dbx.UploadSessionStart(files.NewUploadSessionStartArg(), nil)
	if err != nil {
		log.Fatal("Erro ao iniciar sessão de upload: ", err)
	}

	fileInfo, _ := file.Stat()
	bar := pb.Full.Start64(fileInfo.Size())

	log.Printf("Tamanho do arquivo enviado: %d bytes", fileInfo.Size())

	offset := uint64(0)
	for {
		offset, err = uploadChunk(dbx, file, res.SessionId, offset, bar)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Erro durante o upload: ", err)
		}
	}

	finishArg := files.NewUploadSessionFinishArg(
		&files.UploadSessionCursor{
			SessionId: res.SessionId,
			Offset:    offset,
		},
		&commitInfo)

	_, err = dbx.UploadSessionFinish(finishArg, nil)
	if err != nil {
		log.Fatal("Erro ao finalizar sessão de upload: ", err)
	} else {
		log.Println("Upload concluído com sucesso.")
	}

	telegramTitle := "DB-APP-105"
	telegramMessage := fmt.Sprintf("Backup e Upload DB-UNIX-105 %s concluído. Tamanho do arquivo: %s", mostRecentFile.Name(), formatFileSize(fileInfo.Size()))

	if err := sendTelegramMessage(telegramTitle, telegramMessage); err != nil {
		log.Println("Erro ao enviar mensagem para o Telegram:", err)
	} else {
		log.Println("Mensagem enviada para o Telegram com sucesso.")
	}
}
