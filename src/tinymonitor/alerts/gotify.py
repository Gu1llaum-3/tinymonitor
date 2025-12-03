import requests
from datetime import datetime
from tinymonitor.alerts.provider import AlertProvider, logger
from tinymonitor.utils import system_info

class GotifyAlert(AlertProvider):
    def send(self, component, level, value, title, message):
        server_url = self.config.get("url")
        token = self.config.get("token")
        
        if not server_url or not token:
            logger.error(f"[{self.name}] No url or token provided")
            return

        # Ensure URL ends with /message
        if not server_url.endswith("/message"):
            if not server_url.endswith("/"):
                server_url += "/"
            server_url += "message"

        # Gotify priority mapping (0-10)
        if level == "CRITICAL":
            priority = 8
        elif level == "WARNING":
            priority = 5
        else:
            priority = 2

        # System Info
        hostname = system_info.get_hostname()
        execution_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        ip_private = system_info.get_private_ip()
        ip_public = system_info.get_public_ip()
        load_avg = system_info.get_load_avg()
        uptime_pretty = system_info.get_uptime()
        
        # Enriched message construction (Markdown supported by Gotify)
        # Note: We use two spaces at the end of lines to force a line break in Markdown
        full_message = (
            f"**Component** : {component}  \n"
            f"**Value**     : {value}  \n"
            f"**Level**     : {level}  \n\n"
            f"__Machine Context__  \n"
            f"🖥️ **Server**    : `{hostname}`  \n"
            f"🏠 **Private IP**: `{ip_private}`  \n"
            f"🌍 **Public IP** : `{ip_public}`  \n"
            f"⚙️ **Load Avg**  : `{load_avg}`  \n"
            f"⏱️ **Uptime**    : `{uptime_pretty}`  \n"
            f"🕒 **Time**      : {execution_time}"
        )

        payload = {
            "title": title,
            "message": full_message,
            "priority": priority,
            "extras": {
                "client::display": {
                    "contentType": "text/markdown"
                }
            }
        }

        headers = {
            "X-Gotify-Key": token
        }

        try:
            response = requests.post(server_url, json=payload, headers=headers, timeout=10)
            if response.ok:
                logger.info(f"[{self.name}] Alert sent successfully")
            else:
                logger.error(f"[{self.name}] Failed to send alert: {response.status_code} {response.text}")
        except Exception as e:
            logger.error(f"[{self.name}] Error sending alert: {e}")
