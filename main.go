package main
import (
    "context"
    "os"
    "net/http"
    "encoding/json"
    "log"
    "fmt"
    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/drive/v3"
    "google.golang.org/api/option"
)

func main() {
    router := gin.Default()
    router.POST("/upload", upload)

    router.Run("localhost:8080")
}

func upload(c *gin.Context) {
    service := getDriveService()
    files, _ := os.ReadDir("./images")

    for _, file := range files {
        fmt.Println(file.Name())
        file_name := file.Name()
        f, _ := os.Open("./images/" + file_name)
        go upload_file(service, f)
    }
    c.IndentedJSON(http.StatusOK, "upload starting")
}

func upload_file(service *drive.Service, file *os.File) {
    _, err := service.Files.Create(
        &drive.File{Name: file.Name(), Parents: []string{"1MKf6wfl1exlSR6aa34vnjdkqsYuU8_PM"}}).Media(file).Do()

    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }
}

func getDriveService() *drive.Service {
    ctx := context.Background()
    b, err := os.ReadFile("credentials.json")
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }

    // If modifying these scopes, delete your previously saved token.json.
    config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/drive.file")
    if err != nil {
        log.Fatalf("Unable to parse client secret file to config: %v", err)
    }
    client := getClient(config)

    srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
    if err != nil {
        log.Fatalf("Unable to retrieve Drive client: %v", err)
    }

    return srv
}

func getClient(config *oauth2.Config) *http.Client {
    tokFile := "token.json"
    tok, err := tokenFromFile(tokFile)
    if err != nil {
            tok = getTokenFromWeb(config)
            saveToken(tokFile, tok)
    }
    return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code %v", err)
        }

        tok, err := config.Exchange(context.TODO(), authCode)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web %v", err)
        }
        return tok
}

func saveToken(path string, token *oauth2.Token) {
        fmt.Printf("Saving credential file to: %s\n", path)
        f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
        if err != nil {
                log.Fatalf("Unable to cache oauth token: %v", err)
        }
        defer f.Close()
        json.NewEncoder(f).Encode(token)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
            return nil, err
    }
    defer f.Close()
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    return tok, err
}
