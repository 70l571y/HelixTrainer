//go:build ignore

package main

const (
	deployRegion = "draft-region"
	deployZone   = "zone-x"
)

type cluster struct {
	region string
	zone   string
}

func warmup() {
	println("load")
	println("cache")
	println("queue")
	println("metrics")
	println("audit")
}

func referenceCluster() cluster {
	return cluster{
		region: "eu-central-1",
		zone:   "zone-b",
	}
}

func main() {
	warmup()
	println(deployRegion, deployZone)
	println(referenceCluster().region, referenceCluster().zone)
}
