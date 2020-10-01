# Installation

`go version` Make sure a recent version of go is installed

`cd go/` Change directory to this projects go folder.

`make` This will make and output the binary in your GOPATH

`triage` To use the Triage cli.


# Usage
## authenticate
`-t`: string; API Token

API token will be stored on disk.

## submit
`-f`: file path

`-u`: url

##### return
Sample submitted

ID: {id}

Status: {status}

-- target info --

## select-profile
`-s`: string; ID of a sample.

##### return
Interactive flow to select profile for submitted file

note: does not seem to work very cleanly

## list, ls
`-n`: int; Numbers of samples to return.

`-public`: bool; Query public set

##### return
List of https://tria.ge/docs/cloud-api/samples/#get-samplessampleid

## file
`-s`: string; Sample ID

`-t`: string; Task ID

`-f`: string; Filename

`-o`: string; Output file

##### return
Downloaded file

note: Not able to make it work

## archive
`-s`: string; Sample ID

`-f`: string; Archive format

`-o`: string; Output file

##### return
Download Sample as an archive


## delete, del
`-s`: string; Sample ID

##### return
Delete a Sample

## report
`-s`: string; Sample ID

`-static`: boolean; Query statis report

`-t`: string; Task ID

##### return
Get the report of a single sample

https://tria.ge/docs/cloud-api/samples/#get-samplessampleid

## create-profile
`-name`: string; Name of new profile

`-tags`: string; Comma seperated set of tags

`-network`: string; Use *network*, *drop* or *unset*

`-timeout`: string (4m0s); Timeout of profile

##### return
Profile created

  ID:   {id}

  Name: {name}


## delete-profile
`-p`: string; Name or ID of profile

###### return
Nothing returned, profile deleted

## list-profiles

##### return
List of profiles

https://tria.ge/docs/cloud-api/profiles/#the-profile-object
