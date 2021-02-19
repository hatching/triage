# Copyright (C) 2020 Hatching B.V
# All rights reserved.

import io

from triage import Client

url = "https://api.tria.ge"
token = "<YOUR-APIKEY-HERE>"

c = Client(token, root_url=url)
f = io.StringIO("some initial text data")
r = c.submit_sample_file("test.exe", f)
print(r)
