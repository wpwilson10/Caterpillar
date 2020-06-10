#!/usr/bin/env python3
'''Sets up gRPC server for Caterpillar applications'''
import os
from concurrent import futures

import grpc
import pysbd
from newspaper import Config
from transformers import pipeline

from .setup import is_open, APP_LOG
from .newspaper import newspaper
from .text import sentences, summary, feature_extraction
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
        # setup NLP pipelines, currently using XLNET
        model = 'xlnet-base-cased'
        self.features = pipeline('feature-extraction', model=model, tokenizer=model, device=0)
        self.sentiment = pipeline('sentiment-analysis', device=0)

    # Newspaper3k article extraction
    def Newspaper(self, request, context):
        return newspaper(self, request, context)

    # Sentence parsing
    def Sentences(self, request, context):
        return sentences(self, request, context)

    # Text summarization
    def Summary(self, request, context):
        feature_extraction(self, request, context)
        return summary(request, context)

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
