# pdbval

A CLI to validate PodDisruptionBudget against Deployment.

## Usage

```console
$ cat pdb.yaml
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: test-pdb
  namespace: test-pdb
spec:
  minAvailable: 70%
  # maxUnavailable: 1%
  selector:
    matchLabels:
      app: nginx

$ cat deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  namespace: test-pdb
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2

$ pdbval -pdb pdb.yaml -deploy manifests/test-pdb/deployment.yaml
PodDisruptionBudget(test-pdb): minAvailable(70%) is greater than or equal to Deployment(nginx) replicas(2)
```
