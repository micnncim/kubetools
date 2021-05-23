# restart

A minimal implementation of `kubectl rollout restart`.
`kubectl` is more recommended.
Only supports Deployments.

## Usage

```console
$ restart -n $NAMESPACE $DEPLOYMENT
$ kubectl get po
NAME                     READY   STATUS              RESTARTS   AGE
nginx-5566b966d7-jgffm   1/1     Running             0          10s
nginx-7f47d6fdc9-6cxv5   0/1     ContainerCreating   0          1s

```
