from django.contrib import admin
from ..models import ProxyConfig


@admin.register(ProxyConfig)
class ProxyConfigAdmin(admin.ModelAdmin):
    list_display = ['id', 'address', 'port', 'updated_at']
    search_fields = ["address"]
    