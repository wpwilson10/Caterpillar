#!/usr/bin/env python3
'''Sets up gRPC server for Caterpillar applications'''
import os
from concurrent import futures

import grpc
import pysbd
from newspaper import Config

from .setup import is_open, APP_LOG
from .newspaper import extract_newspaper
from . import caterpillar_pb2
from . import caterpillar_pb2_grpc

class CaterpillarServicer(caterpillar_pb2_grpc.CaterpillarServicer):
    """Provides methods that implement applications for Caterpillar."""

    def __init__(self):
        # setup newspaper configuration
        self.config = Config()
        self.config.browser_user_agent = os.getenv("PY_NEWSPAPER_USER_AGENT")
        self.config.memoize_articles = False
        self.config.fetch_images = False
        self.config.follow_meta_refresh = True
        # setup sentence parser
        self.seg = pysbd.Segmenter(language="en", clean=False)

    # Newspaper3k article extraction
    def Newspaper(self, request, context):
        # call newspaper library
        response = extract_newspaper(link=request.link, config=self.config)
        # check if we failed
        if response is None:
            context.set_code(grpc.StatusCode.INTERNAL)
            return caterpillar_pb2.NewspaperReply()
        # otherwise return data
        return response

    # Sentence parsing
    def Sentences(self, request, context):
        # call sentence parser
        sentences = self.seg.segment(request.text)
        response = caterpillar_pb2.SentenceReply()
        # check if we failed
        if sentences is None or len(sentences) < 1:
            context.set_code(grpc.StatusCode.INTERNAL)
            return response
        # add repeated sentences field
        response.sentences.extend(sentences)
        return response

def run():
    '''Run caterpillar application server'''
     # port to use with app
    port = int(os.getenv("PY_CATERPILLAR_PORT"))

    # if port is not in use, start app
    if is_open("localhost", port):
        # setup server
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=16))
        caterpillar_pb2_grpc.add_CaterpillarServicer_to_server(
            CaterpillarServicer(), server)
        server.add_insecure_port(os.getenv("PY_CATERPILLAR_HOST"))
        # run
        server.start()
        APP_LOG.info("Caterpillar server starting")
        # timeout after 3 days (arbitrary)
        server.wait_for_termination(timeout=60.0*60*24*3)
