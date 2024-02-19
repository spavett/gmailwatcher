Scheduled OpenFAAS Golang function which sets a watch on the authorised Gmail Inbox. This sets a GCP pub/sub topic which will receive a notification whenever a new email is received in the Gmail account. The topic has an HTTP call set to call whenever an event is published on the topic.

Once a watcher is set it will produce notification to the pub/sub topic for the next seven days. GCP reccommend that the watcher is renewed once a day for the duration of the period you want to receive notifications.
