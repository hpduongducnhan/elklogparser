from django.contrib import admin
from ..models import ElkCollectorConfig, ElkCollectResult


@admin.register(ElkCollectorConfig)
class ElkCollectorConfigAdmin(admin.ModelAdmin):
    list_display = ['code', 'name', 'active', 'last_run_at', 'updated_at', ]
    search_fields = ["name", "code"]
    filter_horizontal = ['elk_queries', 'log_filters']



@admin.register(ElkCollectResult)
class ElkCollectResultAdmin(admin.ModelAdmin):
    list_display = ['id', "collector__code", "collector__name", 'success', 'status', 'finish_at', 'query_from_datetime', 'query_to_datetime', ]
    search_fields = ["collector__code", "collector__name"]
    list_filter = ['success', 'status', 'collector__code']

    def get_readonly_fields(self, request, obj=None):
        return [field.name for field in self.model._meta.fields]
