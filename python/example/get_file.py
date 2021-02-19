# Copyright (C) 2020 Hatching B.V
# All rights reserved.

from triage import Client

url = "https://api.tria.ge"
token = "<YOUR-APIKEY-HERE>"

c = Client(token, root_url=url)
sample_id = input("sample id: ") # e.g. 200916-1tmctk8x46
t = c.sample_by_id(sample_id)
task = input("task id: ") # e.g. behavioral1
t = c.task_report(sample_id, task)
file = input("file: ") # e.g. memory/1652-0-0x0000000000000000-mapping.dmp
t = c.sample_task_file(sample_id, task, file)
with open("test.dmp", "wb") as f:
    f.write(t)
print("wrote test.dmp")

