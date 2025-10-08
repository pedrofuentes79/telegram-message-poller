#!/usr/bin/env python3
from send_notification import send_notification, send_audio_notification

# Simple desktop notification
# send_notification("Test", "This is a simple notification")

# Audio notification with forced speaker selection
# First, run list_audio_sinks.py to see your available sinks
send_audio_notification(
    title="Important Alert!",
    message="You have unread messages",
    audio_file="/usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga",
    forced_sink="alsa_output.pci-0000_08_00.3.analog-stereo"  
)

# Without forced sink (uses current default)
# send_audio_notification(
#     title="Alert",
#     message="New notification",
#     audio_file="/usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga",
#     skip_cooldown=True
# )

