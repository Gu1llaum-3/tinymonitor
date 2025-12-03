import time
import psutil
from tinymonitor.metrics.provider import MetricProvider

class IoMetric(MetricProvider):
    def __init__(self, name, config):
        super().__init__(name, config)
        self.last_counters = psutil.disk_io_counters()
        self.last_time = time.time()

    def parse_threshold(self, value, max_value=None):
        if value is None:
            return float('inf')
        
        if isinstance(value, (int, float)):
            return float(value)
            
        if isinstance(value, str):
            value = value.strip().upper()
            if value.endswith('%'):
                if max_value is None:
                    return float('inf') # Cannot calculate percentage without max
                try:
                    percent = float(value[:-1])
                    return max_value * (percent / 100.0)
                except ValueError:
                    return float('inf')
            
            # Handle units
            units = {'K': 1024, 'M': 1024**2, 'G': 1024**3, 'T': 1024**4}
            for unit, multiplier in units.items():
                # Check for MB, GB, etc.
                if value.endswith(f"{unit}B"):
                     try:
                        num_part = value[:-2]
                        return float(num_part) * multiplier
                     except ValueError:
                        pass
                # Check for M, G, etc.
                elif value.endswith(unit):
                    try:
                        num_part = value[:-1]
                        return float(num_part) * multiplier
                    except ValueError:
                        pass
            
            # Try parsing as plain number
            try:
                return float(value)
            except ValueError:
                pass
                
        return float('inf')

    def check(self):
        current_time = time.time()
        current_counters = psutil.disk_io_counters()
        
        if not self.last_counters or not current_counters:
            self.last_counters = current_counters
            self.last_time = current_time
            return []

        time_delta = current_time - self.last_time
        if time_delta <= 0:
            return []

        read_bytes_delta = current_counters.read_bytes - self.last_counters.read_bytes
        write_bytes_delta = current_counters.write_bytes - self.last_counters.write_bytes
        
        # Avoid negative values if counters reset (unlikely but possible)
        if read_bytes_delta < 0: read_bytes_delta = 0
        if write_bytes_delta < 0: write_bytes_delta = 0
        
        read_speed = read_bytes_delta / time_delta
        write_speed = write_bytes_delta / time_delta
        
        self.last_counters = current_counters
        self.last_time = current_time
        
        # Convert to human readable
        def format_bytes(size):
            power = 2**10
            n = 0
            power_labels = {0 : '', 1: 'K', 2: 'M', 3: 'G', 4: 'T'}
            while size > power:
                size /= power
                n += 1
            return f"{size:.1f}{power_labels.get(n, '')}B/s"

        formatted_read = format_bytes(read_speed)
        formatted_write = format_bytes(write_speed)
        
        level = None
        
        # Parse max_speed if present
        max_speed_val = self.config.get("max_speed")
        max_speed = None
        if max_speed_val:
             max_speed = self.parse_threshold(max_speed_val)
             if max_speed == float('inf'):
                 max_speed = None

        # Thresholds
        warning_val = self.config.get("warning")
        critical_val = self.config.get("critical")
        
        warning_threshold = self.parse_threshold(warning_val, max_speed)
        critical_threshold = self.parse_threshold(critical_val, max_speed)
        
        total_speed = read_speed + write_speed
        
        if total_speed >= critical_threshold:
            level = "CRITICAL"
        elif total_speed >= warning_threshold:
            level = "WARNING"
            
        return [("I/O", level, f"R: {formatted_read} W: {formatted_write}")]
