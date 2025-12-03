import unittest
from unittest.mock import MagicMock, patch
import time
import sys
import os

# Add src to path to allow imports
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../src')))

from tinymonitor.monitor import Monitor

class TestMonitorLogic(unittest.TestCase):
    def setUp(self):
        self.config = {
            "refresh": 1,
            "cooldown": 60,
            "cpu": {"enabled": False},
            "memory": {"enabled": False},
            "filesystem": {"enabled": False},
            "load": {"enabled": False},
            "alerts": {}
        }
        
        # Patch AlertManager to avoid real network calls or file loading
        # We patch where it is imported in monitor.py
        with patch('tinymonitor.monitor.AlertManager') as MockAlertManager:
            self.monitor = Monitor(self.config)
            # Replace the instance's alert_manager with a fresh mock for assertions
            self.monitor.alert_manager = MagicMock()

    def test_process_state_immediate(self):
        """Test that duration=0 triggers immediately"""
        # Duration = 0 -> Should return True immediately
        result = self.monitor.process_state("TEST_COMP", "CRITICAL", "99%", duration=0)
        self.assertTrue(result)

    def test_process_state_recovery(self):
        """Test that level=None clears the state"""
        # Set state first
        self.monitor.alert_states["TEST_COMP"] = {"level": "CRITICAL", "start_time": 123}
        
        # Recovery
        result = self.monitor.process_state("TEST_COMP", None, "0%", duration=0)
        self.assertFalse(result)
        self.assertNotIn("TEST_COMP", self.monitor.alert_states)

    @patch('time.time')
    def test_process_state_duration_logic(self, mock_time):
        """Test the wait duration logic"""
        mock_time.return_value = 1000
        
        # First detection
        result = self.monitor.process_state("TEST_COMP", "CRITICAL", "99%", duration=10)
        self.assertFalse(result, "Should wait for duration")
        self.assertIn("TEST_COMP", self.monitor.alert_states)
        self.assertEqual(self.monitor.alert_states["TEST_COMP"]["start_time"], 1000)
        
        # 5 seconds later
        mock_time.return_value = 1005
        result = self.monitor.process_state("TEST_COMP", "CRITICAL", "99%", duration=10)
        self.assertFalse(result, "Should still be waiting")
        
        # 10 seconds later (total)
        mock_time.return_value = 1010
        result = self.monitor.process_state("TEST_COMP", "CRITICAL", "99%", duration=10)
        self.assertTrue(result, "Should trigger after duration")

    @patch('time.time')
    def test_trigger_alert_cooldown(self, mock_time):
        """Test the cooldown logic (anti-spam)"""
        mock_time.return_value = 1000
        self.monitor.config["cooldown"] = 60
        
        # First alert
        self.monitor.trigger_alert("TEST_COMP", "CRITICAL", "99%")
        self.monitor.alert_manager.send_alert.assert_called_with("TEST_COMP", "CRITICAL", "99%")
        self.assertEqual(self.monitor.last_alert["TEST_COMP"], 1000)
        
        # Reset mock
        self.monitor.alert_manager.send_alert.reset_mock()
        
        # 30 seconds later (inside cooldown)
        mock_time.return_value = 1030
        self.monitor.trigger_alert("TEST_COMP", "CRITICAL", "99%")
        self.monitor.alert_manager.send_alert.assert_not_called()
        
        # 61 seconds later (after cooldown)
        mock_time.return_value = 1061
        self.monitor.trigger_alert("TEST_COMP", "CRITICAL", "99%")
        self.monitor.alert_manager.send_alert.assert_called_with("TEST_COMP", "CRITICAL", "99%")

    @patch('time.time')
    def test_trigger_alert_once_mode(self, mock_time):
        """Test the cooldown=-1 logic (alert once per incident)"""
        mock_time.return_value = 1000
        self.monitor.config["cooldown"] = -1
        
        # Setup state (needed for alert once logic)
        self.monitor.alert_states["TEST_COMP"] = {"level": "CRITICAL", "start_time": 1000}
        
        # First alert
        self.monitor.trigger_alert("TEST_COMP", "CRITICAL", "99%")
        self.monitor.alert_manager.send_alert.assert_called_once()
        self.assertEqual(self.monitor.last_alert["TEST_COMP"], 1000)
        
        self.monitor.alert_manager.send_alert.reset_mock()
        
        # Try again later - should be suppressed because last_alert >= start_time
        mock_time.return_value = 2000
        self.monitor.trigger_alert("TEST_COMP", "CRITICAL", "99%")
        self.monitor.alert_manager.send_alert.assert_not_called()
        
        # New problem occurs (reset state)
        self.monitor.alert_states["TEST_COMP"] = {"level": "CRITICAL", "start_time": 3000}
        mock_time.return_value = 3000
        
        # Should alert again
        self.monitor.trigger_alert("TEST_COMP", "CRITICAL", "99%")
        self.monitor.alert_manager.send_alert.assert_called_once()

if __name__ == '__main__':
    unittest.main()
