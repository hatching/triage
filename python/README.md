
# Usage


## triage.Client

### \_\_init__
`token`: String

`root_url:` String
* default: https://api.tria.ge


### submit_sample_file
`filename`: String

`file`: File

`interactive`: Boolean
* default: False

`profiles`: String[]
* default: []

##### Return JSON
https://tria.ge/docs/cloud-api/samples/#submitting-a-file


### submit_sample_url
`url`: String

`interactive`: Boolean
* default: False

`profiles`: String[]
* default: []

##### Return JSON
https://tria.ge/docs/cloud-api/samples/#submitting-a-url

### set_sample_profile
`sample_id`: String

`profiles`: String[]


##### Return JSON
{}

### set_sample_profile_automatically
`sample_id`: String

`pick`: String[]
* default: []

##### Return JSON
{}

### owned_samples

`max`: Int
* default: 20

##### Return Paginator
https://tria.ge/docs/cloud-api/samples/#get-samples

### public_samples

`max`: Int
* default: 20

##### Return Paginator
https://tria.ge/docs/cloud-api/samples/#get-samples

### search_samples

`query`: String
`max`: Int
* default: 20

##### Return Paginator
https://tria.ge/docs/cloud-api/samples/#get-samples

### sample_by_id

`sample_id`: String

##### Return JSON
https://tria.ge/docs/cloud-api/samples/#get-samplessampleid

### delete_sample

`sample_id`: String

##### Return JSON
{}

### static_report

`sample_id`: String

##### Return JSON
https://tria.ge/docs/cloud-api/samples/#get-samplessampleidreportsstatic

### overview_report

`sample_id`: String

##### Return JSON
https://tria.ge/docs/cloud-api/samples/#get-samplessampleidoverviewjson

### task_report

`sample_id`: String

`task_id`: String

##### Return JSON
https://tria.ge/docs/cloud-api/samples/#get-samplessampleidtaskidreport_triagejson


### sample_task_file

`sample_id`: String

`task_id`: String

`filename`: String

##### Return binary data
file contents

### sample_archive_tar

`sample_id`: String

##### Return binary data
file contents

### sample_archive_zip

`sample_id`: String

##### Return binary data
file contents

### create_profile

`name`: String

`tags`: String[]

`network`: String[]
* Either \"internet\", \"drop\" or None

`timeout`: Int

##### Return JSON
https://tria.ge/docs/cloud-api/profiles/

### delete_profile

`profile_id`: String

##### Return JSON
{}

### profiles

##### Return Paginator
https://tria.ge/docs/cloud-api/profiles/


### sample_events

`sample_id`: String

##### Return JSON
https://tria.ge/docs/cloud-api/samples/#get-samplesevents
