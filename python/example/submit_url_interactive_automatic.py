# Copyright (C) 2020 Hatching B.V
# All rights reserved.

import io, time

from triage import Client

url = "https://api.tria.ge"
token = "<YOUR-APIKEY-HERE>"

c = Client(token, root_url=url)
r = c.submit_sample_url("http://google.com", interactive=True)
print(r)
time.sleep(1)
r = c.set_sample_profile_automatically(r["id"], [])
