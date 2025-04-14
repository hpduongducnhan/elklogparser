from django.db import models
from .base import BaseDBModel


class ElkCollectedLog(BaseDBModel):
    """
    Log collection model for ELK.
    """
    collector_code = models.CharField(max_length=100, verbose_name="Collector Code")
    log_id = models.CharField(max_length=100, unique=True)
    log_type = models.CharField(max_length=50, verbose_name="Log Type")
    log_raw = models.JSONField(default=dict, verbose_name="Log Raw")
    log_timestamp = models.FloatField(default=0)

    def __str__(self):
        return f"{self.log_id}[{self.collector_code}]"
    
class ElkGroupLog(BaseDBModel):
    """
    Log collection model for ELK.
    """
    collector_code = models.CharField(max_length=100, verbose_name="Collector Code")
    name = models.CharField(max_length=100, verbose_name="Name")
    key = models.CharField(max_length=100, unique=True, verbose_name="Key")
    logs = models.ManyToManyField(
        ElkCollectedLog,
        related_name="elk_group_logs",
        verbose_name="Logs"
    )


    def __str__(self):
        return f"{self.name}[{self.collector_code}]"