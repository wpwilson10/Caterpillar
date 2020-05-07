#!/usr/bin/env python3

import os
from http.server import HTTPServer, BaseHTTPRequestHandler
from socketserver import ThreadingMixIn
from jsonrpcserver import method, dispatch
from newspaper import Article, Config
from setup import is_open, logger

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
    except (Exception) as error:
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


# Run application
def mainNewspaper():
     # port to use with app
    port = int(os.getenv("PY_NEWSPAPER_PORT"))

    # if port is not in use, start app
    if is_open("localhost", port):
        logger.info("Python Newspaper server starting")
        server = ThreadedHTTPServer(('localhost', port), Handler)
        server.serve_forever()
        logger.info("Server stopping")
