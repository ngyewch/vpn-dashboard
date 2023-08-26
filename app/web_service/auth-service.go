package web_service

import (
	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/github"
	"github.com/dghubble/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
	"net/http"
	"os"
	"strings"
)

const (
	sessionName     = "vpn-dashboard"
	sessionUserKey  = "github.id"
	sessionUsername = "github.username"

	authUriPrefix         = "login"
	authGithubLoginUri    = authUriPrefix + "/github"
	authGithubCallbackUri = authUriPrefix + "/oauth2/code/github"
	authUnauthorizedUri   = authUriPrefix + "/unauthorized"
)

var (
	allowedUsers = strings.Split(os.Getenv("ALLOWED_USERS"), ",")
)

type AuthService struct {
	cookieStore *sessions.CookieStore
}

func NewAuthService(cookieStore *sessions.CookieStore) *AuthService {
	return &AuthService{
		cookieStore: cookieStore,
	}
}

func (service *AuthService) Install(r *gin.Engine) {
	r.Use(service.authMiddleware())

	oauth2Config := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_OAUTH2_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_OAUTH2_REDIRECT_URL"),
		Endpoint:     githubOAuth2.Endpoint,
	}

	stateConfig := gologin.DebugOnlyCookieConfig

	r.Any(authGithubLoginUri,
		gin.WrapH(github.StateHandler(stateConfig, github.LoginHandler(oauth2Config, nil))))
	r.Any(authGithubCallbackUri,
		gin.WrapH(github.StateHandler(stateConfig, github.CallbackHandler(oauth2Config, service.issueSession(), nil))))

	r.GET(authUnauthorizedUri, func(c *gin.Context) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	})
}

func (service *AuthService) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/"+authUriPrefix+"/") {
			_, err := service.cookieStore.Get(c.Request, sessionName)
			if err != nil {
				c.Redirect(http.StatusFound, authGithubLoginUri)
				c.Abort()
				return
			}
		}
	}
}

func (service *AuthService) issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var match = false
		for _, userId := range allowedUsers {
			if userId == *githubUser.Login {
				match = true
				break
			}
		}

		if !match {
			http.Redirect(w, req, authUnauthorizedUri, http.StatusFound)
			return
		}

		session := service.cookieStore.New(sessionName)
		session.Values[sessionUserKey] = *githubUser.ID
		session.Values[sessionUsername] = *githubUser.Login
		err = session.Save(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/public/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
