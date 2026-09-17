# Discord notifications

Posts task events to a Discord channel using an incoming webhook. Core is not patched; this folder is a drop-in extension. Project owners configure a team channel; any member can add a personal “Notify me” destination. Personal inbox tasks (no project) use Settings → Integrations.

Messages include `{content}` plus a small embed (title, status, actor, Open URL). When Discord returns a message id, later `task.commented` events post as thread replies (`wait=true` on create).

## Install

1. Copy this directory to `data/extensions/discord` on the server (the folder name must stay `discord`).
2. Restart Ordryn.
3. Site admin: Admin → Extensions → Discord → Enable, then Save.
4. In Discord: channel settings → Integrations → Webhooks → New Webhook. Copy the URL.
5. Project owner: open the project → Settings → Extensions → Discord. Paste the webhook URL, choose triggers, edit messages if you want, Enable, Save.
6. Optional: set `PUBLIC_URL` in `.env` (for example `https://todo.example.com`) so `{url}` in templates becomes a clickable task link.

The site admin must enable the extension before any project can post. Channel, triggers, and templates live on the project, not in Admin → Extensions.

Sibling examples cover Slack, Microsoft Teams, Google Chat, a generic HTTPS webhook, ntfy, email, due-dates, comments, claimed, board activity, lifecycle, mentions, callback-bot, and join-requests.

With **Only notify when status changes** on, `task.updated` is skipped unless the kanban/list status actually changed.

Templates may use `{task}` `{name}` `{status}` `{old_status}` `{project}` `{actor}` `{url}` `{id}` `{priority}` `{comment}` `{claimed_by}` `{due_date}` `{sprint}` `{tags}` `{mentions}` `{member}` `{count}`. Leave a template blank to skip that trigger. If a project leaves a template unset, the defaults in this manifest are used.

Optional `description` on the extension and on each setting is shown in the UI under the field label.

Optional actor mention map (JSON on the project form) replaces `{actor}` with a Discord mention id when the Ordryn username matches.
