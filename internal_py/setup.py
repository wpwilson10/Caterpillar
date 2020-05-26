#!/usr/bin/env python3
"""Setup contains tool relevant to program initialization"""

import logging
import os
import socket

# global logger, clear out so we can set it to use our logger
APP_LOG = logging.getLogger('ApplicationLogger')

def setup_logger():
    """Configures and initializes a logger"""

    global APP_LOG
    APP_LOG.setLevel(logging.INFO)
    # create file handler which logs
    log_path = os.getenv("LOG_FILEPATH") +"pylog.log"
    file_handler = logging.FileHandler(log_path)
    file_handler.setLevel(logging.INFO)
    # create formatter and add it to the handlers
    formatter = logging.Formatter('%(asctime)s -  %(levelname)s - %(name)s - %(message)s')
    file_handler.setFormatter(formatter)
    # add the handlers to logger
    APP_LOG.addHandler(file_handler)


def is_open(ip_address, port):
    """Checks if the given host and port address is in use"""

    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.settimeout(3)
    result = sock.connect_ex((ip_address, port))
    sock.close()

    # connection should be refused because nothing is checking it
    if result == 0:
        return False
    if result == 111:
        # Ubuntu: errno 111: Connection refused
        return True
    if result == 10061:
        # [WinError 10061]
        # No connection could be made because the target machine actively refused it
        return True
    return False
