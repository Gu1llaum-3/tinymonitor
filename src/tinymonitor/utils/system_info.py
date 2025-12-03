import socket
import psutil
import time
import requests

def get_hostname():
    return socket.gethostname()

def get_private_ip():
    try:
        # Connect to an external server to determine the interface used for default route
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        s.connect(("8.8.8.8", 80))
        ip = s.getsockname()[0]
        s.close()
        return ip
    except Exception:
        return "127.0.0.1"

def get_public_ip():
    try:
        return requests.get("https://api.ipify.org", timeout=3).text
    except Exception:
        return "N/A"

def get_load_avg():
    try:
        load = psutil.getloadavg()
        return f"{load[0]:.2f}, {load[1]:.2f}, {load[2]:.2f}"
    except AttributeError:
        return "N/A" # Windows doesn't have getloadavg

def get_uptime():
    uptime_seconds = time.time() - psutil.boot_time()
    hours, remainder = divmod(uptime_seconds, 3600)
    minutes, seconds = divmod(remainder, 60)
    return f"{int(hours)}h {int(minutes)}m"
