// HTTPS-to-SMTP relay for examples/extensions/due-dates.
// Uses SMTP from the environment. Does not read Ordryn Admin → Email.
// The same listener can also receive posts from the full Email relay example.
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

type payload struct {
	Text    string `json:"text"`
	Content string `json:"content"`
	Event   string `json:"event"`
	Name    string `json:"name"`
	Project string `json:"project"`
	DueDate string `json:"due_date"`
}

func main() {
	listen := getenv("LISTEN", "127.0.0.1:8787")
	_ = mustEnv("SMTP_HOST")
	_ = mustEnv("SMTP_USER")
	_ = mustEnv("SMTP_PASS")
	_ = mustEnv("MAIL_FROM")
	http.HandleFunc("/", handle)
	log.Printf("due-dates relay listening on %s (MAIL_TO=%s)", listen, mustEnv("MAIL_TO"))
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
	text := strings.TrimSpace(p.Text)
	if text == "" {
		text = strings.TrimSpace(p.Content)
	}
	if text == "" {
		http.Error(w, "empty message", http.StatusBadRequest)
		return
	}
	if err := sendSMTP(subjectFor(p), text); err != nil {
		log.Printf("smtp: %v", err)
		http.Error(w, "send failed", http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func subjectFor(p payload) string {
	label := strings.TrimSpace(p.Project)
	if p.Name != "" {
		if label != "" {
			label += ": "
		}
		label += p.Name
	}
	switch p.Event {
	case "task.overdue":
		if label == "" {
			return "Overdue task"
		}
		return "Overdue: " + label
	case "task.due_changed":
		due := strings.TrimSpace(p.DueDate)
		if due == "" {
			if label == "" {
				return "Due date cleared"
			}
			return "Due date cleared: " + label
		}
		if label == "" {
			return "Due " + due
		}
		return "Due " + due + ": " + label
	}
	if label != "" {
		return label
	}
	if p.Event != "" {
		return p.Event
	}
	return "Ordryn"
}

func validSignature(secret string, body []byte, header string) bool {
	header = strings.TrimSpace(header)
	wantMac := hmac.New(sha256.New, []byte(secret))
	_, _ = wantMac.Write(body)
	want := "sha256=" + hex.EncodeToString(wantMac.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(want), []byte(header)) == 1
}

func sendSMTP(subject, body string) error {
	host := mustEnv("SMTP_HOST")
	port := getenv("SMTP_PORT", "587")
	user := mustEnv("SMTP_USER")
	pass := mustEnv("SMTP_PASS")
	from := mustEnv("MAIL_FROM")
	to := mustEnv("MAIL_TO")
	addr := net.JoinHostPort(host, port)

	html := strings.ReplaceAll(body, "\n", "<br/>")
	msg := strings.Join([]string{
		"From: " + from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		html,
	}, "\r\n")

	auth := smtp.PlainAuth("", user, pass, host)
	if port == "465" {
		tlsCfg := &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}
		conn, err := tls.Dial("tcp", addr, tlsCfg)
		if err != nil {
			return err
		}
		defer conn.Close()
		c, err := smtp.NewClient(conn, host)
		if err != nil {
			return err
		}
		defer c.Close()
		if err := c.Auth(auth); err != nil {
			return err
		}
		if err := c.Mail(from); err != nil {
			return err
		}
		if err := c.Rcpt(to); err != nil {
			return err
		}
		wc, err := c.Data()
		if err != nil {
			return err
		}
		if _, err := wc.Write([]byte(msg)); err != nil {
			_ = wc.Close()
			return err
		}
		if err := wc.Close(); err != nil {
			return err
		}
		return c.Quit()
	}
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}

func mustEnv(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		log.Fatalf("missing %s", key)
	}
	return v
}

func getenv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
