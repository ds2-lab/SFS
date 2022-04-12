# Copyright (c) 2019 Princeton University
#
# This source code is licensed under the MIT license found in the
# LICENSE file in the root directory of this source tree.
import json
import markdown
import base64
import os
import sys
def f(n):
    try:
        with open('/openpiton-readme.json') as f:
            data = json.load(f)
    except:
        return {'Error' : 'Possibly lacking markdown parameter in request.'}
    text = data["markdown"]
    decoded_text = base64.b64decode(text.encode()).decode()
    return decoded_text

if __name__ == "__main__":
    n = sys.argv[1]
    f(n)
