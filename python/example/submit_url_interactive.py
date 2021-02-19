# Copyright (C) 2020 Hatching B.V
# All rights reserved.

import io
import time

from triage import Client

url = "https://api.tria.ge"
token = "<YOUR-APIKEY-HERE>"

c = Client(token, root_url=url)
r = c.submit_sample_url("http://google.com", interactive=True)
print(r)
print("waiting..") # triage takes some time to process
profile = input("profile: ")
time.sleep(1)
r = c.set_sample_profile(r["id"], [{"profile": profile}])
print(r)
print("submitted")
