package ttt

import (
	"net/http"

	userTmp "github.com/firefirestyle/go.miniuser/template"
)

var userConfig = userTmp.UserTemplateConfig{
	TwitterConsumerKey:       "dummy",
	TwitterConsumerSecret:    "dummy",
	TwitterAccessToken:       "dummy",
	TwitterAccessTokenSecret: "dummy",
	FacebookAppSecret:        "dummy",
	FacebookAppId:            "dummy",
	GroupName:                "Main",
	KindBaseName:             "FFUser",
	AllowInvalidSSL:          true,
}

var userTemplate = userTmp.NewUserTemplate(userConfig)

func init() {
	userTemplate.InitUserApi()
	initHomepage()
}

func initHomepage() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to FireFireStyle!!"))
	})
}
