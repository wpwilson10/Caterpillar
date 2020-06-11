#!/usr/bin/env python3
"""Text handles string text processing"""

import grpc
from gensim.summarization import summarize, keywords
from . import caterpillar_pb2

def sentences(self, request, context):
    """sentences parses text blocks into individual sentences"""
    # call sentence parser
    sent = self.seg.segment(request.text)
    response = caterpillar_pb2.SentenceReply()
    # check if we failed
    if sent is None or len(sent) < 1:
        context.set_code(grpc.StatusCode.INTERNAL)
        return response
    # add repeated sentences field
    response.sentences.extend(sent)
    return response

def summary(request, context):
    """summary returns a summarization and keywords of the given text string"""
    # uses gensim library

    # figure out a reasonable ratio of summarization
    ratio = 1.0
    if len(request.text) <= 280:
        pass  # 280 is the max size of a tweet, so don't summarize
    elif len(request.text) <= 1000:
        ratio = 0.7
    elif len(request.text) <= 3000:
        ratio = 0.5
    elif len(request.text) <= 10000:
        ratio = 0.3
    else:
        ratio = 0.1

    # call sentence parser
    summ = summarize(request.text, ratio=ratio)
    keys = keywords(request.text, words=20, split=True, lemmatize=True)
    response = caterpillar_pb2.SummaryReply()
    # check if we failed
    if summ is None and keys is None:
        context.set_code(grpc.StatusCode.INTERNAL)
        return response
    # add data fields
    response.summary = summ
    response.keywords.extend(keys)
    return response
