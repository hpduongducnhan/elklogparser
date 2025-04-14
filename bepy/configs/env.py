# -*- coding: utf-8 -*-
import os
from pathlib import Path
from typing import Dict, List, Optional

from pydantic_settings import BaseSettings, SettingsConfigDict  # type: ignore

_current_path = Path(os.getcwd())


class EnvSettings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=os.path.join(_current_path, 'envs', '.env'), env_file_encoding='utf-8'
    )

    DEBUG: bool = True
    BASE_DIR: str = str(_current_path)

    # DJANGO -----------------------------------------------------
    DJANGO_SECRET_KEY: Optional[str] = (
        "django-insecure-dfx_n0$#x$53$k32px9na9@(_fu&og)xz-cmcot7a9-6-ej()="
    )
    DJANGO_DEBUG: Optional[bool] = True
    DJANGO_DB_SELECT_WRITE_MODELS: List[str] = []

    # GUNICORN ---------------------------------------------------
    GUNICORN_HOST: str = '0.0.0.0'
    GUNICORN_PORT: str = '5000'
    GUNICORN_BIND_PATH: Optional[str] = None
    GUNICORN_WORKER_CONCURRENCY: Optional[int] = os.cpu_count() * 2
    GUNICORN_GRACEFUL_TIMEOUT: Optional[int] = 120
    GUNICORN_TIMEOUT: Optional[int] = 500
    GUNICORN_KEEP_ALIVE: Optional[int] = 120
    GUNICORN_LIMIT_REQUEST_LINE: Optional[int] = 0

    GUNICORN_ACCESS_LOG: Optional[str] = '-'  # default log to console
    GUNICORN_ERROR_LOG: Optional[str] = '-'  # default log to console

    # APP --------------------------------------------------------
    APP_PROXY_ENABLE: Optional[bool] = False
    APP_PROXY_ADDRESS: Optional[str] = 'http://proxy.com:80'
    APP_PROXY_BYPASS: Optional[str] = 'localhost,127.0.0.1,172.27.230.15,172.27.230.30'

    APP_DB_MASTER_HOST: Optional[str] = ''
    APP_DB_MASTER_PORT: Optional[int] = 0
    APP_DB_MASTER_ENGINE: Optional[str] = ''
    APP_DB_MASTER_NAME: Optional[str] = ''
    APP_DB_MASTER_USER: Optional[str] = ''
    APP_DB_MASTER_PASSWORD: Optional[str] = ''

    APP_DB_SLAVE_HOST: Optional[str] = ''
    APP_DB_SLAVE_PORT: Optional[int] = 0
    APP_DB_SLAVE_ENGINE: Optional[str] = ''
    APP_DB_SLAVE_NAME: Optional[str] = ''
    APP_DB_SLAVE_USER: Optional[str] = ''
    APP_DB_SLAVE_PASSWORD: Optional[str] = ''

    APP_REDIS_HOST: Optional[str] = 'localhost'
    APP_REDIS_PORT: Optional[int] = 6379
    APP_REDIS_DB: Optional[int] = 10
    APP_REDIS_USER: Optional[str] = ''
    APP_REDIS_PASSWORD: Optional[str] = 'password'
    

    def get_db_conf(self, prefix) -> Dict:
        res = {
            'ENGINE': getattr(self, f'{prefix}_ENGINE'),
            'NAME': getattr(self, f'{prefix}_NAME'),
            'USER': getattr(self, f'{prefix}_USER'),
            'PASSWORD': getattr(self, f'{prefix}_PASSWORD'),
            'HOST': getattr(self, f'{prefix}_HOST'),
            'PORT': getattr(self, f'{prefix}_PORT'),
        }
        if getattr(self, f'{prefix}_ENGINE') == 'django.db.backends.postgresql':
            res.update(
                {
                    'ATOMIC_REQUESTS': True,
                    'CONN_MAX_AGE': 0,
                    'CONN_HEALTH_CHECKS': True,
                }
            )
            res.update(
                {
                    'OPTIONS': {
                        # "pool": True,
                        # 'adapter': 'psycopg'
                    },
                }
            )
        elif getattr(self, f'{prefix}_ENGINE') == 'django.db.backends.mysql':
            res.update(
                {
                    'OPTIONS': {
                        'charset': 'utf8mb4',
                        'init_command': "SET sql_mode='STRICT_TRANS_TABLES'",
                    },
                }
            )
        return res

    def _get_sqlite3(self, name: str) -> Dict:
        return {
            'ENGINE': 'django.db.backends.sqlite3',
            'NAME': os.path.join(self.BASE_DIR, 'data', f'{name}.sqlite3'),
        }

    def get_db_conf_master(self, prefix: str = 'APP_DB_MASTER') -> Dict:
        if not self.APP_DB_MASTER_ENGINE or not self.APP_DB_MASTER_HOST:
            return self._get_sqlite3('master')
        return self.get_db_conf(prefix)

    def get_db_conf_slave(self, prefix: str = 'APP_DB_SLAVE') -> Dict:
        if not self.APP_DB_SLAVE_ENGINE or not self.APP_DB_SLAVE_HOST:
            return self.get_db_conf_master()
        return self.get_db_conf(prefix)

    def set_proxy(self):
        if not self.APP_PROXY_ENABLE:
            return
        os.environ['http_proxy'] = self.APP_PROXY_ADDRESS
        os.environ['https_proxy'] = self.APP_PROXY_ADDRESS
        no_proxy = os.environ.get('no_proxy', '')
        if no_proxy:
            for addr in self.APP_PROXY_BYPASS.split(','):
                if addr not in no_proxy:
                    no_proxy = no_proxy + ',' + addr
        else:
            no_proxy = self.APP_PROXY_BYPASS

        os.environ['no_proxy'] = no_proxy

    def get_redis_url(self) -> str:
        self.CELERY_BROKER_URL = f'redis://{self.APP_REDIS_USER}:{self.APP_REDIS_PASSWORD}@{
            self.APP_REDIS_HOST}:{self.APP_REDIS_PORT}/{self.APP_REDIS_DB}'
        return f'redis://{self.APP_REDIS_USER}:{self.APP_REDIS_PASSWORD}@{self.APP_REDIS_HOST}:{self.APP_REDIS_PORT}/{self.APP_REDIS_DB}'


env = EnvSettings()
