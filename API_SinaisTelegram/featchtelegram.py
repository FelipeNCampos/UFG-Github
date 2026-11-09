import requests
import time

def get_updates(token, offset=None):
    """
    Fetches new messages sent to the bot.
    """
    url = f"https://api.telegram.org/bot{token}/getUpdates"
    params = {"offset": offset, "timeout": 10}
    response = requests.get(url, params=params)
    return response.json()

def send_message(token, chat_id, text):
    """
    Sends a message to a specific chat.
    """
    url = f"https://api.telegram.org/bot{token}/sendMessage"
    payload = {"chat_id": chat_id, "text": text}
    response = requests.post(url, json=payload)
    return response.json()

def main():
    """
    Main bot loop: listens for updates and responds to commands.
    """
    token = "7862862108:AAEARORzD7rLcyl0MflNTC_R8OTnI0g2Oy0"  # Replace with your bot's API token
    offset = None  # Used to track processed updates

    print("Bot is running...")

    while True:
        # Fetch updates
        updates = get_updates(token, offset)

        if updates["ok"]:
            for update in updates["result"]:
                offset = update["update_id"] + 1  # Mark this update as processed

                if "message" in update:
                    chat_id = update["message"]["chat"]["id"]
                    text = update["message"].get("text", "")

                    # Respond to commands
                    if text == "/newchat":
                        send_message(token, chat_id, f"Seu id de chat é {chat_id}")
                    else:
                        send_message(token, chat_id, "Sorry, I don't understand that command.")

if __name__ == "__main__":
    main()