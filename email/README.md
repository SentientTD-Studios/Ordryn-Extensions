# Email relay

Posts task events as JSON to **your** public HTTPS relay. The relay sends mail using SMTP or an email API **you** configure.

Admin → Email (site SMTP/Mailgun) is never used. Password resets and invites stay on that core mailer, which is rate-limited. Project members cannot send through it.

Private, loopback, and link-local relay URLs are rejected.

## Install

1. Copy this directory to `data/extensions/email` (the folder name must stay `email`).
2. Restart Ordryn.
3. Site admin: Admin → Extensions → Email relay → Enable, then Save.
4. Run a relay you control (see `relay/main.go` in this folder, or any HTTPS listener that accepts the JSON body).
5. Project owner: project → Settings → Extensions → Email relay. Paste the public HTTPS URL, choose triggers, Enable, Save.
6. Optional: rotate the signing secret on that form and set `ORDRYN_SIGNING_SECRET` on the relay so it can verify `X-Ordryn-Signature`.

## Sample relay

`relay/main.go` is a tiny HTTPS-to-SMTP forwarder (not part of the Ordryn binary):

```bash
export SMTP_HOST=smtp.example.com
export SMTP_PORT=587
export SMTP_USER=relay@example.com
export SMTP_PASS=secret
export MAIL_FROM=relay@example.com
export MAIL_TO=alerts@example.com
export ORDRYN_SIGNING_SECRET='the secret shown once after Rotate signing secret'
go run ./relay
```

Put that process behind public HTTPS (Caddy, nginx, Cloudflare Tunnel, or similar). Ordryn will not POST to `localhost`.

The JSON body matches the generic webhook: `text`, `content`, `event`, `name`, `project`, `actor`, `url`, and related fields. Use `text` as the email body.

For due dates only, copy the `due-dates` folder instead — same custom-relay pattern, fewer events.

Sibling examples cover Discord, Slack, Teams, a generic HTTPS webhook, ntfy, due-dates, comments, claimed, activity, lifecycle, mentions, callback-bot, and join-requests.
