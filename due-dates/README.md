# Due dates

Emails `task.due_changed` and `task.overdue` using the site Mailgun or SMTP settings. Use this when the full Email extension is too noisy. Copy this folder to `data/extensions/due-dates` (name must stay `due-dates`), restart, enable in Admin → Extensions, then set a notify address on the project Extensions tab (or Notify me).

Overdue fires once per task per due date. Empty `{due_date}` means the date was cleared.
