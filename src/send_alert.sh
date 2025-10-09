#! /bin/bash
should_play_alert=$1


play_alert() {
    paplay --device $1 /usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga
}


# Gets speaker sink                  => filters my speakers => gets the second column (its name)
TARGET_SINK=$(pactl list sinks short | grep "analog-stereo" | awk '{print $2}')
# Same but for active sink
ACTIVE_SINK=$(pactl get-default-sink)

# Send a desktop notification
notify-send -u "critical" -a "Telegram Poller" -i "telegram" "You have unread messages" "Check your telegram"

# Play an alert on the target sink and on the active sink
if [ "$should_play_alert" == "send" ]; then
    echo "playing alert"
    play_alert $TARGET_SINK
    play_alert $ACTIVE_SINK

    # add a new line to the alarm log with date and time
    current_datetime=$(date +%Y-%m-%dT%H:%M:%S)
    echo "$current_datetime" >> logs/alarms.log
fi