# Ordryn Extensions

Example and future Ordryn extensions, kept out of the Ordryn core so they can evolve independently.

An extension is a **folder with a `manifest.json`**. Ordryn loads it at startup, renders the UI itself, and (for notifications) POSTs JSON to a webhook. There is no extension process, no plugin VM, and no custom JavaScript in host API 1. A manifest can also declare a **sandboxed HTML panel**, mint **callback tokens** for relays, and opt into extra **inbound webhook actions**.

**Authoring spec (host API 1):** [wiki](https://github.com/sentientTD-Studios/Ordryn-Extensions/wiki) — how to create an extension, the full `manifest.json` schema, hooks, delivery types, settings, controls, [sandboxed UI](https://github.com/sentientTD-Studios/Ordryn-Extensions/wiki/Sandboxed-UI), and [callbacks / inbound actions](https://github.com/sentientTD-Studios/Ordryn-Extensions/wiki/Callbacks-and-inbound).

## Layout

```
my-extension/
  manifest.json    # required
  README.md        # recommended (install notes for operators)
  icon.svg         # optional
  panel.html       # optional (sandboxed UI; set manifest `ui`)
```

- The **folder name must equal `manifest.id`**.
- Copy the folder to Ordryn’s extensions directory (default `data/extensions/`). Hidden directories (names starting with `.`) are ignored.
- Restart Ordryn after adding or changing files on disk. Manifests are not hot-reloaded.
- Override the scan path with `EXTENSIONS_DIR` (legacy alias: `MODS_DIR`).

## Install an example

1. Copy a folder from this repo to `data/extensions/<id>` on the Ordryn server (the folder name must stay the `id`).
2. Restart Ordryn. Check logs for `extensions: loaded <id>` or `extensions: failed …`.
3. **Site admin:** Admin → Extensions → enable → Save.
4. **Project owner:** Project settings → Extensions → enable and configure → Save.

Both site and project enablement are required for project-scoped behavior. Join-requests is site-scoped only (Admin URL; no project Extensions tab). Member destinations (“Notify me”) are configured in Settings → Integrations when the manifest has `scope: member` settings. Each example README has operator install notes.

## Examples in this repo

### Custom fields

| Folder | Kind | What it shows |
| --- | --- | --- |
| [severity](severity/) | Fields | One enum, sidebar + kanban badge |
| [fields-demo](fields-demo/) | Fields | All eight field types (including `date` and `markdown`) |
| [estimate](estimate/) | Fields | Fibonacci points enum, sidebar + kanban |

### Generic HTTPS (`http.webhook`)

| Folder | Kind | What it shows |
| --- | --- | --- |
| [webhook](webhook/) | Hooks | Large event set, `format: json` (not every host hook) |
| [comments](comments/) | Hooks | `task.commented` and `task.mentioned` |
| [activity](activity/) | Hooks | Board/audit events (move, tag, overdue, reorder, membership, sprints) |
| [lifecycle](lifecycle/) | Hooks | Split moves, archive/restore, membership, sprint dates |
| [email](email/) | Hooks | JSON to your HTTPS email relay (not Admin → Email) |
| [due-dates](due-dates/) | Hooks | Due-date, due-soon, and overdue; quieter email relay |
| [join-requests](join-requests/) | Hooks (site) | Site-wide `join.request` / `join.approved` / `join.denied` |

### Chat and push

| Folder | Kind | What it shows |
| --- | --- | --- |
| [discord](discord/) | Hooks | `discord.webhook` |
| [slack](slack/) | Hooks | `slack.webhook` |
| [teams](teams/) | Hooks | `teams.webhook` (Teams Workflows Adaptive Card) |
| [google-chat](google-chat/) | Hooks | `googlechat.webhook` (Chat cards) |
| [ntfy](ntfy/) | Hooks | `ntfy.webhook` |
| [claimed](claimed/) | Hooks | ntfy for claim / unclaim / overdue; `claimed_is_me` |
| [mentions](mentions/) | Hooks | ntfy for `task.mentioned`; member “Notify me” |

### Advanced

| Folder | Kind | What it shows |
| --- | --- | --- |
| [callback-bot](callback-bot/) | Hooks + fields + UI | Sandboxed panel, callback tokens, inbound `complete` / `set_field` |

Safe to enable `severity`, `fields-demo`, and `estimate` together (different enum keys). Enable as many notification extensions as you want; each project owner still configures a URL per extension.
