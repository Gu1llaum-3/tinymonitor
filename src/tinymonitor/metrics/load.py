import psutil
from tinymonitor.metrics.provider import MetricProvider

class LoadMetric(MetricProvider):
    def check(self):
        try:
            # Load Average (1 min, 5 min, 15 min)
            load_1, _, _ = psutil.getloadavg()
        except AttributeError:
            # Windows doesn't support getloadavg
            return []

        level = None
        
        if load_1 >= self.config.get("critical", 9.0):
            level = "CRITICAL"
        elif load_1 >= self.config.get("warning", 7.0):
            level = "WARNING"

        return [("LOAD", level, f"{load_1:.2f}")]
