// HTTPS relay for examples/extensions/callback-bot.
// Receives Ordryn JSON webhooks, then calls POST /api/v1/ext/callback
// with the project-scoped token from the payload.
package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type payload struct {
	Event         string            `json:"event"`
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Priority      string            `json:"priority"`
	DueDate       string            `json:"due_date"`
	CallbackToken string            `json:"callback_token"`
	CallbackURL   string            `json:"callback_url"`
	Config        map[string]string `json:"config"`
}

func main() {
	listen := getenv("LISTEN", "127.0.0.1:8790")
	http.HandleFunc("/", handle)
	log.Printf("callback-bot relay listening on %s", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<16))
	if err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	if secret := strings.TrimSpace(os.Getenv("ORDRYN_SIGNING_SECRET")); secret != "" {
		if !validSignature(secret, body, r.Header.Get("X-Ordryn-Signature")) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
	}
	var p payload
	if err := json.Unmarshal(body, &p); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if p.Config["alert_mode"] == "urgent" && !strings.EqualFold(strings.TrimSpace(p.Priority), "High") {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	taskID, err := strconv.Atoi(strings.TrimSpace(p.ID))
	if err != nil || taskID <= 0 {
		http.Error(w, "missing task id", http.StatusBadRequest)
		return
	}
	token := strings.TrimSpace(p.CallbackToken)
	callbackURL := strings.TrimSpace(p.CallbackURL)
	if token == "" || callbackURL == "" {
		http.Error(w, "missing callback credentials", http.StatusBadRequest)
		return
	}
	switch p.Event {
	case "task.due_soon", "task.mentioned":
		comment := "callback-bot: " + strings.TrimSpace(p.Event)
		if name := strings.TrimSpace(p.Name); name != "" {
			comment += " on " + name
		}
		if uid := strings.TrimSpace(p.Config["ping_user"]); uid != "" {
			comment += " (ping user id " + uid + ")"
		}
		if err := callback(callbackURL, token, map[string]any{
			"action":  "comment",
			"task_id": taskID,
			"comment": comment,
		}); err != nil {
			log.Printf("callback comment: %v", err)
			http.Error(w, "callback failed", http.StatusBadGateway)
			return
		}
	case "task.created":
		review := strings.TrimSpace(p.DueDate)
		if review == "" {
			review = time.Now().UTC().Add(7 * 24 * time.Hour).Format("2006-01-02")
		}
		if err := callback(callbackURL, token, map[string]any{
			"action":  "set_field",
			"task_id": taskID,
			"field":   "callback-bot.review_by",
			"value":   review,
		}); err != nil {
			log.Printf("callback set_field: %v", err)
			http.Error(w, "callback failed", http.StatusBadGateway)
			return
		}
	default:
		// Acknowledge other subscribed events without writing back.
	}
	w.WriteHeader(http.StatusNoContent)
}

func callback(url, token string, body map[string]any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 1<<14))
	if res.StatusCode >= 300 {
		return fmt.Errorf("status %d", res.StatusCode)
	}
	return nil
}

func validSignature(secret string, body []byte, header string) bool {
	header = strings.TrimSpace(header)
	wantMac := hmac.New(sha256.New, []byte(secret))
	_, _ = wantMac.Write(body)
	want := "sha256=" + hex.EncodeToString(wantMac.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(want), []byte(header)) == 1
}

func getenv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
