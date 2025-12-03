import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from datetime import datetime
from tinymonitor.alerts.provider import AlertProvider, logger
from tinymonitor.utils import system_info

class SMTPAlert(AlertProvider):
    def send(self, component, level, value, title, message):
        smtp_host = self.config.get("host")
        smtp_port = self.config.get("port", 587)
        smtp_user = self.config.get("user")
        smtp_pass = self.config.get("password")
        from_addr = self.config.get("from_addr")
        to_addrs = self.config.get("to_addrs") # Can be a list or a comma-separated string
        use_tls = self.config.get("use_tls", True)

        if not all([smtp_host, smtp_user, smtp_pass, from_addr, to_addrs]):
            logger.error(f"[{self.name}] Missing SMTP configuration (host, user, password, from_addr, or to_addrs)")
            return

        if isinstance(to_addrs, str):
            to_addrs = [addr.strip() for addr in to_addrs.split(",")]

        # System Info
        hostname = system_info.get_hostname()
        execution_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        ip_private = system_info.get_private_ip()
        ip_public = system_info.get_public_ip()
        load_avg = system_info.get_load_avg()
        uptime_pretty = system_info.get_uptime()

        # Email Subject
        subject = f"[{level}] {component} on {hostname} - {value}"

        # Email Body (HTML)
        html_content = f"""
        <html>
        <body>
            <h2>{title}</h2>
            <p><strong>Component:</strong> {component}</p>
            <p><strong>Value:</strong> {value}</p>
            <p><strong>Level:</strong> {level}</p>
            <hr>
            <h3>Machine Context</h3>
            <ul>
                <li><strong>Server:</strong> {hostname}</li>
                <li><strong>Private IP:</strong> {ip_private}</li>
                <li><strong>Public IP:</strong> {ip_public}</li>
                <li><strong>Load Avg:</strong> {load_avg}</li>
                <li><strong>Uptime:</strong> {uptime_pretty}</li>
                <li><strong>Time:</strong> {execution_time}</li>
            </ul>
        </body>
        </html>
        """

        msg = MIMEMultipart()
        msg['From'] = from_addr
        msg['To'] = ", ".join(to_addrs)
        msg['Subject'] = subject

        msg.attach(MIMEText(html_content, 'html'))

        try:
            server = smtplib.SMTP(smtp_host, smtp_port)
            server.ehlo()
            if use_tls:
                server.starttls()
                server.ehlo()
            
            server.login(smtp_user, smtp_pass)
            server.sendmail(from_addr, to_addrs, msg.as_string())
            server.close()
            logger.info(f"[{self.name}] Email sent to {to_addrs}")
        except Exception as e:
            logger.error(f"[{self.name}] Failed to send email: {e}")
