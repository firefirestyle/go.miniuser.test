package ttt

import (
	"net/http"

	"errors"

	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	"github.com/firefirestyle/go.minioauth/twitter"
	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"

	userhundler "github.com/firefirestyle/go.miniuser/handler"
	//
	"github.com/firefirestyle/go.minioauth/facebook"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	//

	"io/ioutil"
)

const (
	UrlTwitterTokenUrlRedirect  = "/api/v1/twitter/tokenurl/redirect"
	UrlTwitterTokenCallback     = "/api/v1/twitter/tokenurl/callback"
	UrlFacebookTokenUrlRedirect = "/api/v1/facebook/tokenurl/redirect"
	UrlFacebookTokenCallback    = "/api/v1/facebook/tokenurl/callback"
	UrlUserGet                  = "/api/v1/user/get"
	UrlUserFind                 = "/api/v1/user/find"
	UrlUserBlobGet              = "/api/v1/user/getblob"
	UrlUserRequestBlobUrl       = "/api/v1/user/requestbloburl"
	UrlUserCallbackBlobUrl      = "/api/v1/user/callbackbloburl"
	UrlMeLogout                 = "/api/v1/me/logout"
)

var twitterHandlerObj *twitter.TwitterHandler = nil
var userHandlerObj *userhundler.UserHandler = nil

func CheckLogin(r *http.Request, input *miniprop.MiniProp) minisession.CheckLoginIdResult {
	ctx := appengine.NewContext(r)
	token := input.GetString("token", "")
	return GetUserHundlerObj(ctx).GetSessionMgr().CheckLoginId(ctx, token, minisession.MakeAccessTokenConfigFromRequest(r))
}

func GetUserHundlerObj(ctx context.Context) *userhundler.UserHandler {
	if userHandlerObj == nil {
		v := appengine.DefaultVersionHostname(ctx)
		if v == "127.0.0.1:8080" {
			v = "localhost:8080"
		}
		userHandlerObj = userhundler.NewUserHandler(UrlUserCallbackBlobUrl,
			userhundler.UserHandlerManagerConfig{ //
				ProjectId:   "firefirestyle",
				UserKind:    "user",
				RelayIdKind: "relayId",
				SessionKind: "session",
			}, //
			twitter.TwitterOAuthConfig{
				ConsumerKey:       TwitterConsumerKey,
				ConsumerSecret:    TwitterConsumerSecret,
				AccessToken:       TwitterAccessToken,
				AccessTokenSecret: TwitterAccessTokenSecret,
				CallbackUrl:       "http://" + appengine.DefaultVersionHostname(ctx) + "" + UrlTwitterTokenCallback,
				SecretSign:        appengine.VersionID(ctx),
			}, //
			facebook.FacebookOAuthConfig{
				ConfigFacebookAppSecret: ConfigFacebookAppSecret,
				ConfigFacebookAppId:     ConfigFacebookAppId,
				SecretSign:              appengine.VersionID(ctx),
				CallbackUrl:             "http://" + v + "" + UrlFacebookTokenCallback,
			},
			userhundler.UserHandlerOnEvent{}, //
			blobhandler.BlobHandlerOnEvent{
				OnBlobRequest: func(w http.ResponseWriter, r *http.Request, input *miniprop.MiniProp, output *miniprop.MiniProp, h *blobhandler.BlobHandler) (string, map[string]string, error) {
					ret := CheckLogin(r, input)
					if ret.IsLogin == false {
						return "", map[string]string{}, errors.New("Failed in token check")
					}
					return ret.AccessTokenObj.GetLoginId(), map[string]string{}, nil
				},
			})
	}
	return userHandlerObj
}

func init() {
	initApi()
	initHomepage()
}

func initHomepage() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to FireFireStyle!!"))
	})
}

func initApi() {
	// twitter
	http.HandleFunc(UrlTwitterTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleTwitterRequestToken(w, r)
	})
	http.HandleFunc(UrlTwitterTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleTwitterCallbackToken(w, r)
	})
	// facebook
	http.HandleFunc(UrlFacebookTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleFacebookRequestToken(w, r)
	})
	http.HandleFunc(UrlFacebookTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleFacebookCallbackToken(w, r)
	})
	// user
	http.HandleFunc(UrlUserGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleGet(w, r)
	})
	http.HandleFunc(UrlUserFind, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleFind(w, r)
	})
	http.HandleFunc(UrlUserRequestBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleBlobRequestToken(w, r)
	})
	http.HandleFunc(UrlUserCallbackBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleBlobUpdated(w, r)
	})
	http.HandleFunc(UrlUserBlobGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleBlobGet(w, r)
	})

	// me
	http.HandleFunc(UrlMeLogout, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		propObj := miniprop.NewMiniPropFromJson(bodyBytes)
		token := propObj.GetString("token", "")
		ctx := appengine.NewContext(r)
		GetUserHundlerObj(ctx).GetSessionMgr().Logout(ctx, token, minisession.MakeAccessTokenConfigFromRequest(r))
	})

}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
