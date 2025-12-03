import psutil
from tinymonitor.metrics.provider import MetricProvider

class FilesystemMetric(MetricProvider):
    def check(self):
        results = []
        excludes = self.config.get("exclude", [])
        
        for part in psutil.disk_partitions(all=False): # all=False filters out some virtual fs
            try:
                # Filter out Snap loops and read-only squashfs
                if 'loop' in part.device or 'squashfs' in part.fstype:
                    continue
                
                # Filter out Docker overlay/containers if not explicitly requested
                if 'docker' in part.mountpoint or 'overlay' in part.fstype:
                    continue

                # Skip cdrom or other weird devices if needed
                if 'cdrom' in part.opts or part.fstype == '':
                    continue
                
                # User defined excludes (partial match on mountpoint)
                if excludes and any(ex in part.mountpoint for ex in excludes):
                    continue
                    
                usage = psutil.disk_usage(part.mountpoint)
                usage_percent = usage.percent
                level = None

                if usage_percent >= self.config.get("critical", 90):
                    level = "CRITICAL"
                elif usage_percent >= self.config.get("warning", 80):
                    level = "WARNING"

                component_name = f"DISK:{part.mountpoint}"
                results.append((component_name, level, f"{usage_percent}%"))
            except (PermissionError, OSError):
                continue
        return results
