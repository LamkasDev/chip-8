package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

func SelectROM() (string, error) {
	paths := make(map[string]string)
	list := []string{}
	filepath.Walk("data", func(path string, file os.FileInfo, err error) error {
		if err == nil && strings.Contains(file.Name(), ".ch8") {
			name := strings.Replace(file.Name(), ".ch8", "", -1)
			paths[name] = path
			list = append(list, name)
		}

		return nil
	})
	prompt := promptui.Select{
		Label:             "Select ROM",
		Items:             list,
		Size:              10,
		StartInSearchMode: true,
		Searcher: func(input string, index int) bool {
			return strings.Contains(strings.ToLower(list[index]), strings.ToLower(input))
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return paths[result], nil
}
