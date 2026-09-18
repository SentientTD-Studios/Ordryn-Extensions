# Retro

Host API 2 example: a **kanban tab** (`surfaces.at: kanban.tab`) plus the project **document store**. It is not a chat webhook.

Project members open **Retro** next to List / Board. Cards live in `extension_store` under `sprint:{id}` or `sprint:backlog`, keyed to the sprint switcher. The iframe still cannot read the session; the parent Vue app forwards `store.get` / `store.put` over `postMessage`.

## Install

1. Copy this directory to `data/extensions/retro` (folder name must stay `retro`).
2. Restart Ordryn. Logs should show `extensions: loaded retro`.
3. Site admin: Admin → Extensions → Retro → Enable, then Save.
4. Project owner: a **kanban** project → Settings → Extensions → Retro → Enable, Save.
5. Open the project board and click **Retro**.

Editors and owners can add cards, drag them between columns, edit or delete their own cards, and vote. Viewers can read the board.

## Store document

Key: `sprint:12` or `sprint:backlog`.

```json
{
  "columns": ["went-well", "improve", "actions"],
  "cards": [
    {
      "id": "…",
      "column": "went-well",
      "text": "Shipped invites",
      "author_id": 1,
      "votes": [1, 2]
    }
  ]
}
```

Writes use optimistic `revision`. A conflict reloads the server document and retries.
