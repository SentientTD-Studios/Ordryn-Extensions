# Comments webhook

Posts `task.commented` and `task.mentioned` as structured JSON (`event`, `comment`, `mentions`, `actor`, `url`, and the usual task fields). Use this for a discussion pipeline without status-change noise.

Copy this folder to `data/extensions/comments` (name must stay `comments`), restart, enable in Admin → Extensions, then paste an HTTPS URL on the project Extensions tab.

Outbound HMAC: rotate the signing secret on the Extensions panel; requests include `X-Ordryn-Signature: sha256=…` over the raw body.

`{comment}` is truncated. Slack/Discord incoming webhooks should use those dedicated extensions instead.
