package bridge

import (
	"git.codecoop.org/systemli/ticker/internal/model"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/go-twitter/twitter"
	"strconv"
)

var Twitter *TwitterBridge

type bridge interface {
	Initialized() bool
	Update(ticker model.Ticker, message model.Message) (string, error)
}

//
type TwitterBridge struct {
	ConsumerKey    string
	ConsumerSecret string
}

//
func NewTwitterBridge(key, secret string) *TwitterBridge {
	return &TwitterBridge{
		ConsumerKey:    key,
		ConsumerSecret: secret,
	}
}

//
func (tb *TwitterBridge) Initialized() bool {
	return tb.ConsumerKey != "" && tb.ConsumerSecret != ""
}

//
func (tb *TwitterBridge) Update(ticker model.Ticker, message model.Message) (string, error) {
	client := tb.client(ticker.Twitter.Token, ticker.Twitter.Secret)

	tweet, _, err := client.Statuses.Update(message.Text, nil)
	if err != nil {
		return "", err
	}
	id := strconv.FormatInt(tweet.ID, 10)

	return id, nil
}

//User returns the user information.
func (tb *TwitterBridge) User(ticker model.Ticker) (*twitter.User, error) {
	token := oauth1.NewToken(ticker.Twitter.Token, ticker.Twitter.Secret)
	httpClient := tb.config().Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	user, _, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{
		IncludeEmail:    twitter.Bool(false),
		IncludeEntities: twitter.Bool(false),
		SkipStatus:      twitter.Bool(true),
	})

	if err != nil {
		return user, err
	}

	return user, nil
}

func (tb *TwitterBridge) config() *oauth1.Config {
	return oauth1.NewConfig(tb.ConsumerKey, tb.ConsumerSecret)
}

func (tb *TwitterBridge) client(accessToken, accessSecret string) *twitter.Client {
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := tb.config().Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}