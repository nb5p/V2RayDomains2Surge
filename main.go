package main

import (
	"fmt"
	"os"
	"path"

	"v2raydomains2surge/rule"
)

func main() {
	if len(os.Args) < 3 {
		println("Usage: <v2ray-domains-path> <output-path>")

		os.Exit(1)
	}

	data := path.Join(os.Args[1], "data")
	generated := os.Args[2]

	_ = os.MkdirAll(generated, 0755)

	ruleSets, err := rule.ParseDirectory(data)
	if err != nil {
		println("Load domains: " + err.Error())

		os.Exit(1)
	}

	for name := range ruleSets {
		tags, err := rule.Resolve(ruleSets, name)
		if err != nil {
			println("Resolve " + name + ": " + err.Error())

			continue
		}

		for tag, rules := range tags {
			var outputPath string

			if tag == "" {
				outputPath = path.Join(generated, fmt.Sprintf("%s.txt", name))
			} else {
				outputPath = path.Join(generated, fmt.Sprintf("%s@%s.txt", name, tag))
			}

			file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				println("Write file " + outputPath + ": " + err.Error())

				continue
			}

			_, _ = file.WriteString(fmt.Sprintf("# Generated from https://github.com/v2fly/domain-list-community/tree/master/data/%s\n\n", name))
			_, _ = file.WriteString("# Behavior: domain-set\n\n")

			for _, domain := range rules {
				_, _ = file.WriteString(fmt.Sprintf("%s\n", domain))
			}

			_ = file.Close()
		}
	}
}
