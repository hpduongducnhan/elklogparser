# -*- coding: utf-8 -*-
import threading


_USE_WRITE = 'use_write'


class DbPrimaryReplicaSelector:
    """[summary]"""

    def _singleton_init(self, **kwargs):
        """[summary]"""
        self._local_db_selector = threading.local()
        self.set_write()

    def reset(self):
        setattr(self._local_db_selector, _USE_WRITE, False)

    def handle_new_request(self):
        self.set_read()
        # print(f'db selector handle new request: {getattr(self._local_db_selector, _USE_WRITE, False)}')

    def set_read(self):
        setattr(self._local_db_selector, _USE_WRITE, False)

    def set_write(self):
        setattr(self._local_db_selector, _USE_WRITE, True)

    def is_write(self) -> bool:
        # print(f'db selector is write: {getattr(self._local_db_selector, _USE_WRITE, False)}')
        return getattr(self._local_db_selector, _USE_WRITE, False)


db_selector: DbPrimaryReplicaSelector = DbPrimaryReplicaSelector()
