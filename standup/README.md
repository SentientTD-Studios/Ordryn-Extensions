# Standup

Example of the **sandboxed HTML panel** (`ui` in the manifest). Project members see `panel.html` on the project Extensions tab in an iframe.

This iframe cannot read the Ordryn session, cookies, or callback token (`sandbox` without `allow-same-origin`). The panel is a daily check-in worksheet: fill done / today / blockers, copy markdown into a task comment, or copy a `curl` that posts through `POST /api/v2/ext/callback`.

Callback-bot is the kitchen-sink relay demo. This folder is the UI-first example.

## Install

1. Copy this directory to `data/extensions/standup` (folder name must stay `standup`).
2. Restart Ordryn. Logs should show `extensions: loaded standup`.
3. Site admin: Admin → Extensions → Standup → Enable, then Save.
4. Project owner: project → Settings → Extensions → Standup. Enable for the project. Rotate the callback token if you will post from a terminal or relay (copy it once). Save.

Then open the **Extension panel** on that same tab. Edit a task to see `standup.last` (list + sidebar) and `standup.notes` (sidebar).

## What the panel can and cannot do

The host serves the HTML with a tight CSP (inline script/style only, `connect-src https:`, no `'self'`). Combined with the unique iframe origin, the document **cannot** call `/api/v2` with the signed-in session.

Use the panel to compose text. To write into Ordryn:

- Paste the markdown into a task comment yourself, or
- Rotate the project callback token (shown once) and run the generated `curl` against `https://your-host/api/v2/ext/callback`, or
- Enable the project inbound webhook; with this extension on, `POST /api/v2/webhooks/inbound` accepts `comment` and `set_field` (`project_id` required).

Do not paste long-lived secrets into the panel; it has nowhere safe to store them.

## Callback API

```json
POST /api/v2/ext/callback
Authorization: Bearer <callback_token>
{
  "action": "comment" | "set_field" | "get",
  "task_id": 42,
  "comment": "optional",
  "field": "standup.last",
  "value": "2026-09-17"
}
```

The panel stamps `standup.last` (today) and `standup.notes` (blockers, or the full check-in). The token is project-scoped, not a user API key.
