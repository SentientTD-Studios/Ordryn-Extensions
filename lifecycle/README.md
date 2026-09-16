# Project lifecycle

JSON webhook for sprint dates, membership, archive/restore, and the `task.moved` split (`task.project_changed` / `task.sprint_changed`). `task.moved` still fires; the extra events let a relay tell project moves apart from sprint changes.

Copy this folder to `data/extensions/lifecycle` (name must stay `lifecycle`), restart, enable in Admin → Extensions, then paste a public HTTPS URL on the project Extensions tab.

`{member}` is set on join/leave. `{sprint}` is the sprint name. Rotate the signing secret for `X-Ordryn-Signature`.
