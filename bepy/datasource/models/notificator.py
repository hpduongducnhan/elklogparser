from django.db import models
from .base import BaseDBModel

# class Notificator(BaseDBModel):
#     name = models.CharField(max_length=255, unique=True)
#     description = models.TextField(max_length=255, blank=True)
#     ntf_type = models.CharField(
#         choices=[
#             ("api", "Api"),
#             ("chat", "Chat"),
#             ("email", "Email"),
#         ],
#         default="chat",
#         max_length=20,
#     )
#     ntf_sub_type = models.CharField(
#         choices=[
#             ("telegram", "Telegram"),
#             ("slack", "Slack"),
#             ("email", "Email"),
#         ],
#         default="telegram",
#         max_length=20,
#     )