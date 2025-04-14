from django.db import models
from .base import BaseDBModel


class ProxyConfig(BaseDBModel):
    """
    Base model for all database models.
    """
    address = models.CharField(max_length=255, help_text="Host of the ELK server")
    port = models.IntegerField(null=True, blank=True, help_text="Port of the ELK server")
    username = models.CharField(max_length=255, null=True, blank=True, help_text="Username for ELK server")
    password = models.CharField(max_length=255, null=True, blank=True, help_text="Password for ELK server")
    options = models.JSONField(default=dict, help_text="Additional options for ELK server")

    def __str__(self):
        return f"{self.__class__.__name__}({self.pk})"