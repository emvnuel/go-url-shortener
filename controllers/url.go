package controllers

import (
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/emvunel/urlshortener/models"
	"github.com/emvunel/urlshortener/utils"
	"github.com/gin-gonic/gin"
)

var regexAlphanum = regexp.MustCompile("^[a-zA-Z0-9-]+$")

type UrlRequest struct {
	OriginalUrl     string `json:"originalUrl" binding:"required,url"`
	Autogenerate    bool   `json:"autogenerate"`
	ShortenedSuffix string `json:"shortenedSuffix"`
}

func ShortenUrl(c *gin.Context) {
	var input UrlRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Autogenerate {
		input.ShortenedSuffix = utils.RandStringBytes(6)
	}

	if input.ShortenedSuffix == "" || !regexAlphanum.Match([]byte(input.ShortenedSuffix)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "suffix can not be empty, contains spaces or special simbols"})
		return
	}

	shortenedUrl := models.ShortenedUrl{ID: input.ShortenedSuffix, OriginalUrl: input.OriginalUrl}
	if err := models.DB.Create(&shortenedUrl).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "suffix already exists"})
		return
	}

	c.Header("Location", "http://"+c.Request.Host+"/"+shortenedUrl.ID)
	c.Status(http.StatusCreated)
}

func RedirectUrl(c *gin.Context) {
	var shortenedUrl models.ShortenedUrl
	if err := models.DB.Where("id = ?", c.Param("id")).First(&shortenedUrl).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
		return
	}

	go func() {

		requestDump, err := httputil.DumpRequest(c.Request, true)
		dump := string(requestDump)

		if err != nil {
			panic("err " + err.Error())
		}

		click := models.Click{ShortenedUrl: shortenedUrl, RawRequest: dump}

		models.DB.Create(&click)

	}()

	c.Header("Location", shortenedUrl.OriginalUrl)
	c.Status(http.StatusMovedPermanently)
}
