package mesh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"local.only/gmcat/internal/catalog"
	"local.only/gmcat/internal/utils"
)

func (h *Handler) UseCatalog(method, endpoint string, body interface{}) (resp *http.Response) {
	// var buff *bytes.Buffer = nil
	var buff = bytes.NewBuffer([]byte(nil))
	buff.Reset()
	if body != nil {
		data, err := json.Marshal(body)
		utils.Check(err)

		buff = bytes.NewBuffer(data)
	}

	url := h.config.Edge + h.config.Catalog + endpoint
	fmt.Printf("\nUsing:\n  URL: %s\n  Buff: %s\n", url, buff)

	req, err := http.NewRequest(method, url, buff)
	utils.Check(err)

	resp, err = h.client.Do(req)
	utils.Check(err)
	return
}

// parameters are amended to /clusters
// to search all clusters use ""
// to find specific cluster use "/<catalogName>"
func (h *Handler) ListCatalogEntries() (results []string) {
	resp := h.UseCatalog("GET", "/clusters", nil)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)

	var clusters []catalog.Cluster
	err = json.Unmarshal(data, &clusters)
	utils.Check(err)

	for i, _ := range clusters {
		results = append(results, clusters[i].ClusterName)
	}
	return
}

//Get a catalog entry
func (h *Handler) GetCatalogEntry(clusterName string) (c *catalog.Cluster) {
	endpoint := fmt.Sprintf("/clusters/%s", clusterName)
	resp := h.UseCatalog("GET", endpoint, nil)
	defer resp.Body.Close()

	data2, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)

	var catalog_entry []catalog.Cluster
	err = json.Unmarshal(data2, &catalog_entry)
	if err != nil {
		fmt.Printf("GetCatalogEntry issue: \n %s", err.Error())
	}

	c = &catalog_entry[0]

	catalog_entry[0].PrintCluster()

	return
}

//Create a catalog entry
func (h *Handler) CreateCatalogEntry(cluster *catalog.Cluster) (success bool) {
	fmt.Println("Starting Create Entry")
	success = false
	if len(cluster.ClusterName) > 0 {
		endpoint := "/clusters"
		resp := h.UseCatalog("POST", endpoint, cluster)
		if resp.Status == "200 OK" {
			success = true
			fmt.Printf("Created %s entry from catalog\n\n", cluster.ClusterName)
		} else {
			fmt.Printf("There was an issue creating %s entry from catalog\n\n", cluster.ClusterName)
		}
		defer resp.Body.Close()
	}
	return
}

//Delete a catalog entry
func (h *Handler) DeleteCatalogEntry(clusterName, zone string) (success bool) {
	fmt.Println("Starting Delete Entry")
	success = false
	if len(clusterName) > 0 && len(zone) > 0 {
		endpoint := "/clusters/" + clusterName + "?zoneName=" + zone
		resp := h.UseCatalog("DELETE", endpoint, nil)
		if resp.Status == "200 OK" {
			success = true
			fmt.Printf("Deleted %s entry from catalog\n\n", clusterName)
		} else {
			fmt.Printf("There was an issue deleting %s entry from catalog\n\n", clusterName)
		}
		defer resp.Body.Close()
	}
	return
}
