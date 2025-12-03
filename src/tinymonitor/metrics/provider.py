from abc import ABC, abstractmethod

class MetricProvider(ABC):
    def __init__(self, name, config):
        self.name = name
        self.config = config

    @abstractmethod
    def check(self):
        """
        Executes the check.
        Must return a list of tuples: (component_name, level, value)
        
        - component_name : Name of the component (ex: "CPU", "DISK:/")
        - level : "CRITICAL", "WARNING" or None
        - value : Formatted value (ex: "85%")
        """
        pass
