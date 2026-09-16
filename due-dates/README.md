# Due dates

Posts `task.due_changed` and `task.overdue` as JSON to **your** public HTTPS relay. The relay sends mail using SMTP or an email API **you** configure. Use this when the full Email relay extension is too noisy.

Admin → Email (site SMTP/Mailgun) is never used. Password resets and invites stay on that core mailer, which is rate-limited. Project members cannot send through it.

Private, loopback, and link-local relay URLs are rejected. Overdue fires once per task per due date. Empty `{due_date}` means the date was cleared.

## Install

1. Copy this directory to `data/extensions/due-dates` (the folder name must stay `due-dates`).
2. Restart Ordryn.
3. Site admin: Admin → Extensions → Due dates → Enable, then Save.
4. Run a relay you control (see `relay/main.go` in this folder, or any HTTPS listener that accepts the JSON body).
5. Project owner: project → Settings → Extensions → Due dates. Paste the public HTTPS URL, choose triggers, Enable, Save.
6. Optional: rotate the signing secret on that form and set `ORDRYN_SIGNING_SECRET` on the relay so it can verify `X-Ordryn-Signature`.

## Sample relay

`relay/main.go` is a tiny HTTPS-to-SMTP forwarder (not part of the Ordryn binary). It reads the JSON `text` field as the email body and sets a subject from `event`, `project`, and `name`:

```bash
export SMTP_HOST=smtp.example.com
export SMTP_PORT=587
export SMTP_USER=relay@example.com
export SMTP_PASS=secret
export MAIL_FROM=relay@example.com
export MAIL_TO=alerts@example.com
export ORDRYN_SIGNING_SECRET='the secret shown once after Rotate signing secret'
go run ./examples/extensions/due-dates/relay
```

If you already copied this folder into `data/extensions/due-dates`, run `go run ./relay` from that directory instead.

Put that process behind public HTTPS (Caddy, nginx, Cloudflare Tunnel, or similar). Ordryn will not POST to `localhost`.

The JSON body matches the generic webhook: `text`, `content`, `event`, `name`, `project`, `actor`, `url`, `due_date`, and related fields. Use `text` as the email body.

You can point Due dates and the full Email relay at the same listener; this sample understands both.

Sibling examples cover Discord, Slack, Teams, a generic HTTPS webhook, ntfy, email relay, comments, claimed, activity, and join-requests.
