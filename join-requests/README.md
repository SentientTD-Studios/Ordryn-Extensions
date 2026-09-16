# Join requests

Site-wide `join.request`, `join.approved`, and `join.denied` JSON webhook. Copy this folder to `data/extensions/join-requests` (name must stay `join-requests`), restart, then in Admin → Extensions enable it and paste a public HTTPS URL. It does not appear on project Extensions tabs.

The payload includes `join_email`, `join_message`, and `{url}` pointing at Admin → Join requests when `PUBLIC_URL` is set.
