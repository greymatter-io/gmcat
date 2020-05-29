package catalog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/r3labs/diff"
	"local.only/gmcat/internal/utils"
)

type Cluster struct {
	ClusterName             string     `json:"clusterName"`
	ZoneName                string     `json:"zoneName"`
	Name                    string     `json:"name"`
	Version                 string     `json:"version"`
	Owner                   string     `json:"owner"`
	Capability              string     `json:"capability"`
	Runtime                 string     `json:"runtime"`
	Documentation           string     `json:"documentation"`
	PrometheusJob           string     `json:"prometheusJob"`
	MinInstances            int        `json:"minInstances"`
	MaxInstances            int        `json:"maxInstances"`
	Authorized              bool       `json:"authorized"`
	EnableInstanceMetrics   bool       `json:"enableInstanceMetrics"`
	EnableHistoricalMetrics bool       `json:"enableHistoricalMetrics"`
	MetricsPort             uint32     `json:"metricsPort"`
	Instances               []Instance `json:"instances"`
}

type Instance struct {
	Name      string `json:"name"`
	StartTime uint64 `json:"startTime"`
}

func NewEntry(clusterName, zoneName, displayName, version, owner, capability, runtime, documentation, prometheusJob string, minInstances, maxInstances int, metricsPort uint32, authorized, enableHistoricalMetrics, enableInstanceMetrics bool, instances []Instance) (c *Cluster) {
	c = &Cluster{
		ClusterName:             clusterName,
		ZoneName:                zoneName,
		Name:                    displayName,
		Version:                 version,
		Owner:                   owner,
		Capability:              capability,
		Runtime:                 runtime,
		Documentation:           documentation,
		PrometheusJob:           prometheusJob,
		MinInstances:            minInstances,
		MaxInstances:            maxInstances,
		Authorized:              authorized,
		EnableInstanceMetrics:   enableInstanceMetrics,
		EnableHistoricalMetrics: enableHistoricalMetrics,
		MetricsPort:             metricsPort,
		Instances:               instances,
	}
	return
}

func NewInstance(name string, st uint64) (i *Instance) {
	i = &Instance{
		Name:      name,
		StartTime: st,
	}
	return
}

type CatalogEntryMap map[string]*Cluster

func (c *Cluster) PrintCluster() {
	prettyJSON, err := json.MarshalIndent(c, "", "")
	utils.Check(err)
	fmt.Printf("%s\n", string(prettyJSON))
}

//takes input file path and will create a catalog Cluster struct from that
func FileToEntryCluster(filePath string) (c *Cluster, err error) {
	filePath = strings.TrimSpace(filePath)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading in file %s\n", err.Error())
		return
	}

	err = json.Unmarshal(data, &c)
	utils.Check(err)
	return
}

//Write cluster to file
func (c *Cluster) ClusterToFile(path, name string) (absPathName string, err error) {
	absPathName = fmt.Sprint(path, "/", name)

	err = os.MkdirAll(path, 0755)
	utils.Check(err)

	file, err := json.MarshalIndent(c, "", "")
	utils.Check(err)

	err = ioutil.WriteFile(absPathName, file, 0644)
	utils.Check(err)

	return
}

func DiffCatalogCluster(from, to *Cluster) {
	//print old and new entries
	fmt.Println("------------- OLD Catalog Entry ------------- ")
	from.PrintCluster()
	fmt.Println("------------- NEW Catalog Entry ------------- ")
	to.PrintCluster()
	fmt.Println("--------------------------------------------- ")

	//diff from local entry to in-mesh entry
	changelog, _ := diff.Diff(from, to)
	fmt.Printf("The change from local to in-mesh is:\n%s\n\n", changelog)
}
