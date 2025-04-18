from django_redis import get_redis_connection
from django.conf import settings
import redis
import logging


logger = logging.getLogger(__name__)


def request_elk_log_collector_reload_config():
    try:
        redis_conn: redis.Redis = get_redis_connection("default")
        res = redis_conn.publish('elkConfigChange', 1)
        print(f"request_elk_log_collector_reload_config {res=}")
        return True
    except Exception as e:
        logger.error(f"request_elk_log_collector_reload_config {e=}")
        return False