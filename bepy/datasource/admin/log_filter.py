from django.contrib import admin
from ..models import LogFilter


@admin.register(LogFilter)
class LogFilterAdmin(admin.ModelAdmin):
    list_display = ['id', 'name', 'code', 'filter_type', 'compare_type', 'priority', 'updated_at']
    search_fields = ["name", "code"]