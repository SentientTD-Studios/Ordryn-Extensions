# Board activity

Posts `task.moved`, `task.tagged`, `task.overdue`, `task.reordered`, and `project.updated` as JSON. Use this for an audit or automation feed that should ignore ordinary status edits.

Copy this folder to `data/extensions/activity` (name must stay `activity`), restart, enable in Admin → Extensions, then paste an HTTPS URL on the project Extensions tab.

`task.reordered` is one project-level event with `{count}`. `project.updated` covers name, members, workflow, sprints, and GitHub link changes. Rotate the signing secret for `X-Ordryn-Signature`.
