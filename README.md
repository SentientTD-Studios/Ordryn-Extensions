# Ordryn Extensions

Example and future Ordryn extensions, kept out of the Ordryn core so they can evolve independently.

An extension is a **folder with a `manifest.json`**. Ordryn loads it at startup, renders the UI itself, and (for notifications) POSTs JSON to a webhook. There is no extension process, no plugin VM, and no custom JavaScript in host API 1.

## Layout

```
my-extension/
  manifest.json    # required
  README.md        # recommended (install notes for operators)
```

Rules:

- The **folder name must equal `manifest.id`**.
- Drop the folder into Ordryn’s extensions directory (default `data/extensions/`). Hidden directories (names starting with `.`) are ignored.
- Restart Ordryn after adding or changing files on disk. Manifests are not hot-reloaded.
- One bad extension fails that folder only; others still load. Duplicate ids disable **both** copies.

Override the scan path with `EXTENSIONS_DIR` (legacy alias: `MODS_DIR`).

## Create an extension

1. Choose an `id`: lowercase letters, digits, and hyphens; must start with a letter; 2–33 characters. Example: `severity`, `fields-demo`.
2. Create a directory named exactly that `id`.
3. Add `manifest.json` with the required identity fields, plus either `fields` (custom task fields), `hooks` + `delivery` (outbound notifications), or both.
4. Copy the directory to `data/extensions/<id>` on the Ordryn server.
5. Restart Ordryn. Check server logs for `extensions: loaded <id>` or `extensions: failed …`.
6. **Site admin:** Admin → Extensions → enable the extension → Save.
7. **Project owner:** Project settings → Extensions → enable it for that project → Save. Configure secrets, triggers, and templates there if the extension has them.

Both site and project enablement are required. Personal tasks (no project) never get custom fields or outbound posts.

## Required files and fields

The only required file is `manifest.json`.

| Field | Required | Notes |
| --- | --- | --- |
| `id` | yes | Must match the folder name. Pattern `^[a-z][a-z0-9-]{1,32}$`. |
| `name` | yes | Display name in Admin and project settings. |
| `version` | yes | Opaque string shown in the UI (semver is conventional, not enforced). |
| `host_api` | yes | Integer. Use `1`. Higher than the running Ordryn build will fail to load. |
| `description` | no | Shown under the name. Max 400 characters. |

Minimal identity:

```json
{
  "id": "severity",
  "name": "Severity",
  "version": "1.0.0",
  "host_api": 1,
  "description": "Adds a severity field to project tasks."
}
```

That manifest loads, but it does nothing until you add `fields` and/or `hooks`.

## Custom fields

Declare task fields in `fields`. Ordryn stores values as `{extension id}.{field key}` (for example `severity.level`) and renders them in the task sidebar, list, and/or kanban. No `ui.js` is needed.

```json
{
  "id": "severity",
  "name": "Severity",
  "version": "1.0.0",
  "host_api": 1,
  "fields": [
    {
      "key": "level",
      "type": "enum",
      "label": "Severity",
      "description": "How bad it is if this work slips.",
      "show_on": ["sidebar", "kanban"],
      "options": [
        { "value": "low", "label": "Low" },
        { "value": "medium", "label": "Medium" },
        { "value": "high", "label": "High" },
        { "value": "critical", "label": "Critical" }
      ]
    }
  ]
}
```

### Field object

| Field | Required | Notes |
| --- | --- | --- |
| `key` | yes | Local key. Pattern `^[a-z][a-z0-9_-]{0,32}$`. Combined with the extension id for storage. |
| `type` | yes | One of `string`, `number`, `boolean`, `enum`, `url`, `user`. |
| `label` | yes | Sidebar / badge label. |
| `description` | no | Help text under the control. Max 400 characters. |
| `required` | no | If true, the value cannot be cleared once set. |
| `show_on` | no | Surfaces: `sidebar`, `kanban`, `list`. Defaults to `["sidebar"]`. |
| `options` | enum only | Non-empty list of `{ "value", "label?", "color?" }`. Duplicate `value`s are rejected. Empty `label` falls back to `value`. |

`show_on` controls where the host UI shows the value:

- `sidebar` — editor on the task detail sidebar (always the place to set it).
- `kanban` — badge on board cards.
- `list` — badge on list rows.

### Value rules

| Type | Stored as | Validation |
| --- | --- | --- |
| `string` | trimmed string | Empty clears the field. |
| `number` | JSON number | Must be numeric. |
| `boolean` | JSON boolean | Must be `true` or `false`. |
| `enum` | option `value` | Must match an option. |
| `url` | string | `http` or `https` with a host. |
| `user` | project member user id | Must belong to the project. |

Clear a value with JSON `null` or an empty string (except `required` fields).

See [fields-demo](fields-demo/) for one field of each type, and [severity](severity/) for a single enum used as a kanban badge.

## Notifications (hooks)

Hook extensions subscribe to task events and POST JSON to a webhook URL the project owner configures. Host API 1 has four delivery types:

| `delivery.type` | Use for | Body |
| --- | --- | --- |
| `http.webhook` | Any public HTTPS receiver (n8n, Zapier, Mattermost, Google Chat, custom listeners) | Chosen by `format` (see below) |
| `discord.webhook` | Discord incoming webhooks | `{"content": "…"}` |
| `slack.webhook` | Slack incoming webhooks | `{"text": "…"}` |
| `teams.webhook` | Teams Workflows incoming webhooks | Adaptive Card |

Prefer the dedicated Discord / Slack / Teams examples when targeting those services: they send the payload the service expects and restrict the URL host. Use `http.webhook` for everything else.

```json
{
  "id": "webhook",
  "name": "Generic webhook",
  "version": "1.0.0",
  "host_api": 1,
  "hooks": [
    { "on": "task.created" },
    { "on": "task.updated" },
    { "on": "task.deleted" },
    { "on": "task.commented" }
  ],
  "delivery": {
    "type": "http.webhook",
    "url_from": "webhook_url",
    "format": "json"
  },
  "settings": [
    {
      "key": "webhook_url",
      "type": "secret",
      "label": "Webhook URL",
      "required": true,
      "scope": "project"
    }
  ],
  "templates": {
    "task.created": "New task {name} in {project} ({url})"
  }
}
```

A hook fires only when **all** of these are true:

1. The extension is enabled for the site.
2. The extension is enabled for the project.
3. The event is listed in `hooks`.
4. The project owner included it in **Triggers**.
5. A template exists for that event (project override, else manifest default). An empty template skips the send.
6. For `task.updated`, if the project has **Only notify when status changes** on, the kanban/list status actually changed.

Personal tasks never post.

### Hook names

Declared in `hooks[].on`. No duplicates.

| Event | Dispatched today |
| --- | --- |
| `task.created` | yes |
| `task.updated` | yes |
| `task.deleted` | yes |
| `task.commented` | yes |
| `task.reordered` | accepted in the manifest, not dispatched |
| `project.updated` | accepted in the manifest, not dispatched |
| `join.request` | accepted in the manifest, not dispatched |

Only the four `task.*` events that say “yes” produce outbound messages. The others are reserved names; listing them will not send anything until Ordryn wires them up.

### Delivery

| Field | Required | Notes |
| --- | --- | --- |
| `type` | yes | `http.webhook`, `discord.webhook`, `slack.webhook`, or `teams.webhook`. |
| `url_from` | yes | Must be a `settings` key (almost always a `secret`). |
| `format` | no | Only for `http.webhook`. One of `json`, `text`, or `content`. Defaults to `text` (`{"text": "…"}`) if omitted. Rejected on provider-specific types. |

Every delivery is `POST` with `Content-Type: application/json`, a 5s timeout, and no redirects. Interpolated message text is trimmed to 2000 characters; `@everyone` / `@here` are stripped.

Project owners can send a test post from Project settings → Extensions (bypasses enable/trigger filters; still needs a webhook URL).

#### `http.webhook` formats

| `format` | Body |
| --- | --- |
| `json` | Structured event (see [webhook](webhook/)). Includes `text` and `content` plus `event`, `id`, `name`, `task`, `status`, `old_status`, `project`, `actor`, `url`, `priority`. |
| `text` | `{"text": "…"}` (Slack-style). Default if `format` is omitted. |
| `content` | `{"content": "…"}` (Discord-style). |

Example `format: json` body:

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

#### URL rules

All types require **HTTPS**, no credentials in the URL, and no fragment. Loopback, private, link-local, and similar destinations are rejected.

| Type | Extra host / path rules |
| --- | --- |
| `http.webhook` | Any public HTTPS host. Non-443 ports are allowed. |
| `discord.webhook` | `discord.com`, `discordapp.com`, `canary.discord.com`, or `ptb.discord.com`; path `/api/webhooks/…`; no query string. |
| `slack.webhook` | `hooks.slack.com`; path `/services/…`, `/triggers/…`, or `/workflows/…`. |
| `teams.webhook` | `*.logic.azure.com`, `*.api.powerplatform.com`, `webhook.office.com`, or Outlook Office hosts; path required. Query strings are allowed (Workflows URLs include them). |

### Settings

Settings drive the Admin / project forms. Keys use `^[a-z][a-z0-9_-]{0,32}$`.

| Field | Required | Notes |
| --- | --- | --- |
| `key` | yes | Unique within the extension. |
| `type` | yes | `secret`, `hook_select`, `bool`, or `project_ids`. |
| `label` | yes | Form label. |
| `description` | no | Help text. Max 400 characters. |
| `required` | no | UI hint. An empty secret still skips delivery rather than erroring. |
| `scope` | no | `site` (default) or `project`. |

| Type | Purpose |
| --- | --- |
| `secret` | Encrypted value (webhook URL). The API never returns the secret, only whether it is set. |
| `hook_select` | Project trigger checkboxes, one per declared hook. |
| `bool` | Boolean flag. The notification examples use `status_only` for “only notify when status changes”. |
| `project_ids` | Accepted in the manifest; no host UI in API 1. Do not rely on it. |

Site admin currently persists **enabled** only. Webhook URL, triggers, templates, and `status_only` live on the **project**. Use `"scope": "project"` for those.

An extension appears on the project Extensions tab if it has any project-scoped setting **or** any custom fields.

### Templates

`templates` is a map of hook name → message string. Keys must be known hook names.

Project owners can override them. If a project sets a template to blank, that trigger is skipped. If the project leaves a template unset, the manifest default is used. If neither has text, nothing is sent.

Placeholders (case-insensitive). Unknown tokens become empty.

| Token | Value |
| --- | --- |
| `{task}`, `{name}` | Task title (max 120 characters; newlines flattened) |
| `{status}` | Current status name |
| `{old_status}` | Previous status (updates) |
| `{project}` | Project name |
| `{actor}` | Username or email of the person who caused the event |
| `{url}` | Task link, from `PUBLIC_URL` (or a URL-like `BASE_PATH`) + `/tasks/{id}` |
| `{id}` | Task id |
| `{priority}` | `None`, `Low`, `Medium`, or `High` |

Set `PUBLIC_URL` in Ordryn’s `.env` (for example `https://todo.example.com`) so `{url}` is a clickable link. Markup is passed through unchanged — use Discord `**bold**`, Slack `*bold*`, or plain text depending on the destination.

See [webhook](webhook/) for a generic HTTPS receiver, or [discord](discord/), [slack](slack/), and [teams](teams/) for provider-specific payloads.

## Enablement

| Layer | Who | What |
| --- | --- | --- |
| Site | Admin | Admin → Extensions. Toggle **Enable**, Save. Failed extensions cannot be configured. |
| Project | Project **owner** only | Project settings → Extensions. Enable, secrets, triggers, templates. Blocked until the site toggle is on. |

Editors cannot manage extensions.

## Reserved: `ui`

`ui` is an optional relative path inside the extension folder (no `..`, no leading `/`). If set, the file must exist or load fails.

Host API 1 does **not** serve or execute that file. Custom fields and webhook settings are rendered by Ordryn. Leave `ui` omitted unless you are experimenting against a newer host.

## Checklist

- [ ] Folder name equals `id`
- [ ] `id`, `name`, `version`, `host_api` present (`host_api` is `1`)
- [ ] Field / setting keys are unique and match the key pattern
- [ ] Enums have `options`; non-enums do not
- [ ] Hook extensions declare `hooks`, `delivery.type` (`http.webhook` or a provider type), and a matching settings key for `url_from`
- [ ] `http.webhook` sets `format` if you need `json` or `content` (otherwise you get `{"text": "…"}`)
- [ ] Project-facing settings use `"scope": "project"`
- [ ] Copied to `data/extensions/<id>` and Ordryn restarted
- [ ] Enabled in Admin, then enabled (and configured) on a project

## Examples in this repo

| Folder | Kind | What it shows |
| --- | --- | --- |
| [severity](severity/) | Fields | One enum, sidebar + kanban badge |
| [fields-demo](fields-demo/) | Fields | All six field types |
| [webhook](webhook/) | Hooks | Generic `http.webhook` with `format: json` |
| [discord](discord/) | Hooks | `discord.webhook` (`{"content": "…"}`) |
| [slack](slack/) | Hooks | `slack.webhook` (`{"text": "…"}`) |
| [teams](teams/) | Hooks | `teams.webhook` (Adaptive Card) |

Each example has its own README with install steps. Safe to enable `severity` and `fields-demo` together (different enum keys). Enable as many notification extensions as you want; each project owner still configures a URL per extension.
