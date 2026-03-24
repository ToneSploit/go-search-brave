package functions

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const pushoverAPI = "https://api.pushover.net/1/messages.json"

// PushoverConfig holds the credentials and enabled flag for Pushover notifications.
type PushoverConfig struct {
	AppToken string
	UserKey  string
	Enabled  bool
}

// SendPushoverNotification sends a single notification for a news article.
// The article title becomes the notification title, and the article URL is attached.
func SendPushoverNotification(cfg *PushoverConfig, article *NewsArticle) error {
	if !cfg.Enabled {
		return nil
	}

	body := article.Description
	if body == "" {
		body = article.Title
	}

	form := url.Values{}
	form.Set("token", cfg.AppToken)
	form.Set("user", cfg.UserKey)
	form.Set("title", article.Title)
	form.Set("message", body)
	form.Set("url", article.URL)
	form.Set("url_title", "Read article")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(pushoverAPI, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("pushover request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pushover returned unexpected status: %s", resp.Status)
	}

	return nil
}
