# PoC kops Cluster Data Sourcing

This is a PoC for using kops resource and utils libraries to source information about the kops cluster. 

The tool sources any cloud profiles needed in a similar way to kops. It takes a `KOPS_CLUSTER` environment variable that names the cluster to source information for. 
An example run:

```
go build -o kops-data main.go
KOPS_CLUSTER=test.k8s.local KOPS_STATE_STORE=s3://abucketwithaconfig ./kops-data
```

## Notes

I intentionally left the commented out kops imports as notes of other options. This implementation is cloud provider agnostic. If we prefer an aws cluster specific tool, the aws resources package can provide more tailored ListX functions. 

I committed the vendor directory purposfully and you should not run go mod vendor while using this PoC. 

Kops modules have two dependencies that would make maintenance of this type of code sligtly more burdonsome. 

1. Kops depend on k8s modules that are not intened to be libraries. The packages are published with a v0.0.0 release that doesn't result in a real dependency. This requires replace directives in the go.mod file to keep the version of these libraries in step with kops and the targeted cluster version. 

2. Kops modules use generated bindata.go files that are consumed as packages but need to be generated by the version of kops being used. see vendor/k8s.io/kops/upup/models/bindata.go for an example. This was generated from the kops v1.15.3 source using the cloudup and nodeup direcories that exist in the source but not the vendor versions of kops. 
