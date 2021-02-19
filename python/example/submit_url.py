# Copyright (C) 2021 Hatching B.V
# All rights reserved.

import io

from triage import Client

url = "https://api.tria.ge"
token = "<YOUR-APIKEY-HERE>"

c = Client(token, root_url=url)
r = c.submit_sample_url("http://google.com")
print(r)
