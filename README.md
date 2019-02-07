# sgk 

```shell
NAME:
   sgk - Quickly manage scratch GCP K8s clusters to test Sourcegraph on

USAGE:
   sgk [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     create, up, start            create a new k8s cluster in GCP
     delete, destroy, down, stop  delete a k8s cluster from GCP
     help, h                      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --name value, -n value  name of the k8s cluster (default: "geoffrey-cluster-test") [$NAME]
   --zone value            zone to create the cluster in (default: "us-central1-a") [$ZONE]
   --project value         GCP project to create the cluster in (default: "sourcegraph-server") [$PROJECT]
   --help, -h              show help
   --version, -v           print the version
```
