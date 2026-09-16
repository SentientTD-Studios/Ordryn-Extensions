# Mentions

ntfy alerts for `task.mentioned` (comment @-mentions of project members). Copy this folder to `data/extensions/mentions` (name must stay `mentions`), restart, enable in Admin → Extensions, then paste a public HTTPS topic URL.

A team topic receives every mention in the project. Members can open **Notify me**, paste their own topic, and only receive events where they were named. Skip-self ignores mentions you wrote.

`{mentions}` is the comma-separated usernames from the comment. LAN ntfy hosts are rejected (SSRF protection).
