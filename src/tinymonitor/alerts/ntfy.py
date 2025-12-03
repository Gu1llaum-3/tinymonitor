import requests
from datetime import datetime
from tinymonitor.alerts.provider import AlertProvider, logger
from tinymonitor.utils import system_info

class NtfyAlert(AlertProvider):
    def send(self, component, level, value, title, message):
        topic_url = self.config.get("topic_url")
        if not topic_url:
            logger.error(f"[{self.name}] No topic_url provided")
            return

        # Ntfy priority mapping (1=min, 3=default, 5=max)
        if level == "CRITICAL":
            priority = 5
            tags = "rotating_light,critical"
        elif level == "WARNING":
            priority = 3
            tags = "warning"
        else:
            priority = 1
            tags = "information_source"

        # System Info
        hostname = system_info.get_hostname()
        execution_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        ip_private = system_info.get_private_ip()
        ip_public = system_info.get_public_ip()
        load_avg = system_info.get_load_avg()
        uptime_pretty = system_info.get_uptime()
        
        # Enriched message construction (Markdown supported by ntfy)
        full_message = (
            f"**Component** : {component}\n"
            f"**Value**     : {value}\n"
            f"**Level**     : {level}\n\n"
            f"__Machine Context__\n"
            f"🖥️ **Server**    : `{hostname}`\n"
            f"🏠 **Private IP**: `{ip_private}`\n"
            f"🌍 **Public IP** : `{ip_public}`\n"
            f"⚙️ **Load Avg**  : `{load_avg}`\n"
            f"⏱️ **Uptime**    : `{uptime_pretty}`\n"
            f"🕒 **Time**      : {execution_time}"
        )

        headers = {
            "Title": title,
            "Priority": str(priority),
            "Tags": tags,
            "Markdown": "yes"
        }

        # Authentication handling if token is provided
        token = self.config.get("token")
        if token:
            headers["Authorization"] = f"Bearer {token}"

        try:
            response = requests.post(
                topic_url,
                data=full_message.encode('utf-8'),
                headers=headers,
                timeout=10
            )
            if response.status_code == 200:
                logger.info(f"[{self.name}] Alert sent successfully")
            else:
                logger.error(f"[{self.name}] Failed to send alert: {response.text}")
        except Exception as e:
            logger.error(f"[{self.name}] Error sending alert: {e}")
