import asyncio
from telethon import TelegramClient
from secrets import API_ID, API_HASH, SESSION_NAME, TARGET
from send_notification import send_alarm



async def main():
    client = TelegramClient(SESSION_NAME, API_ID, API_HASH)
    await client.start()  # si no hay sesión te pedirá el código por consola/telegram

    # obtener todos los diálogos (paginado internamente)
    dialogs = await client.get_dialogs()

    # intentar encontrar por username, title o id
    target_dialog = None
    for d in dialogs:
        # d.entity puede ser User/Chat/Channel
        ent = d.entity
        ent_id = getattr(ent, 'id', None)

        if (ent_id is not None and ent_id == TARGET):
            target_dialog = d
            break

    if not target_dialog:
        print("No se encontró el chat objetivo. Puedes usar username, title exacto o id numérico.")
        print("Chats encontrados con sus ids:")
        for d in dialogs:
            print(f"{d.name}: {d.entity.id}")
        await client.disconnect()
        return

    unread = target_dialog.unread_count or 0
    print(f"Chat: {target_dialog.name}")
    print(f"Unread count: {unread}")

    if unread > 0:
        msgs = await client.get_messages(target_dialog.entity, limit=unread)
        for m in reversed(msgs):  # del más antiguo al más reciente
            sender = await m.get_sender()
            sender_name = getattr(sender, 'first_name', None) or getattr(sender, 'username', None) or 'Unknown'
            print(f"- [{m.date}] {sender_name}: {m.text[:200]!r}")

        send_alarm()
    await client.disconnect()

if __name__ == '__main__':
    asyncio.run(main())
