# ntfy

Pushes project task events to an [ntfy](https://ntfy.sh) topic. Copy this folder to `data/extensions/ntfy` (name must stay `ntfy`), restart, enable in Admin → Extensions, then paste a public HTTPS topic URL on the project Extensions tab.

LAN/self-hosted ntfy on a private IP is rejected (SSRF protection). Optional access token is sent as `Authorization: Bearer …`.
