package main
import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    router.POST("/upload", upload)

    router.Run("localhost:8080")
}

func upload(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, "upload starting")
}
