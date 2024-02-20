package function

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	appCreds, err := readSecret("wordleboard-client-id")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON([]byte(appCreds), gmail.GmailReadonlyScope)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	client, err := getClient(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	user := "me"
	pubsubTopic, err := readSecret("wordleboard-pubsub-topic")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}
	watch, err := srv.Users.Watch(user, &gmail.WatchRequest{
		LabelIds:  []string{"INBOX"},
		TopicName: pubsubTopic,
	}).Do()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Watch: %v\n", watch)))
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	tok, err := tokenFromSecret("wordleboard-token")
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), tok), err
}

func tokenFromSecret(name string) (*oauth2.Token, error) {
	data, err := readSecret(name)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.Unmarshal([]byte(data), tok)
	return tok, err
}

// readSecret reads a file from /var/lib/faasd-provider/secrets/<name> or
// returns an error
func readSecret(name string) (string, error) {
	data, err := os.ReadFile("/var/openfaas/secrets/" + name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
