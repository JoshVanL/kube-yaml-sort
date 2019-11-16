package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

type options struct {
	outFile string
}

func (o *options) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.outFile, "file", "o",
		"", "file to send to the output rather than stdout.")
}

func main() {
	opts := new(options)

	cmd := &cobra.Command{
		Use:   "kube-yaml-sort",
		Short: "This command takes in Kubernetes YAML objects and outputs the manifests in alphabetical order.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var yaml []byte
			var err error

			if len(args) == 0 {
				yaml, err = readStdin()
				if err != nil {
					return err
				}

			} else {
				yaml, err = readFiles(args)
				if err != nil {
					return err
				}
			}

			sorted, err := SortYAMLObjects(yaml)
			if err != nil {
				return err
			}

			if cmd.Flag("file").Changed {
				err = ioutil.WriteFile(opts.outFile, sorted, 0644)
				if err != nil {
					return err
				}

			} else {
				fmt.Printf("%s", sorted)
			}

			return nil
		},
	}

	opts.addFlags(cmd)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func readStdin() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)

	var out []byte
	for {
		in, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		out = append(out, in)
	}

	return out, nil
}

func readFiles(fs []string) ([]byte, error) {
	if len(fs) == 0 {
		return nil, errors.New("at least one file needed as input")
	}

	out, err := ioutil.ReadFile(fs[0])
	if err != nil {
		return nil, err
	}

	for _, f := range fs[1:] {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}

		out = append(out, yamlsepnl...)
		out = append(out, b...)
	}

	return out, nil
}
