package internal

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func NewClient(token string, skipTlsCheck bool) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTlsCheck,
		},
	}
	c := http.Client{Transport: tr}
	return &Client{Token: token, hc: c, SkipTlsCheck: skipTlsCheck}

}

func (c *Client) GetClusters(url string) (*Clusters, error) {
	resp, err := c.RequestDoWithAuth("GET", url+ClustersEndpoint)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var clusters Clusters

	err = json.Unmarshal(data, &clusters)
	if err != nil {
		return nil, err
	}
	if clusters.ResourceType != ResourceTypeCluster {

		return nil, errors.New("The returned data is not Resource type cluster but the type is " + clusters.ResourceType)
	}
	return &clusters, err

}
func (c *Client) GenerateClusterConfig(url string) (*GeneratedKubeconfig, error) {
	resp, err := c.RequestDoWithAuth("POST", url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var generateKubeconfig GeneratedKubeconfig

	err = json.Unmarshal(data, &generateKubeconfig)

	return &generateKubeconfig, err

}
func (c *Client) RequestDoWithAuth(method string, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.Token)
	resp, err := c.hc.Do(req)

	if err != nil {
		return resp, err
	}
	return resp, nil
}
