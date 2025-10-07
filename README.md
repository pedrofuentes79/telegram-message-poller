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
