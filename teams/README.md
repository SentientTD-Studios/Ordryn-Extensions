# Microsoft Teams notifications

Posts task events to a Teams channel using a Workflows incoming webhook (Adaptive Card with title and Open action). Core is not patched; this folder is a drop-in extension. Project owners configure a team channel; any member can add a personal “Notify me” destination.

Office 365 Connector incoming webhooks are retired; use a Teams **Workflow**.

## Install

1. Copy this directory to `data/extensions/teams` on the server (the folder name must stay `teams`).
2. Restart Ordryn.
3. Site admin: Admin → Extensions → Microsoft Teams → Enable, then Save.
4. In Teams: channel → Workflows → **Post to a channel when a webhook request is received**. Copy the HTTP URL (often `*.logic.azure.com` or `*.api.powerplatform.com`, and it may include a query string).
5. Project owner: open the project → Settings → Extensions → Microsoft Teams. Paste the webhook URL, choose triggers, edit messages if you want, Enable, Save.
6. Optional: set `PUBLIC_URL` in `.env` (for example `https://todo.example.com`) so `{url}` in templates becomes a clickable task link.

The site admin must enable the extension before any project can post. Channel, triggers, and templates live on the project, not in Admin → Extensions.

With **Only notify when status changes** on, `task.updated` is skipped unless the kanban/list status actually changed.

Templates may use `{task}` `{name}` `{status}` `{old_status}` `{project}` `{actor}` `{url}` `{id}` `{priority}` `{comment}` `{claimed_by}` `{due_date}` `{sprint}` `{tags}` `{count}`. Leave a template blank to skip that trigger. If a project leaves a template unset, the defaults in this manifest are used.
