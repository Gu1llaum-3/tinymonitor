import time
import argparse
import logging
from tinymonitor import __version__
from tinymonitor.config import load_config
from tinymonitor.alerts.manager import AlertManager
from tinymonitor.metrics import CpuMetric, MemoryMetric, FilesystemMetric, LoadMetric, RebootMetric, IoMetric

logger = logging.getLogger("tinymonitor")

class Monitor:
    def __init__(self, config):
        self.config = config
        self.alert_manager = AlertManager(config)
        self.last_alert = {} # Flat dict: component_name -> timestamp
        self.alert_states = {}
        self.metrics = []
        self.load_metrics()

    def load_metrics(self):
        # Load standard metrics
        # In the future, this could be dynamic
        if self.config.get("cpu", {}).get("enabled", True):
            self.metrics.append(CpuMetric("cpu", self.config.get("cpu", {})))
        
        if self.config.get("memory", {}).get("enabled", True):
            self.metrics.append(MemoryMetric("memory", self.config.get("memory", {})))
            
        if self.config.get("filesystem", {}).get("enabled", True):
            self.metrics.append(FilesystemMetric("filesystem", self.config.get("filesystem", {})))
            
        if self.config.get("load", {}).get("enabled", True):
            self.metrics.append(LoadMetric("load", self.config.get("load", {})))

        if self.config.get("reboot", {}).get("enabled", True):
            self.metrics.append(RebootMetric("reboot", self.config.get("reboot", {})))

        if self.config.get("io", {}).get("enabled", True):
            self.metrics.append(IoMetric("io", self.config.get("io", {})))

    def process_state(self, component, level, value, duration):
        """
        Manages alert state persistence.
        Returns True if the alert should be triggered, False otherwise.
        """
        if level is None:
            # Return to normal: clear state
            if component in self.alert_states:
                # Only log recovery if we actually triggered an alert state (duration passed)
                if self.alert_states[component].get("alert_triggered", False):
                    logger.info(f"RECOVERY: {component} is back to normal")
                del self.alert_states[component]
            return False

        now = time.time()
        current_state = self.alert_states.get(component)
        
        if not current_state or current_state["level"] != level:
            self.alert_states[component] = {
                "level": level,
                "start_time": now,
                "alert_triggered": False
            }
            if duration <= 0:
                self.alert_states[component]["alert_triggered"] = True
                return True
            else:
                logger.debug(f"Detected {component} {level}, waiting for duration {duration}s")
                return False
        
        elapsed = now - current_state["start_time"]
        if elapsed >= duration:
            self.alert_states[component]["alert_triggered"] = True
            return True
            
        return False

    def trigger_alert(self, component, level, value):
        current_time = time.time()
        last_time = self.last_alert.get(component, 0)

        cooldown = self.config.get("cooldown", 60)
        should_alert = False

        if cooldown < 0:
            # "Alert Once" Mode
            state = self.alert_states.get(component)
            if state:
                start_time = state["start_time"]
                if last_time < start_time:
                    should_alert = True
        else:
            # Classic Mode
            if (current_time - last_time) > cooldown:
                should_alert = True

        if should_alert:
            logger.info(f"ALERT: {component} is {level} ({value})")
            self.alert_manager.send_alert(component, level, value)
            self.last_alert[component] = current_time
        else:
            if cooldown >= 0:
                logger.debug(f"Alert suppressed (cooldown): {component}")
            else:
                logger.debug(f"Alert suppressed (already sent): {component}")

    def run(self):
        logger.info("Starting TinyMonitor...")
        try:
            while True:
                for metric in self.metrics:
                    try:
                        results = metric.check()
                        for component, level, value in results:
                            duration = metric.config.get("duration", 0)
                            if self.process_state(component, level, value, duration):
                                self.trigger_alert(component, level, value)
                    except Exception as e:
                        logger.error(f"Error checking metric {metric.name}: {e}")
                
                time.sleep(self.config.get("refresh", 5))
        except KeyboardInterrupt:
            try:
                logger.info("Stopping TinyMonitor...")
            except KeyboardInterrupt:
                pass

def main():
    parser = argparse.ArgumentParser(description="TinyMonitor System Monitor")
    parser.add_argument("-c", "--config", help="Path to configuration file")
    parser.add_argument("-v", "--version", action="version", version=f"TinyMonitor {__version__}")
    args = parser.parse_args()

    config = load_config(args.config)

    # Configure logging
    handlers = [logging.StreamHandler()]
    log_file = config.get("log_file")
    if log_file:
        handlers.append(logging.FileHandler(log_file))

    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        handlers=handlers,
        force=True
    )

    monitor = Monitor(config)
    try:
        monitor.run()
    except KeyboardInterrupt:
        pass

if __name__ == "__main__":
    main()
