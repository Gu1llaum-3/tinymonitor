import logging
from concurrent.futures import ThreadPoolExecutor
from tinymonitor.alerts.google_chat import GoogleChatAlert
from tinymonitor.alerts.ntfy import NtfyAlert
from tinymonitor.alerts.gotify import GotifyAlert
from tinymonitor.alerts.smtp import SMTPAlert
from tinymonitor.alerts.webhook import WebhookAlert

logger = logging.getLogger("tinymonitor")

class AlertManager:
    def __init__(self, config):
        self.providers = []
        self.executor = ThreadPoolExecutor(max_workers=5)
        self.load_providers(config.get("alerts", {}))

    def load_providers(self, alerts_config):
        # Here we manually instantiate known plugins.
        # Later, we could do dynamic loading if needed.
        
        # 1. Google Chat
        gc_config = alerts_config.get("google_chat")
        if gc_config and gc_config.get("enabled", False):
            self.providers.append(GoogleChatAlert("google_chat", gc_config))
            logger.info("Alert Provider loaded: Google Chat")

        # 2. Ntfy
        ntfy_config = alerts_config.get("ntfy")
        if ntfy_config and ntfy_config.get("enabled", False):
            self.providers.append(NtfyAlert("ntfy", ntfy_config))
            logger.info("Alert Provider loaded: Ntfy")

        # 3. SMTP
        smtp_config = alerts_config.get("smtp")
        if smtp_config and smtp_config.get("enabled", False):
            self.providers.append(SMTPAlert("smtp", smtp_config))
            logger.info("Alert Provider loaded: SMTP")

        # 4. Webhook (Generic)
        webhook_config = alerts_config.get("webhook")
        if webhook_config and webhook_config.get("enabled", False):
            self.providers.append(WebhookAlert("webhook", webhook_config))
            logger.info("Alert Provider loaded: Webhook")

        # 5. Gotify
        gotify_config = alerts_config.get("gotify")
        if gotify_config and gotify_config.get("enabled", False):
            self.providers.append(GotifyAlert("gotify", gotify_config))
            logger.info("Alert Provider loaded: Gotify")

    def send_alert(self, component, level, value):
        """Distributes the alert to all configured and eligible providers asynchronously."""
        
        title = f"ALERT {level} : {component}"
        message = f"Component {component} is in state {level}. Value: {value}"

        for provider in self.providers:
            if provider.should_send(component, level):
                logger.info(f"Triggering alert via {provider.name} for {component} ({level})")
                self.executor.submit(provider.send, component, level, value, title, message)
