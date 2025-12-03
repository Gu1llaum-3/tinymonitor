import psutil
from tinymonitor.metrics.provider import MetricProvider

class MemoryMetric(MetricProvider):
    def check(self):
        mem = psutil.virtual_memory()
        mem_percent = mem.percent
        level = None

        if mem_percent >= self.config.get("critical", 90):
            level = "CRITICAL"
        elif mem_percent >= self.config.get("warning", 70):
            level = "WARNING"

        return [("MEMORY", level, f"{mem_percent}%")]
