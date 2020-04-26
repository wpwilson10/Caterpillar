import pysbd

with open('./test/foxnews1clean.txt', 'r') as file:
    text = file.read()

    seg = pysbd.Segmenter(language="en", clean=False)
    one = seg.segment(text)

    seg2 = pysbd.Segmenter(language="en", clean=True)
    two = seg2.segment(text)

    i = 0
    length = 0
    for each in one:
        i += 1
        length += len(each)
        print(i, each)

    print("**************", length)

    i = 0
    length = 0
    for each in two:
        i += 1
        length += len(each)
        print(i, each)

    print("--------------", length)