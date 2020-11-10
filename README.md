# send_oem_alerts
Telegram bot accepts metrics from Oracle OEM (Cloud Control) and send it to telegram via env variables.
Alert messages from OEM also can be redirected to PagerDuty API for Voice Call or other tasks

# Requirements:
 OEM must be 13c+

# Quick start guide: 
1. Register your bot in @BotFather, get bot token
2. Build binary from source (ex for Linux) or  take pre-build binary (ex send_oem_alert.linux64.v.0.3) \
	env GOOS=linux GOARCH=amd64 go build -o send_oem_alert
3. Put binary to **some** directory on OEM host
4. Put your bot token and chatId to config.yml
5. Config notification rule in OEM (Settings -> Notifications -> Scripts and SNMPv1 Traps -> add OS Commands -> add **some** path to binary to "OS Command" field)
6. Add new notification method to your alert rules and/or rulesets. (Settings -> Incidents -> Incident Rules -> edit Rules and add newly created notification method to this action to "Action Summary"
7. You are awesome =)
