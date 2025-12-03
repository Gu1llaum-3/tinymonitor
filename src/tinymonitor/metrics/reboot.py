import os
from tinymonitor.metrics.provider import MetricProvider

class RebootMetric(MetricProvider):
    def check(self):
        reboot_required = False
        details = "OK"
        
        # Debian / Ubuntu / Mint standard
        # Check for the flag file created by apt/dpkg
        if os.path.exists('/var/run/reboot-required') or os.path.exists('/run/reboot-required'):
            reboot_required = True
            details = "System requires a reboot (updates installed)"
            
        # Future: Add RHEL/CentOS support (needs-restarting command)
        
        if reboot_required:
            # We use WARNING level for pending reboot
            return [("REBOOT", "WARNING", details)]
        
        return [("REBOOT", None, details)]
