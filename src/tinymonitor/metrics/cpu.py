import psutil
from tinymonitor.metrics.provider import MetricProvider

class CpuMetric(MetricProvider):
    def __init__(self, name, config):
        super().__init__(name, config)
        # Premier appel pour initialiser les compteurs de psutil
        # Sinon le premier check() renverrait 0.0
        psutil.cpu_percent(interval=None)

    def check(self):
        # interval=None car appelé en boucle
        cpu_percent = psutil.cpu_percent(interval=None)
        level = None
        
        if cpu_percent >= self.config.get("critical", 90):
            level = "CRITICAL"
        elif cpu_percent >= self.config.get("warning", 70):
            level = "WARNING"
            
        return [("CPU", level, f"{cpu_percent}%")]
