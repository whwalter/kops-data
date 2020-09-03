package main

import (
	resourceops "k8s.io/kops/pkg/resources/ops"
//	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/pkg/resources"
//	"k8s.io/kops/pkg/resources/aws"
	"k8s.io/kops/upup/pkg/fi/cloudup"
//	"k8s.io/kops/upup/pkg/fi/cloudup/awsup"
//	"k8s.io/kops/pkg/resources/aws"
	"k8s.io/kops/cmd/kops/util"
	"k8s.io/kops/pkg/client/simple"
	kopsapi "k8s.io/kops/pkg/apis/kops"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"os"
	"fmt"
)

func main() {

	clusterName := os.Getenv("KOPS_CLUSTER")
	if clusterName == "" {
		panic(fmt.Errorf("KOPS_CLUSTER environment variable is required: got %s", clusterName))
	}

	stateStore := os.Getenv("KOPS_STATE_STORE")
	if clusterName == "" {
		panic(fmt.Errorf("KOPS_CLUSTER environment variable is required: got %s", clusterName))
	}

	fmt.Printf("Cluster: %s\n", clusterName)
	opts := util.FactoryOptions{ RegistryPath: stateStore}
	factory := util.NewFactory(&opts)

	cluster, err := GetCluster(*factory, clusterName)
	if err != nil {
		panic(err)
	}

	cloud, err := cloudup.BuildCloud(cluster)
	if err != nil {
		panic(err)
	}

	resourcesList, err := resourceops.ListResources(cloud, clusterName, "us-east-1")
	if err != nil {
		panic(err)
	}

	outputs  := map[string]map[string]*resources.Resource{}
	for _, r := range resourcesList {
		if outputs[r.Type] == nil {
			outputs[r.Type] = map[string]*resources.Resource{ r.Name: r }
		} else {
			outputs[r.Type][r.Name] = r
		}
	}
	for k,v := range outputs {
		fmt.Printf("%s\n", k)
		for _, v := range v {
			fmt.Printf("Name: %s\tID: %s\n", v.Name, v.ID)
		}
		fmt.Println()
	}
}

type Factory interface {
	Clientset() (simple.Clientset, error)
}

func GetCluster(factory util.Factory, clusterName string) (*kopsapi.Cluster, error) {
	if clusterName == "" {
		return nil, field.Required(field.NewPath("ClusterName"), "Cluster name is required")
	}

	clientset, err := factory.Clientset()
	if err != nil {
		return nil, err
	}

	cluster, err := clientset.GetCluster(clusterName)
	if err != nil {
		return nil, fmt.Errorf("error reading cluster configuration: %v", err)
	}
	if cluster == nil {
		return nil, fmt.Errorf("cluster %q not found", clusterName)
	}

	if clusterName != cluster.ObjectMeta.Name {
		return nil, fmt.Errorf("cluster name did not match expected name: %v vs %v", clusterName, cluster.ObjectMeta.Name)
	}
	return cluster, nil
}

