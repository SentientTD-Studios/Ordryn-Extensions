# Google Chat

Posts task events to a Google Chat space using an incoming webhook. Messages use Chat’s `text` field plus a `cardsV2` card (header, status/actor fields, Open button). The generic webhook example cannot do this: Chat rejects extra JSON fields.

Comments on the same task reuse a `threadKey` so they stay in one thread when Google Chat allows it.

## Install

1. Copy this directory to `data/extensions/google-chat` (the folder name must stay `google-chat`).
2. Restart Ordryn.
3. Site admin: Admin → Extensions → Google Chat → Enable, then Save.
4. In Google Chat: open the space → **Apps & integrations** → **Webhooks** → add a webhook. Copy the URL (`https://chat.googleapis.com/v1/spaces/…/messages?key=…&token=…`).
5. Project owner: project → Settings → Extensions → Google Chat. Paste the URL, choose triggers, Enable, Save.
6. Optional: set `PUBLIC_URL` in `.env` so `{url}` becomes a clickable task link.

The site admin must enable the extension before any project can post. Space URL, triggers, and templates live on the project.

Templates use Google Chat text formatting (`*bold*`). Tokens: `{task}` `{name}` `{status}` `{old_status}` `{project}` `{actor}` `{url}` `{id}` `{priority}` `{comment}` `{claimed_by}` `{due_date}` `{sprint}` `{tags}` `{count}`. Leave a template blank to skip that trigger.

Optional actor mention map (JSON) replaces `{actor}` with a Chat user mention such as `<users/123456789>`.
