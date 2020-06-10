#!/usr/bin/env python3
'''Uses the newspaper3k library to collect and parse article data'''

import grpc
from newspaper import Article

from . import caterpillar_pb2
from .setup import APP_LOG

def newspaper(self, request, context):
    '''newspaper handles caterpillar calls to newspaper3k library'''
    # call newspaper library
    response = extract(link=request.link, config=self.config)
    # check if we failed
    if response is None:
        context.set_code(grpc.StatusCode.INTERNAL)
        return caterpillar_pb2.NewspaperReply()
    # otherwise return data
    return response

def extract(link, config):
    '''extract parses an article from the given link'''
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
