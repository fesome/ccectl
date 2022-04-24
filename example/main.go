package main

import (
	"context"
	"flag"
	"fmt"

	v1 "github.com/xiaohuifirst/ccectl/apis/calico/v1"
	calicov1 "github.com/xiaohuifirst/ccectl/client/clientset/versioned/typed/calico/v1"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := calicov1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ipPools, err := clientset.IPPools().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d ipPools in the cluster\n", len(ipPools.Items))

	for _, pool := range ipPools.Items {
		poolBytes, err := yaml.Marshal(pool)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(poolBytes))
	}
	
	poolInfo, err := clientset.IPPools().Create(context.TODO(), &v1.IPPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "pool-test",
			Labels:                     map[string]string{
				"organization": "test",
			},
		},
		Spec:       v1.IPPoolSpec{
			CIDR:             "192.171.0.0/16",
			VXLANMode:        v1.VXLANModeNever,
			IPIPMode:         v1.IPIPModeNever,
			NATOutgoing:      false,
			Disabled:         false,
			DisableBGPExport: false,
			BlockSize:        26,
			NodeSelector:     "cop!=monitor",
		},
	}, metav1.CreateOptions{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(poolInfo)

	poolBytes, err := yaml.Marshal(poolInfo)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(poolBytes))
}
