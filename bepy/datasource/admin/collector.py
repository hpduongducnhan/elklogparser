from django.contrib import admin
from ..models import ElkCollectorConfig, ElkCollectResult


@admin.register(ElkCollectorConfig)
class ElkCollectorConfigAdmin(admin.ModelAdmin):
    list_display = ['id', 'name', 'code', 'last_run_at', 'updated_at', ]
    search_fields = ["name", "code"]
    filter_horizontal = ['elk_queries', 'log_filters']



@admin.register(ElkCollectResult)
class ElkCollectResultAdmin(admin.ModelAdmin):
    list_display = ['id', "collector__code", "collector__name", 'success', 'status', 'finish_at']
    search_fields = ["collector__code", "collector__name"]