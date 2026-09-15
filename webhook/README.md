# Generic webhook

Posts task events as JSON to any **public HTTPS** webhook. Use this for n8n, Zapier, Mattermost, Google Chat, or a listener you write. For Slack, Discord, or Teams, the dedicated example extensions send the payload those services expect and restrict the URL host.

Core is not patched; this folder is a drop-in extension. Each **project owner** configures their own URL. Personal tasks (no project) are never posted. Loopback, private, and link-local destinations are rejected.

This manifest uses `"delivery": { "type": "http.webhook", "format": "json" }`. Authors writing their own extension can instead set `format` to `text` (`{"text": "…"}`, Slack-style) or `content` (`{"content": "…"}`, Discord-style).

## Install

1. Copy this directory to `data/extensions/webhook` on the server (the folder name must stay `webhook`).
2. Restart Ordryn.
3. Site admin: Admin → Extensions → Generic webhook → Enable, then Save.
4. Project owner: open the project → Settings → Extensions → Generic webhook. Paste the HTTPS URL, choose triggers, edit messages if you want, Enable, Save.
5. Optional: set `PUBLIC_URL` in `.env` (for example `https://todo.example.com`) so `{url}` in templates becomes a clickable task link.

## JSON body (`format: json`)

`POST` with `Content-Type: application/json`:

```json
{
  "text": "Task Ship updated to Done in project Ordryn",
  "content": "Task Ship updated to Done in project Ordryn",
  "event": "task.updated",
  "id": "42",
  "name": "Ship",
  "task": "Ship",
  "status": "Done",
  "old_status": "In progress",
  "project": "Ordryn",
  "actor": "ada",
  "url": "https://todo.example.com/tasks/42",
  "priority": "High"
}
```

Templates may use `{task}` `{name}` `{status}` `{old_status}` `{project}` `{actor}` `{url}` `{id}` `{priority}`. Leave a template blank to skip that trigger.
