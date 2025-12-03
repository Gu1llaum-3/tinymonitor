from abc import ABC, abstractmethod
import logging

logger = logging.getLogger("tinymonitor")

class AlertProvider(ABC):
    def __init__(self, name, config):
        self.name = name
        self.config = config

    def _normalize_component_name(self, component):
        """
        Converts the technical component name to a configuration key.
        Ex: 'DISK:/Volumes/Data' -> 'filesystem'
        Ex: 'LOAD' -> 'load'
        Ex: 'cpu' -> 'cpu'
        """
        if component.startswith("DISK:"):
            return "filesystem"
        if component == "LOAD":
            return "load"
        return component.lower()

    def should_send(self, component, level):
        """
        Checks if this provider should send an alert for this component and level.
        """
        # 1. Global check
        if not self.config.get("enabled", True):
            return False

        # 2. Rules retrieval
        rules = self.config.get("rules")
        
        # If no rules defined, use old system (or all by default)
        if not rules:
            accepted_levels = self.config.get("levels", ["WARNING", "CRITICAL"])
            return level in accepted_levels

        # 3. Component logic
        config_key = self._normalize_component_name(component)
        
        # Look for specific rule, else default, else refuse for safety
        allowed_levels = rules.get(config_key, rules.get("default", []))
        
        return level in allowed_levels

    @abstractmethod
    def send(self, component, level, value, title, message):
        """Abstract method to be implemented by plugins."""
        pass
