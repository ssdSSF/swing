## Problem Statement

[Swing Education](https://swingeducation.com) is not sending the notifications through their App right away when a new opening is posted. It waits about 10 - 15 minutes if no one is signing up for the opening then it pushes the notification. It leads to a situation where a user must constantly refresh the App before an openings is taken by other user because of high demand (not enough openings).

## Solution

This project is to poll openings and send new notifications to the Slack channel that you designated so that your can recieve openings immediately through Slack notifications. The fastest it can poll is every 10 seconds. To avoid bot detection, the polling will slow down between 00:00 to 04:59 in the midnight to every 10 minutes.

Please make sure you have a [~/.swing-secrets.yaml](./swing-secrets.yaml) in your home dir:
```
% cat ~/.swing-secrets.yaml 
googleToken: AMf-vBwpSdT-6EjzxaVQms9E.... # get it from the Google SSO
slackToken: xoxb-261105... # get it from Slack integration
slackChannel: C08T5TXXXXX # Slack channel ID
slackHeartbeatChannel: C08T85YYY # Slack channel for sending the heartbeat
interval: 30 # in seconds
citiesToSkip: # skip the cities that you don't want to go
  - San Francisco
  - Daly City
```

To build (with `env GOOS=linux GOARCH=arm`)
```
make build
```

To deploy (to a Raspberry Pi, with `export PI_AT_HOST=pi@<your Raspberry Pi hostname>.local`)
```
make deploy
```

Run the worker:
```
./swing-cli list-worker
```


Sample Notification
![Sample Notification](./example-slack-notification.png)