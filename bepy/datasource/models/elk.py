from django.db import models
from .base import BaseDBModel
from .utils import request_elk_log_collector_reload_config


class ElkQuery(BaseDBModel):
    name = models.CharField(max_length=255, help_text="Name of the query")
    code = models.CharField(max_length=255, unique=True, help_text="Code of the query")
    description = models.TextField(help_text="Description of the query")
    index = models.CharField(max_length=255, help_text="Index to be queried")
    query = models.JSONField(default=dict, help_text="Query to be executed")
    
    class Meta:
        ordering = ["-updated_at"]
        indexes = [
            models.Index(fields=["code"]),
        ]

    def __str__(self):
        return f"{self.name}-{self.index}[{self.code}]"
    
    def save(self, *args, **kwargs):
        request_elk_log_collector_reload_config()
        return super().save(*args, **kwargs)

class ElkConfig(BaseDBModel):
    """
    Base model for all database models.
    """
    code = models.CharField(max_length=255, unique=True, help_text="Code of the ELK server")
    host = models.CharField(max_length=255, help_text="Host of the ELK server")
    port = models.IntegerField(help_text="Port of the ELK server")
    username = models.CharField(max_length=255, null=True, blank=True, help_text="Username for ELK server")
    password = models.CharField(max_length=255, null=True, blank=True, help_text="Password for ELK server")
    options = models.JSONField(default=dict, help_text="Additional options for ELK server")

    proxy_enable = models.BooleanField(default=False, help_text="Is proxy enabled?")
    proxy = models.ForeignKey(
        "ProxyConfig",
        on_delete=models.CASCADE,
        null=True,
        blank=True,
        help_text="Proxy config for ELK server"
    )

    class Meta:
        ordering = ["-updated_at"]
        indexes = [
            models.Index(fields=["code"]),
        ]

    def __str__(self):
        return f"{self.__class__.__name__}[{self.code}]"
    
    def save(self, *args, **kwargs):
        request_elk_log_collector_reload_config()
        return super().save(*args, **kwargs)