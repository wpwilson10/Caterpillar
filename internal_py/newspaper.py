#!/usr/bin/env python3
'''Uses the newspaper3k library to collect and parse article data'''
import os
from concurrent import futures

import grpc
from newspaper import Article, Config

from .setup import is_open, APP_LOG
from . import caterpillar_pb2
from . import caterpillar_pb2_grpc

def extract_newspaper(link, config):
    '''extract_newspaper parses an article from the given link'''
    try:
        # setup article
        article = Article(link, config=config)

        # extract article
        article.download()
        article.parse()

        # create reply with information to return
        response = caterpillar_pb2.NewspaperReply()

        if article is None:
            # case where we didn't get an article
            response.link = link
            return response
        elif len(article.text) < 10 or len(article.title) < 3:
            # case where there is not enough article information
            response.link = link
            return response
        # at this point, we have information worth returning
        # published_date is commonly not returned
        if article.publish_date is None:
            response.link = link
            response.title = article.title
            response.text = article.text
            response.canonical = article.canonical_link
            response.pubdate = ""
        else:
            # got eveything
            response.link = link
            response.title = article.title
            response.text = article.text
            response.canonical = article.canonical_link
            response.pubdate = article.publish_date.isoformat()

        # add repeated authors field
        if len(article.authors) > 0:
            # ignore error, field is there
            response.authors.extend(article.authors)

        return response
    except (Exception) as error:
        APP_LOG.error(error)

class NewspaperServicer(caterpillar_pb2_grpc.NewspaperServicer):
    """Provides methods that implement functionality of Newspaper server."""

    def __init__(self):
        # setup newspaper configuration
        self.config = Config()
        self.config.browser_user_agent = os.getenv("PY_NEWSPAPER_USER_AGENT")
        self.config.memoize_articles = False
        self.config.fetch_images = False
        self.config.follow_meta_refresh = True

    # Main server application call
    def Request(self, request, context):
        # call newspaper library
        response = extract_newspaper(link=request.link, config=self.config)
        # check if we failed
        if response is None:
            context.set_code(grpc.StatusCode.INTERNAL)
            return caterpillar_pb2.NewspaperReply()
        # otherwise return data
        return response

def run():
    '''Run newspaper application server'''
     # port to use with app
    port = int(os.getenv("PY_NEWSPAPER_PORT"))

    # if port is not in use, start app
    if is_open("localhost", port):
        # setup server
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        caterpillar_pb2_grpc.add_NewspaperServicer_to_server(
            NewspaperServicer(), server)
        server.add_insecure_port(os.getenv("NEWSPAPER_HOST"))
        # run
        server.start()
        APP_LOG.info("Newspaper server starting")
        # timeout after 3 days (arbitrary)
        server.wait_for_termination(timeout=60.0*60*24*3)
