# Generic webhook

Posts task events as JSON to any **public HTTPS** webhook. Use this for n8n, Zapier, Mattermost, or a listener you write. For Slack, Discord, Teams, or Google Chat, the dedicated example extensions send the payload those services expect and restrict the URL host.

Core is not patched; this folder is a drop-in extension. Each **project owner** configures a team URL; members can add their own destination. Personal inbox tasks (no project) use Settings → Integrations. Loopback, private, and link-local destinations are rejected.

This manifest uses `"delivery": { "type": "http.webhook", "format": "json" }`. Authors writing their own extension can instead set `format` to `text` (`{"text": "…"}`, Slack-style) or `content` (`{"content": "…"}`, Discord-style).

Outbound requests include `X-Ordryn-Signature: sha256=<hex>` when a signing secret is rotated on the Extensions panel. The HMAC is over the raw JSON body.

## Install

1. Copy this directory to `data/extensions/webhook` on the server (the folder name must stay `webhook`).
2. Restart Ordryn.
3. Site admin: Admin → Extensions → Generic webhook → Enable, then Save.
4. Project owner: open the project → Settings → Extensions → Generic webhook. Paste the HTTPS URL, choose triggers, edit messages if you want, Enable, Save. Rotate the signing secret and copy it once.
5. Optional: set `PUBLIC_URL` in `.env` (for example `https://todo.example.com`) so `{url}` in templates becomes a clickable task link.

## JSON body (`format: json`)

`POST` with `Content-Type: application/json`:

```json
{
  "text": "Task Ship updated to Done in project Ordryn",
  "content": "Task Ship updated to Done in project Ordryn",
  "event": "task.updated",
  "event_id": "00000000-0000-0000-0000-000000000001",
  "occurred_at": "2026-09-15T16:00:00Z",
  "changed": ["status"],
  "id": "42",
  "name": "Ship",
  "task": "Ship",
  "status": "Done",
  "old_status": "In progress",
  "project": "Ordryn",
  "actor": "ada",
  "url": "https://todo.example.com/tasks/42",
  "priority": "High",
  "comment": "",
  "claimed_by": "ada",
  "due_date": "2026-09-20",
  "sprint": "Sprint 1",
  "tags": "release",
  "fields": {"severity.level": "high"}
}
```

Templates may use `{task}` `{name}` `{status}` `{old_status}` `{project}` `{actor}` `{url}` `{id}` `{priority}` `{comment}` `{claimed_by}` `{due_date}` `{sprint}` `{tags}` `{mentions}` `{member}` `{count}` `{event_id}`. Leave a template blank to skip that trigger.

`join.request`, `join.approved`, and `join.denied` are site-level only (Admin → Extensions site URL). They are never posted into every project channel.
