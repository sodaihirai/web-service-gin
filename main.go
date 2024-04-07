package main
import (
    "context"
    "os"
    "net/http"
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
        service.Files.Create(&drive.File{Name: file.Name()})
    }
    c.IndentedJSON(http.StatusOK, "upload starting")
}

func getDriveService() *drive.Service {
    ctx := context.Background()
    b, err := os.ReadFile("credentials.json")
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }

    // If modifying these scopes, delete your previously saved token.json.
    config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/drive.install")
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
    tok := tokenFromFile(tokFile)
    return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token) {
    f, _ := os.Open(file)
    defer f.Close()
    tok := &oauth2.Token{}
    return tok
}
