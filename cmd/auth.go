package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anton-johansson/k8s-login/kubernetes"
	"github.com/coreos/go-oidc"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type loginApp struct {
	// Parameters
	clientID           string
	clientSecret       string
	port               int
	kubeconfigFileName string
	updateContext      bool
	shutdownChannel    chan bool

	// State
	kubeconfig *kubernetes.KubeConfig
	client     *http.Client
	provider   *oidc.Provider
	verifier   *oidc.IDTokenVerifier
}

var app loginApp
var authCommand = &cobra.Command{
	Use:   "auth [server]",
	Short: "Attempts to authenticate to the given server",
	Args: func(command *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a server")
		}
		if len(args) > 1 {
			return errors.New("too many arguments")
		}
		return nil
	},
	RunE: func(command *cobra.Command, args []string) error {
		kubeconfig, error := kubernetes.GetKubeConfig(app.kubeconfigFileName)
		if error != nil {
			return error
		}
		app.kubeconfig = kubeconfig

		serverName := args[0]
		server, error := getServerByName(kubeconfig, serverName)
		if error != nil {
			return error
		}
		location := convertServerAddressToDex(server)

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		context := oidc.ClientContext(context.Background(), client)
		provider, error := oidc.NewProvider(context, location)
		if error != nil {
			return error
		}

		app.client = client
		app.provider = provider
		app.verifier = provider.Verifier(&oidc.Config{ClientID: app.clientID})
		app.shutdownChannel = make(chan bool)

		http.HandleFunc("/", app.initiateLogin)
		http.HandleFunc("/callback", app.finalizeLogin)

		open("http://localhost:" + strconv.Itoa(app.port))
		go awaitShutdown(app.shutdownChannel)
		return http.ListenAndServe("localhost:"+strconv.Itoa(app.port), nil)
	},
}

func awaitShutdown(shutdownChannel chan bool) {
	<-shutdownChannel
	os.Exit(0)
}

func (app *loginApp) oauth() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     app.clientID,
		ClientSecret: app.clientSecret,
		Endpoint:     app.provider.Endpoint(),
		RedirectURL:  "http://localhost:" + strconv.Itoa(app.port) + "/callback",
		Scopes:       []string{"groups", "openid", "profile", "email", "offline_access"},
	}
}

func (app *loginApp) initiateLogin(writer http.ResponseWriter, request *http.Request) {
	authCodeURL := app.oauth().AuthCodeURL("login")
	http.Redirect(writer, request, authCodeURL, http.StatusSeeOther)
}

func (app *loginApp) finalizeLogin(writer http.ResponseWriter, request *http.Request) {
	context := oidc.ClientContext(request.Context(), app.client)
	oauth := app.oauth()

	if errorMessage := request.FormValue("error"); errorMessage != "" {
		http.Error(writer, errorMessage+": "+request.FormValue("error_description"), http.StatusBadRequest)
		return
	}
	code := request.FormValue("code")
	if code == "" {
		http.Error(writer, "No code in the request", http.StatusBadRequest)
		return
	}
	if state := request.FormValue("state"); state != "login" {
		http.Error(writer, "Incorrect client application", http.StatusBadRequest)
		return
	}
	token, error := oauth.Exchange(context, code)
	if error != nil {
		fmt.Println(error)
		http.Error(writer, "Error when exchanging", http.StatusBadRequest)
		return
	}

	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(writer, "No ID token in the token response", http.StatusInternalServerError)
		return
	}

	idToken, error := app.verifier.Verify(context, rawIdToken)
	if error != nil {
		fmt.Println(error)
		http.Error(writer, "Failed to verify token", http.StatusInternalServerError)
		return
	}

	var claims json.RawMessage
	if error := idToken.Claims(&claims); error != nil {
		fmt.Println(error)
		http.Error(writer, "Could not get claims for token", http.StatusInternalServerError)
		return
	}

	buffer := new(bytes.Buffer)
	if error := json.Indent(buffer, []byte(claims), "", "  "); error != nil {
		fmt.Println(error)
		http.Error(writer, "Could not parse JSON from claim", http.StatusInternalServerError)
		return
	}

	var unmarshalledClaims map[string]interface{}
	if error := json.Unmarshal(claims, &unmarshalledClaims); error != nil {
		fmt.Println(error)
		http.Error(writer, "Could not parse JSON from claim", http.StatusInternalServerError)
		return
	}

	user := kubernetes.UserData{
		Name:         unmarshalledClaims["name"].(string),
		ClientID:     app.clientID,
		ClientSecret: app.clientSecret,
		IDToken:      rawIdToken,
		RefreshToken: token.RefreshToken,
		IssuerURL:    unmarshalledClaims["iss"].(string),
	}

	kubernetes.UpdateKubeConfig(app.kubeconfig, user, app.updateContext)

	message := "Successfully updated '" + app.kubeconfig.FileName + "'"
	fmt.Println(message)
	writer.Write([]byte(message + ". You can now close this browser tab."))
	app.shutdownChannel <- true
}

func getServerByName(kubeconfig *kubernetes.KubeConfig, serverName string) (kubernetes.Server, error) {
	servers := kubernetes.GetServers(kubeconfig)
	for _, server := range servers {
		if server.Name == serverName {
			return server, nil
		}
	}
	return kubernetes.Server{}, errors.New("No server with name '" + serverName + "'")
}

func convertServerAddressToDex(server kubernetes.Server) string {
	address := server.Address
	address = strings.Replace(address, "://k8s.svc", "://dex.svc", 1)
	address = strings.Replace(address, ":6443", "", 1)
	return address
}

func open(location string) error {
	command, arguments := (func() (string, []string) {
		switch runtime.GOOS {
		case "windows":
			return "cmd", []string{"/c", "start"}
		case "darwin":
			return "open", []string{}
		default:
			return "xdg-open", []string{}
		}
	})()

	arguments = append(arguments, location)
	return exec.Command(command, arguments...).Start()
}

func init() {
	authCommand.Flags().StringVarP(&app.clientID, "client-id", "i", "k8s-login", "The OAuth2 client ID of this application")
	authCommand.Flags().StringVarP(&app.clientSecret, "client-secret", "s", "lhHN7keNTf4MXEIH3WF4NUL701qITv9Q", "The OAuth2 client secret of this application")
	authCommand.Flags().IntVarP(&app.port, "local-port", "p", 5555, "The local port to host the temporary web server on")
	authCommand.Flags().StringVarP(&app.kubeconfigFileName, "kubeconfig", "k", kubernetes.GetDefaultKubeConfigFileName(), "The path to the kubeconfig")
	authCommand.Flags().BoolVarP(&app.updateContext, "update-context", "c", true, "Indicates whether or not to update the actual context of the kubeconfig")
	rootCommand.AddCommand(authCommand)
}
