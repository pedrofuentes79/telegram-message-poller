# Telegram Message Poller

A simple Python script to check and display unread messages from a specific Telegram chat.

## Setup

1. **Create a Telegram App**
   - Go to [my.telegram.org](https://my.telegram.org)
   - Create a new app to get your `API_ID` and `API_HASH`
   - Note: If you get an unexpected error, try using an Android phone or emulator

2. **Install Dependencies**
   ```bash
   python -m venv env
   source env/bin/activate 
   pip install telethon
   ```

3. **Configure Secrets**
   - Copy `secrets.example.py` to `secrets.py`
   - Fill in your credentials:
     - `API_ID` - Your Telegram app ID
     - `API_HASH` - Your Telegram app hash
     - `SESSION_NAME` - Name for session file (default: 'session')
     - `TARGET` - The chat/user ID you want to monitor
        - If you don't know the ID, just run the script with a sample one, and it will display all the chats with their IDs.

## Usage

```bash
python main.py
```

On first run, you'll be prompted to authenticate with your Telegram account. The script will then display unread messages from your target chat.

## Notifications

The project includes a notification system with desktop and audio alerts.

### Desktop Notifications

```python
from send_notification import send_notification

send_notification("Title", "Message", urgency="normal")  # normal, low, or critical
```

### Audio Notifications with Speaker Control

Force audio to play on specific speakers (useful when headphones are connected):

1. **List available audio sinks:**
   ```bash
   python list_audio_sinks.py
   ```

2. **Use audio notifications:**
   ```python
   from send_notification import send_audio_notification
   
   send_audio_notification(
       title="Alert!",
       message="Important notification",
       audio_file="/usr/share/sounds/freedesktop/stereo/alarm-clock-elapsed.oga",
       forced_sink="alsa_output.pci-0000_08_00.3.iec958-stereo"  # Your speakers
   )
   ```

**Features:**
- Automatically switches to specified speakers
- Plays audio notification
- Restores previous audio sink
- Logs alarms with 24-hour cooldown (prevents spam)
- Falls back to silent notification if on cooldown
