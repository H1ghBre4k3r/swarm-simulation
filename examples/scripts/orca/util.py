import json
import sys


def log(type: str, msg: str) -> None:
    """
    Logs a message to stderr.
    """
    val = {
        "action": type,
        "payload": msg
    }
    print(json.dumps(val), file=sys.stderr)
