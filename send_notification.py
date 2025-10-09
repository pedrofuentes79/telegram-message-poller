import subprocess
import os
from datetime import datetime, timedelta
from pathlib import Path

ALARM_LOG_FILE = Path.home() / '.telegram_poller_alarms.log'
ALARM_COOLDOWN_HOURS = 24

def send_notification(title, message, urgency="normal"):
    try:
        subprocess.run([
            'notify-send',
            '-u', urgency,
            '-a', 'Telegram Poller',
            '-i', 'telegram',
            title,
            message
        ], check=True)
        print(f"Notification sent: {title} - {message}")
    except subprocess.CalledProcessError as e:
        print(f"Failed to send notification: {e}")
    except FileNotFoundError:
        print("notify-send not found. Install libnotify package.")


def get_current_sink():
    result = subprocess.run(['pactl', 'get-default-sink'], 
                          capture_output=True, text=True, check=True)
    return result.stdout.strip()

def get_speaker_sink():
    sinks = subprocess.run(['pactl', 'list', 'sinks', 'short'], capture_output=True, text=True, check=True)
    sink = None
    for line in sinks.stdout.strip().split('\n'):
        if "analog-stereo" in line:
            sink = line.split('\t')[1]
            break

    return sink

def set_sink(sink_name):
    subprocess.run(['pactl', 'set-default-sink', sink_name], check=True)


def play_audio(audio_file, sink_name=None):
    if sink_name:
        subprocess.run(['paplay', '--device', sink_name, audio_file], check=True)
    else:
        subprocess.run(['paplay', audio_file], check=True)


def should_send_alarm():
    if not ALARM_LOG_FILE.exists():
        return True
    
    with open(ALARM_LOG_FILE, 'r') as f:
        last_alarm = f.read().strip()
    
    if not last_alarm:
        return True
    
    try:
        last_alarm_time = datetime.fromisoformat(last_alarm)
        time_since_alarm = datetime.now() - last_alarm_time
        return time_since_alarm > timedelta(hours=ALARM_COOLDOWN_HOURS)
    except ValueError:
        return True


def log_alarm():
    with open(ALARM_LOG_FILE, 'w') as f:
        f.write(datetime.now().isoformat())


def send_alarm():
    forced_sink = get_speaker_sink()
    send_audio_notification(
        title="Alert",
        message="You have unread messages",
        audio_file="/usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga",
        forced_sink=forced_sink,
        skip_cooldown=True
    )

def send_audio_notification(title, message, audio_file, forced_sink=None, urgency="critical", skip_cooldown=False):
    if not should_send_alarm() and not skip_cooldown:
        print("Alarm on cooldown. Skipping audio notification.")
        send_notification(title, message, urgency="normal")
        return
    
    original_sink = None
    try:
        # Send desktop notification
        send_notification(title, message, urgency)
        
        # Save current sink and switch to forced sink if specified
        original_sink = get_current_sink()
        
        if forced_sink:
            print(f"Switching from {original_sink} to {forced_sink}")
            set_sink(forced_sink)
            play_audio(audio_file)
        else:
            play_audio(audio_file)
        
        log_alarm()
        print(f"Audio notification played. Logged at {datetime.now()}")
        
    except Exception as e:
        print(f"Failed to play audio notification: {e}")
    
    finally:
        # Restore original sink
        if original_sink and forced_sink:
            print(f"Restoring sink to {original_sink}")
            set_sink(original_sink)
