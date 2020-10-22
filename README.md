# send_oem_alerts
Telegram bot accepts metrics from Oracle OEM (Cloud Control) and send it to telegram via env variables.

# Requirements:
1. OEM must be 13c+

# Quick start guide: 
2. Register your bot in @BotFather, get bot token
3. Build binary from source (ex for Linux) or  take pre-build binary (send_oem_alert_v.0.2) \
	env GOOS=linux GOARCH=amd64 go build -o send_oem_alert
4. Put binary to **some** directory on OEM host
5. Put your bot token and chatId to config.yml
6. Config notification rule in OEM (Settings -> Notifications -> Scripts and SNMPv1 Traps -> add OS Commands -> add **some** path to binary to "OS Command" field)
7. Add new notification method to your alert rules and/or rulesets. (Settings -> Incidents -> Incident Rules -> edit Rules and add newly created notification method to this action to "Action Summary"
8. You are awesome =)
