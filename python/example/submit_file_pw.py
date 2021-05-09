# Copyright (C) 2021 Hatching B.V
# All rights reserved.

import io
import sys

import triage

url = "https://api.tria.ge"
token = "<YOUR-APIKEY-HERE>"

c = triage.Client(token, root_url=url)
f = open(sys.argv[1], "rb")
r = c.submit_sample_file("test.zip", f, password="password")
print(r)
