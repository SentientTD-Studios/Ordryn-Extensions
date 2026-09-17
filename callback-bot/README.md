# Callback bot

Drop-in example of the extension surface that is more than a chat webhook: identity + icon, a sandboxed HTML panel, `select` / `status` / `user` settings (sent as JSON `config`), `date` / `markdown` fields, scoped callback tokens, and inbound `complete` / `set_field`.

This iframe cannot read the Ordryn session. Relays use the `callback_token` on outbound JSON instead.

## Install

1. Copy this directory to `data/extensions/callback-bot` (folder name must stay `callback-bot`).
2. Restart Ordryn.
3. Site admin: Admin → Extensions → Callback bot → Enable, then Save.
4. Project owner: project → Settings → Extensions → Callback bot. Enable for the project, paste a public HTTPS relay URL, rotate the callback token (copy it once), optionally rotate signing, Save.
5. Run the sample relay (below) behind public HTTPS. Ordryn will not POST to `localhost`.

Optional: turn on the project's inbound webhook. With this extension enabled, `POST /api/v2/webhooks/inbound` also accepts `complete` and `set_field` (plus the usual `create` / `comment` flags). `/api/v1` remains a compatibility alias.

## Sample relay

`relay/main.go` receives the JSON event, then calls `POST /api/v2/ext/callback` with the token from the payload. On `task.due_soon` it posts a comment. On `task.created` it sets `callback-bot.review_by`.

```bash
export ORDRYN_SIGNING_SECRET='the secret shown once after Rotate signing secret'
export LISTEN=127.0.0.1:8790
go run ./relay
```

If `config.alert_mode` is `urgent`, the relay skips events unless `priority` is High. `config.ping_user` is included in comments when set.

## Callback API

```json
POST /api/v2/ext/callback
Authorization: Bearer <callback_token>
{
  "action": "get" | "complete" | "comment" | "set_field",
  "task_id": 42,
  "comment": "optional",
  "field": "callback-bot.notes",
  "value": "optional"
}
```

The token is project-scoped, not a user API key. Rotate it from the Extensions tab.
