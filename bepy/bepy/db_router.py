# -*- coding: utf-8 -*-
from configs.env import env
from django.db import connections
from .db_router_selector import db_selector


class PrimaryReplicaRouter:
    def __init__(self) -> None:
        self._has_replica = False
        if connections['replica']:
            self._has_replica = True
            print('PrimaryReplicaRouter: replica is available')

    def db_for_read(self, model, **hints):
        """
        Reads go to a randomly-chosen replica.
        """
        if (
            not self._has_replica
            or model._meta.db_table in env.DJANGO_DB_SELECT_WRITE_MODELS
            or db_selector.is_write()
        ):
            return 'default'

        return 'replica'

    def db_for_write(self, model, **hints):
        """
        Writes always go to primary.
        """
        db_selector.set_write()
        return "default"

    def allow_relation(self, obj1, obj2, **hints):
        """
        Relations between objects are allowed if both objects are
        in the primary/replica pool.
        """
        db_set = {"default", "replica"}
        if obj1._state.db in db_set and obj2._state.db in db_set:
            return True
        return None

    def allow_migrate(self, db, app_label, model_name=None, **hints):
        """
        All non-auth models end up in this pool.
        """
        return True
