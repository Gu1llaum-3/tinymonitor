import json
import os
import psutil

# Calculate default load thresholds based on CPU count
cpu_count = psutil.cpu_count() or 1
load_warning = cpu_count * 0.7
load_critical = cpu_count * 0.9

DEFAULT_CONFIG = {
    "refresh": 2,
    "cooldown": 60,
    "log_file": "tinymonitor.log",
    "load": {
        "warning": load_warning,
        "critical": load_critical,
        "enabled": True,
        "duration": 60
    },
    "cpu": {
        "warning": 70,
        "critical": 90,
        "enabled": True,
        "duration": 0
    },
    "memory": {
        "warning": 70,
        "critical": 90,
        "enabled": True,
        "duration": 0
    },
    "filesystem": {
        "warning": 80,
        "critical": 90,
        "enabled": True,
        "duration": 0
    },
    "reboot": {
        "enabled": True,
        "duration": 0
    },
    "alerts": {
        "google_chat": {
            "enabled": False,
            "webhook_url": ""
        }
    }
}

def load_config(config_path=None):
    config = DEFAULT_CONFIG.copy()
    
    # 1. Explicit path from CLI
    if config_path:
        if os.path.exists(config_path):
            print(f"Loading config from: {config_path}")
            try:
                with open(config_path, 'r') as f:
                    user_config = json.load(f)
                    config.update(user_config)
            except Exception as e:
                print(f"Error loading config file {config_path}: {e}")
        else:
            print(f"Warning: Config file {config_path} not found. Using defaults.")
        return config

    # 2. Search in standard locations (Priority order)
    search_paths = [
        os.path.join(os.getcwd(), "config.json"),                 # Current directory
        os.path.expanduser("~/.config/tinymonitor/config.json"),  # User config
        "/etc/tinymonitor/config.json"                            # System config
    ]

    for path in search_paths:
        if os.path.exists(path):
            print(f"Loading config from: {path}")
            try:
                with open(path, 'r') as f:
                    user_config = json.load(f)
                    config.update(user_config)
                return config
            except Exception as e:
                print(f"Error loading config file {path}: {e}")
                # If a file is found but invalid, we stop searching to avoid confusion
                return config
            
    print("No config file found. Using internal defaults.")
    return config
