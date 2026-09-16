# Board activity

Posts board-level JSON: `task.moved` plus `task.project_changed` / `task.sprint_changed`, archive/restore, tags, overdue, reorder, sprint lifecycle, membership, and `project.updated`. Use this for an audit or automation feed that should ignore ordinary status edits.

Copy this folder to `data/extensions/activity` (name must stay `activity`), restart, enable in Admin → Extensions, then paste an HTTPS URL on the project Extensions tab.

`task.reordered` is one project-level event with `{count}`. `{member}` is set on join/leave. Rotate the signing secret for `X-Ordryn-Signature`.
