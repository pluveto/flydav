package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pluveto/flydav/internal/config"
	"github.com/pluveto/flydav/internal/logger"
	"github.com/pluveto/flydav/pkg/authenticator"
)

type AuthModule struct {
	Config        config.AuthConfig
	Authenticator authenticator.Authenticator
	// 可能还有其他的状态或依赖
}

func NewAuthModule(cfg config.AuthConfig) *AuthModule {
	var auth authenticator.Authenticator
	if cfg.Backends.Static.Enabled {
		auth = authenticator.NewStaticAuthenticator(cfg.Backends.Static.Users)
	}
	// 这里可以添加其他认证方法的初始化
	return &AuthModule{
		Config:        cfg,
		Authenticator: auth,
	}
}

func (as *AuthModule) Start() error {
	logger.Info("starting auth module")
	return nil
}

func (as *AuthModule) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", as.handleLogin).Methods("POST")
	router.HandleFunc("/internal/authenticate", as.handleAuthenticate).Methods("GET")
}

func (as *AuthModule) IsAuthenticated(r *http.Request) bool {
	switch as.Config.GetEnabledAuthBackend() {
	case config.StaticAuthBackend:
		return as.isAuthenticatedStatic(r)
	default:
		return false
	}
}

func (as *AuthModule) isAuthenticatedStatic(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	payload, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		return false
	}

	username := pair[0]
	password := pair[1]

	ok, err := as.Authenticator.Authenticate(username, password)

	if err != nil {
		logger.Error("Error authenticating user: ", err)
		return false
	}

	return ok
}

func (as *AuthModule) handleLogin(w http.ResponseWriter, r *http.Request) {
	switch as.Config.GetEnabledAuthBackend() {
	case config.StaticAuthBackend:
		as.handleAuthenticate(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (as *AuthModule) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username := pair[0]
	path := r.URL.Query().Get("path")
	permission := config.Permission(r.URL.Query().Get("permission"))

	/**
	ReadPermission
	WritePermission
	*/
	ok, err := as.Authenticator.Authorize(username, path, permission)
	if err != nil {
		logger.Error("Error authenticating user: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
