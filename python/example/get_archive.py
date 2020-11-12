# Copyright (C) 2020 Hatching B.V
# All rights reserved.

from triage import Client

url = "https://api.tria.ge"
token = "<YOUR-APIKEY-HERE>"

c = Client(token, root_url=url)
sample_id = input("sample id: ")
t = c.sample_archive_tar(sample_id)
with open("test.tar", "wb") as f:
    f.write(t)
print("wrote test.tar")

