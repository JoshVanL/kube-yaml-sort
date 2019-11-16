# kube-yaml-sort

This is a small CLI tool to sort Kubernetes manifests. The order of sorting is:
- API Version
- Kind
- Namespace
- Name

The tool accepts multiple files as input arguments or Stdin with no arguments.

The tool will output to Stdout by default or to a file (--file -o).

This tool should work for all Kubernetes manifests that have objects that have
any of the API version, kind or metadata stanzas at the top level.

To build this tool simply run `go build` or `go install`.

Here is an example of the tool:

```
$ cat in.yaml
apiVersion: api-B
kind: kind-A
metadata:
  name: name-A
  namespace: name-B
---
apiVersion: api-A
kind: kind-A
metadata:
  name: name-B
  namespace: name-A
---
apiVersion: api-B
kind: kind-A
metadata:
  name: name-B
  namespace: name-A
---
apiVersion: api-A
kind: kind-B
metadata:
  name: name-A
  namespace: name-A

$ ./kube-yaml-sort in.yaml -o out.yaml

$ cat out.yaml
apiVersion: api-A
kind: kind-A
metadata:
  name: name-B
  namespace: name-A
---
apiVersion: api-A
kind: kind-B
metadata:
  name: name-A
  namespace: name-A
---
apiVersion: api-B
kind: kind-A
metadata:
  name: name-B
  namespace: name-A
---
apiVersion: api-B
kind: kind-A
metadata:
  name: name-A
  namespace: name-B
```
