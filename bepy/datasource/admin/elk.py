from django.contrib import admin
from ..models import ElkConfig, ElkQuery


@admin.register(ElkConfig)
class ElkConfigAdmin(admin.ModelAdmin):
    list_display = ['id', 'code', 'host', 'port']
    search_fields = ["code", "host"]


@admin.register(ElkQuery)
class ElkQueryAdmin(admin.ModelAdmin):
    list_display = ['id', 'name', 'code', 'updated_at']
    search_fields = ["name", "code"]