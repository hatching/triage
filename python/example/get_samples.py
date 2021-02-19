# Copyright (C) 2020 Hatching B.V
# All rights reserved.

from triage import Client

url = "https://api.tria.ge"
token = "<YOUR-APIKEY-HERE>"

c = Client(token, root_url=url)
for i in c.public_samples():
    print(i)

