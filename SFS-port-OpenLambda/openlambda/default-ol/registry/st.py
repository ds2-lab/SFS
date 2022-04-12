import sys
import json
from textblob import TextBlob
import os
import nltk
nltk.data.path.append("/nltk_data")
def f(n):
    words = "PATRIOTS BRAVE & BOLD 4400 University Drive, Fairfax, Virginia 22030 © 2020 George Mason University – Call: +1 (703) 993-1000".split()
    for word in words:
        try:
            analyse = TextBlob(word)
        except:
            return {'Error' : 'Input parameters should include a string to sentiment analyse.'}

        sentences = len(analyse.sentences)

        retVal = {}

        retVal["subjectivity"] = sum([sentence.sentiment.subjectivity for sentence in analyse.sentences]) / sentences
        retVal["polarity"] = sum([sentence.sentiment.polarity for sentence in analyse.sentences]) / sentences
        retVal["sentences"] = sentences

    return retVal

if __name__ == "__main__":
    event = sys.argv[1]
    f(event)
