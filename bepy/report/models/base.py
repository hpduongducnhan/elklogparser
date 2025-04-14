from django.db import models

# Create your models here.
class BaseDBModel(models.Model):
    """
    Base model for all database models.
    """
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        abstract = True

    def __str__(self):
        return f"{self.__class__.__name__}({self.pk})"