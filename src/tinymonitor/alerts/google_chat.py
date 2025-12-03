import requests
import json
import re
from datetime import datetime
from tinymonitor.utils import system_info
from tinymonitor.alerts.provider import AlertProvider, logger

class GoogleChatAlert(AlertProvider):
    def send(self, component, level, value, title, message):
        webhook_url = self.config.get("webhook_url")
        if not webhook_url:
            logger.error(f"[{self.name}] No webhook_url provided")
            return

        # Visual decoration based on status
        if level == "CRITICAL":
            icon = "🚨"
            font_color = "#FF0000"
            title_text = f"CRITICAL ALERT : {component}"
        elif level == "WARNING":
            icon = "⚠️"
            font_color = "#FFA500"
            title_text = f"WARNING : {component}"
        else:
            icon = "ℹ️"
            font_color = "#000000"
            title_text = f"INFO : {component}"

        # System Info
        hostname = system_info.get_hostname()
        execution_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        ip_private = system_info.get_private_ip()
        ip_public = system_info.get_public_ip()
        load_avg = system_info.get_load_avg()
        uptime_pretty = system_info.get_uptime()

        # Sanitize for cardId (only alphanumeric, hyphens, underscores allowed)
        safe_hostname = re.sub(r'[^a-zA-Z0-9_-]', '_', hostname)
        safe_component = re.sub(r'[^a-zA-Z0-9_-]', '_', component)

        # Construct Payload
        payload = {
            "cardsV2": [
                {
                    "cardId": f"tinymonitor-{safe_hostname}-{safe_component}",
                    "card": {
                        "header": {
                            "title": f"{icon} {title_text}",
                            "subtitle": f"Server : {hostname}",
                            "imageUrl": "https://upload.wikimedia.org/wikipedia/commons/thumb/3/35/Tux.svg/1200px-Tux.svg.png",
                            "imageType": "CIRCLE"
                        },
                        "sections": [
                            {
                                "header": "Incident Details",
                                "widgets": [
                                    {
                                        "decoratedText": {
                                            "topLabel": "Monitored Component",
                                            "text": f"<b>{component}</b>",
                                            "startIcon": {"knownIcon": "MEMBERSHIP"}
                                        }
                                    },
                                    {
                                        "decoratedText": {
                                            "topLabel": "Current Value",
                                            "text": f"<font color=\"{font_color}\"><b>{value}</b></font>",
                                            "startIcon": {"knownIcon": "DESCRIPTION"}
                                        }
                                    },
                                    {
                                        "decoratedText": {
                                            "topLabel": "Alert Level",
                                            "text": f"<b>{level}</b>",
                                            "startIcon": {"knownIcon": "STAR"}
                                        }
                                    }
                                ]
                            },
                            {
                                "header": "Machine Context",
                                "collapsible": True,
                                "uncollapsibleWidgetsCount": 2,
                                "widgets": [
                                    {
                                        "textParagraph": {
                                            "text": f"<b>Private IP:</b> {ip_private}<br><b>Public IP:</b> {ip_public}"
                                        }
                                    },
                                    {
                                        "textParagraph": {
                                            "text": f"<b>Load:</b> {load_avg}<br><b>Uptime:</b> {uptime_pretty}"
                                        }
                                    },
                                    {
                                        "textParagraph": {
                                            "text": f"<font color=\"#808080\">Alert Time: {execution_time}</font>"
                                        }
                                    }
                                ]
                            }
                        ]
                    }
                }
            ]
        }

        try:
            response = requests.post(
                webhook_url,
                headers={'Content-Type': 'application/json'},
                data=json.dumps(payload),
                timeout=10
            )
            if response.status_code == 200:
                logger.info(f"[{self.name}] Alert sent successfully")
            else:
                logger.error(f"[{self.name}] Failed to send alert: {response.text}")
        except Exception as e:
            logger.error(f"[{self.name}] Error sending alert: {e}")
