from django.contrib import admin
from ..models import ElkCollectedLog, ElkGroupLog


@admin.register(ElkCollectedLog)
class ElkCollectedLogAdmin(admin.ModelAdmin):
    list_display = ('collector_code', 'log_id', 'log_type', 'log_timestamp', 'created_at')
    search_fields = ("collector_code", "log_id")
    list_filter = ('collector_code', 'log_type',)


@admin.register(ElkGroupLog)
class ElkGroupLogAdmin(admin.ModelAdmin):
    list_display = ('collector_code', 'name', 'key', 'updated_at', 'created_at')
    search_fields = ("collector_code", "name", 'key')