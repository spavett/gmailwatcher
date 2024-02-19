package function

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		log.Fatalf("Unable to read client secret: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON([]byte(appCreds), gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	pubsubTopic, err := readSecret("wordleboard-pubsub-topic")
	if err != nil {
		log.Fatalf("Unable to read pubsub topic: %v", err)
	}
	watch, err := srv.Users.Watch(user, &gmail.WatchRequest{
		LabelIds:  []string{"INBOX"},
		TopicName: pubsubTopic,
	}).Do()
	if err != nil {
		log.Fatalf("Unable to watch: %v", err)
	}
	fmt.Printf("Watch: %v\n", watch)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromSecret("wordleboard-token")
	if err != nil {
		log.Fatalf("Unable to retrieve token from secret: %v", err)
	}
	return config.Client(context.Background(), tok)
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
	data, err := os.ReadFile("/var/lib/faasd-provider/secrets/" + name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
