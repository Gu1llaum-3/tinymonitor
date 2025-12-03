import requests
import json
import socket
from datetime import datetime
from tinymonitor.alerts.provider import AlertProvider, logger
from tinymonitor.utils.system_info import get_private_ip, get_public_ip, get_uptime, get_load_avg

class WebhookAlert(AlertProvider):
    def send(self, component, level, value, title, message):
        url = self.config.get("url")
        if not url:
            logger.error(f"[{self.name}] No url provided")
            return

        # Custom headers (e.g. Authorization)
        headers = self.config.get("headers", {})
        if "Content-Type" not in headers:
            headers["Content-Type"] = "application/json"

        timeout = self.config.get("timeout", 10)

        # Standardized JSON Payload
        payload = {
            "timestamp": datetime.now().isoformat(),
            "alert": {
                "level": level,
                "component": component,
                "value": str(value),
                "title": title,
                "message": message
            },
            "host": {
                "hostname": socket.gethostname(),
                "ip_private": get_private_ip(),
                "ip_public": get_public_ip(),
                "uptime": get_uptime(),
                "load_average": get_load_avg()
            }
        }

        try:
            response = requests.post(
                url, 
                json=payload, 
                headers=headers, 
                timeout=timeout
            )
            
            if response.ok: # Status 200-299
                logger.info(f"[{self.name}] Alert sent successfully")
            else:
                logger.error(f"[{self.name}] Failed to send alert: {response.status_code} {response.text}")
        except Exception as e:
            logger.error(f"[{self.name}] Error sending alert: {e}")
