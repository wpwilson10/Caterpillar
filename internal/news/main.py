#!/usr/bin/env python3

from http.server import HTTPServer, BaseHTTPRequestHandler
from socketserver import ThreadingMixIn
from jsonrpcserver import method, dispatch
from newspaper import Article, Config
from dotenv import load_dotenv
from pathlib import Path
from datetime import date
import threading
import json
import logging
import os
import socket
import time

# global logger
logger = 0

# extractNewspaper parses an article from the given link
# input parameter name must match the arguments from client exactly
@method
def extractNewspaper(Link):
    try:
        # setup configuration
        config = Config()
        config.browser_user_agent = os.getenv("PY_NEWSPAPER_USER_AGENT")
        config.memoize_articles = False
        config.fetch_images = False
        config.follow_meta_refresh = True

        # setup article
        article = Article(Link, config = config)

        # extract article
        article.download()
        article.parse()

        # published_date is commonly not returned
        if article.publish_date is None:
            response = {
                "Title": article.title,
                "Text": article.text,
                "Authors": article.authors,
                "Canonical": article.canonical_link,
                "PubDate": ""
            }

            return response

        # default case, we got all data
        response = {
            "Title": article.title,
            "Text": article.text,
            "Authors": article.authors,
            "Canonical": article.canonical_link,
            "PubDate": article.publish_date.isoformat()
        }

        return response
    
    except (Exception) as error :
        logger.error(error)

# What the HTTP server calls to process requests
class Handler(BaseHTTPRequestHandler):
    def do_POST(self):
        # Read the request
        request = self.rfile.read(int(self.headers["Content-Length"])).decode()
        # Send it to the requested method (Should just be newspaper for now)
        response = dispatch(request)
        # Return response
        self.send_response(response.http_status)
        self.send_header("Content-type", "application/json")
        self.end_headers()
        self.wfile.write(str(response).encode())

# ThreadedHTTPServer based on the blog post:
# https://pymotw.com/2/BaseHTTPServer/index.html#module-BaseHTTPServer
class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
    """Handle requests in a separate thread."""

# Checks if the given host and port address is in use
def isOpen(ip, port):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(3)
    result = s.connect_ex((ip, port))
    s.close()

    # connection should be refused because nothing is checking it
    if result == 0: 
        return False
    elif result == 111:
        # Ubuntu: errno 111: Connection refused
        return True
    elif result == 10061:
        # [WinError 10061] No connection could be made because the target machine actively refused it
        return True
    
    return False

# configure logger
def setupLogger():
    global logger
    logger = logging.getLogger('NewspaperLogger')
    logger.setLevel(logging.INFO)
    # create file handler which logs
    log_path = os.getenv("LOG_FILEPATH") +"pyNews.log"
    fh = logging.FileHandler(log_path)
    fh.setLevel(logging.DEBUG)
    # create formatter and add it to the handlers
    formatter = logging.Formatter('%(asctime)s -  %(levelname)s - %(name)s - %(message)s')
    fh.setFormatter(formatter)
    # add the handlers to logger
    logger.addHandler(fh)

# Run application
def mainNewspaper():
    # load configuration
    env_path = Path('.') / 'configs' / '.env'
    load_dotenv(dotenv_path=env_path)
    # port to use with app
    port = int(os.getenv("PY_NEWSPAPER_PORT"))
    # setup logger
    setupLogger()
    # if port is not in use, start app
    if isOpen("localhost", port):
        logger.info("Python Newspaper server starting")
        server = ThreadedHTTPServer(('localhost', port), Handler)
        server.serve_forever()
        logger.info("Server stopping")

# startup server
if __name__ == "__main__":
    mainNewspaper()
    