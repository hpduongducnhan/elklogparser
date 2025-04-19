from django_redis import get_redis_connection
from typing import List
import redis
import logging
import json

logger = logging.getLogger(__name__)


def request_elk_log_collector_reload_config(collector_codes: List[str]) -> bool:
    try:
        redis_conn: redis.Redis = get_redis_connection("default")
        res = redis_conn.publish('elkCollectorConfigChange', json.dumps({'codes': collector_codes}))
        print(f"request_elk_log_collector_reload_config {res=}")
        return True
    except Exception as e:
        logger.error(f"request_elk_log_collector_reload_config {e=}")
        return False