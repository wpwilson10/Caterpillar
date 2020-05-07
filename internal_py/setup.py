#!/usr/bin/env python3

import logging
import os
import socket

# global logger
logger = 0

# configure logger
def setup_logger():
    global logger
    logger = logging.getLogger('NewspaperLogger')
    logger.setLevel(logging.INFO)
    # create file handler which logs
    log_path = os.getenv("LOG_FILEPATH") +"pyNews.log"
    file_handler = logging.FileHandler(log_path)
    file_handler.setLevel(logging.DEBUG)
    # create formatter and add it to the handlers
    formatter = logging.Formatter('%(asctime)s -  %(levelname)s - %(name)s - %(message)s')
    file_handler.setFormatter(formatter)
    # add the handlers to logger
    logger.addHandler(file_handler)

# Checks if the given host and port address is in use
def is_open(ip, port):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.settimeout(3)
    result = sock.connect_ex((ip, port))
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
