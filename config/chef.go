package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vadv/chef"
)

// load hosts from chef

func (config *Config) fetchHostsFromChef(q string) error {

	if config.chefKey == `` {
		data, err := ioutil.ReadFile(config.defaultchefKeyPath)
		if err == nil {
			config.chefKey = string(data)
		} else {
			fmt.Fprintf(os.Stderr, "Chef key file access error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	client, err := chef.NewClient(&chef.Config{
		SkipSSL: true,
		Name:    config.chefClient,
		Key:     config.chefKey,
		BaseURL: config.chefUrl,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Start chef `%s` search query: `%s`\n", config.chefUrl, q)

	part := make(map[string]interface{})
	part["attr"] = []string{config.chefAttr}

	res, err := client.Search.PartialExec("node", q, part)
	if err != nil {
		return err
	}

	for _, row := range res.Rows {
		// row = {"url": "http://chef/node", "data": {"attr": "<response>"}}
		line, ok := row.(map[string]interface{})
		if !ok {
			fmt.Fprintf(os.Stderr, "Bad chef response #1: %#v\n", line)
			os.Exit(1)
		}
		dataRaw, found := line["data"]
		if !found {
			fmt.Fprintf(os.Stderr, "Bad chef response #2: %#v\n", line)
			os.Exit(1)
		}
		data, transform := dataRaw.(map[string]interface{})
		if !transform {
			fmt.Fprintf(os.Stderr, "Bad chef response #3: %#v\n", dataRaw)
			os.Exit(1)
		}
		host, completed := data["attr"]
		if !completed {
			fmt.Fprintf(os.Stderr, "Empty attribute from chef: %#v\n", data)
			os.Exit(1)
		}
		config.hosts = append(config.hosts, fmt.Sprintf("%v", host))
	}

	fmt.Fprintf(os.Stderr, "Chef search return %d items\n", len(config.hosts))
	if len(config.hosts) == 0 {
		os.Exit(1)
	}

	return nil
}
