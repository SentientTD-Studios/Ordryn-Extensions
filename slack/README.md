# Slack notifications

Posts task events to a Slack channel using an incoming webhook (Block Kit with an Open button). Core is not patched; this folder is a drop-in extension. Project owners configure a team channel; any member can add a personal “Notify me” destination.

Incoming webhooks cannot thread reliably without a Slack bot. Ordryn stores a provider message id when the webhook response includes `ts` and will thread later comments only in that case.

## Install

1. Copy this directory to `data/extensions/slack` on the server (the folder name must stay `slack`).
2. Restart Ordryn.
3. Site admin: Admin → Extensions → Slack → Enable, then Save.
4. In Slack: create an incoming webhook for the channel (Slack app incoming webhooks, or channel Integrations). Copy the URL (`https://hooks.slack.com/services/…`).
5. Project owner: open the project → Settings → Extensions → Slack. Paste the webhook URL, choose triggers, edit messages if you want, Enable, Save.
6. Optional: set `PUBLIC_URL` in `.env` (for example `https://todo.example.com`) so `{url}` in templates becomes a clickable task link.

The site admin must enable the extension before any project can post. Channel, triggers, and templates live on the project, not in Admin → Extensions.

With **Only notify when status changes** on, `task.updated` is skipped unless the kanban/list status actually changed.

Templates use Slack mrkdwn (`*bold*`). Tokens: `{task}` `{name}` `{status}` `{old_status}` `{project}` `{actor}` `{url}` `{id}` `{priority}` `{comment}` `{claimed_by}` `{due_date}` `{sprint}` `{tags}` `{mentions}` `{member}` `{count}`. Leave a template blank to skip that trigger. If a project leaves a template unset, the defaults in this manifest are used.
