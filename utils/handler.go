package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if _, ok := UrlToLink[id]; !ok {
			ctx.Data(http.StatusOK, ContentTypeHTML, []byte("<html><h1>Not Valid Url</h1></html>"))
			return
		}
		baseUrl := "https://objectstorage.ap-mumbai-1.oraclecloud.com"
		pdfPath := UrlToLink[id]
		fullUrl := baseUrl + pdfPath
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Url": fullUrl,
		})
	})

	r.Run()
}
