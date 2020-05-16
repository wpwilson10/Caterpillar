#!/usr/bin/env python3

import os
from concurrent import futures

import pysbd
import grpc

from .setup import is_open, APP_LOG
from . import caterpillar_pb2
from . import caterpillar_pb2_grpc

class TextServicer(caterpillar_pb2_grpc.TextServicer):
    """Provides methods that implement functionality of Text server."""

    # Don't need to initialize anything
    def __init__(self):
        # call pysbd library - https://github.com/nipunsadvilkar/pySBD
        self.seg = pysbd.Segmenter(language="en", clean=False)

    # Main server application call
    def Sentences(self, request, context):
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
    '''Run text application server'''
     # port to use with app
    port = int(os.getenv("PY_TEXT_PORT"))

    # if port is not in use, start app
    if is_open("localhost", port):
        # setup server
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        caterpillar_pb2_grpc.add_TextServicer_to_server(
            TextServicer(), server)
        server.add_insecure_port(os.getenv("TEXT_HOST"))
        # run
        server.start()
        APP_LOG.info("Text server starting")
        # timeout after 3 days (arbitrary)
        server.wait_for_termination(timeout=60.0*60*24*3)
