## Benefits
`esops` takes the minimalist apporach to manage  the elasticsearch cluster.
* It relies on the ElasticSearch REST API only, i.e. reduces upgrade overhead as it does not use any elasticsearch go client api. 
* Works for 1.5 to 6.x API unless the API is changed or deprecated, no need to manage different versions. 


```
$./esops
NAME:
   esops - A new cli application

USAGE:
   esops [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     repo, r       managed s3_repo
     snapshot, ss  snapshot/list/status
     restore, r    restore/status
     cluster, c    disable_shard/enable_shard/setting/shutdown/
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --server value, -s value  elasticsearch server (default: "localhost")
   --help, -h                show help
   --version, -v             print the version
```