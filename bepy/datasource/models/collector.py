from django.db import models
from .base import BaseDBModel


class ElkCollectorConfig(BaseDBModel):
    """
    Base model for all database models.
    """
    code = models.CharField(max_length=255, unique=True, help_text="Code of the Collector")
    name = models.CharField(max_length=255, help_text="Host of the ELK server")
    elk_config = models.ForeignKey(
        "ElkConfig",
        on_delete=models.CASCADE,
        related_name="collector_configs",
        help_text="The ELK client config this config belongs to",
    )
    elk_queries = models.ManyToManyField(
        "ElkQuery",
        related_name="collector_configs",
        null=True,
        blank=True,
        help_text="The ELK queries this config belongs to",
    )
    log_filters = models.ManyToManyField(
        "LogFilter",
        related_name="collector_configs",
        null=True,
        blank=True,
        help_text="The log filters this config belongs to",
    )

    interval_value = models.IntegerField(
        default=30,
        help_text="Interval value for the collector config"
    )
    interval_unit = models.CharField(
        choices=[
            ("seconds", "Seconds"),
            ("minutes", "Minutes"),
            ("hours", "Hours"),
        ],
        default="seconds",
        max_length=10,
        help_text="Interval unit for the collector config"
    )

    last_run_at = models.DateTimeField(
        null=True,
        blank=True,
        help_text="Last time the collector config was run"
    )

    active = models.BooleanField(
        default=True,
        help_text="Is this collector config active?"
    )
    
    def __str__(self):
        return f"{self.__class__.__name__}({self.pk})"
    

class ElkCollectResult(BaseDBModel):
    collector = models.ForeignKey(
        ElkCollectorConfig,
        on_delete=models.CASCADE,
        related_name="collect_results",
        help_text="The collector config this result belongs to",
    )
    start_at = models.DateTimeField(
        null=True, blank=True,
        help_text="Start time of the collection"
    )
    finish_at = models.DateTimeField(
        null=True, blank=True,
        help_text="Finish time of the collection"
    )
    success = models.BooleanField(
        default=False,
        help_text="Was the collection successful?"
    )
    status = models.CharField(
        max_length=255,
        help_text="Status of the collection result"
    )
    detail = models.JSONField(
        default=dict,
        help_text="Details of the collection result"
    )

    def __str__(self):
        return f'{self.__class__.__name__}({self.pk})'
    