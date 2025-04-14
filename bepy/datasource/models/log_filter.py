from django.db import models
from .base import BaseDBModel


class LogFilter(BaseDBModel):
    """
    Model to store log filters.
    """
    code = models.CharField(max_length=255, unique=True)
    name = models.CharField(max_length=255)
    description = models.TextField(blank=True, null=True)
    filter_type = models.CharField(max_length=255)
    priority = models.IntegerField(default=0)
    compare_type = models.CharField(
        choices=[
            ("equal", "Equal"),
            ("startswith", "Startswith"),
            ("endswith", "Endswith"),
            ("include", "Include"),
            ("exclude", "Exclude"),
            ("search", "Search"),
        ],
        default="seconds",
        max_length=10,
        help_text="Interval unit for the collector config"
    )
    active = models.BooleanField(default=True)

    def __str__(self):
        return self.name