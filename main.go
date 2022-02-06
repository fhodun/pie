package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const ignoreFile = "pieignore"

type Project struct {
	DirName string
	Path    string
}

func main() {
	var rootCmd = &cobra.Command{
		Use:     "pie",
		Short:   "Refactor definitely not yours C# code",
		Example: "pie [options] <folder>",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectPath := strings.TrimSuffix(args[0], "/") // remove trailing slash

			// check if the folder exists
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				log.Panic("Path does not exist", err)
			}
			log.Println("Starting pie...")

			// parse path to get the project name
			projPathArr := strings.Split(projectPath, "/")
			projectDir := projPathArr[len(projPathArr)-1]
			project := Project{
				DirName: projectDir,
				Path:    projectPath,
			}

			readFile, err := os.Open(ignoreFile)
			if err != nil {
				log.Panic("Could not open "+ignoreFile, err)
			}
			fileScanner := bufio.NewScanner(readFile)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				ignore := fileScanner.Text()

				// check line is commented or empty
				if strings.HasPrefix(ignore, "#") || ignore == "" {
					continue
				}

				// parse line if it contains a wildcard
				ut, err := template.New("ignore").Parse(ignore)
				if err != nil {
					log.Panic(err)
				}
				var ignoreBuffer bytes.Buffer
				if err := ut.Execute(&ignoreBuffer, project); err != nil {
					log.Panic(err)
				}
				ignore = ignoreBuffer.String()

				ignoredPath := fmt.Sprintf("%s/%s", project.Path, ignore)
				// check if the path exists
				if _, err := os.Stat(ignoredPath); !os.IsNotExist(err) {
					log.Printf("Removing useless %s thing", ignore)
					// remove path
					if err := os.RemoveAll(ignoredPath); err != nil {
						log.Panic(err)
					}
				}
			}
			readFile.Close()
		},
	}
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
