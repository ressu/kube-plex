# kube-plex

kube-plex is a scalable Plex Media Server solution for Kubernetes. It
distributes transcode jobs by creating jobs in a Kubernetes cluster to perform
transcodes, instead of running transcodes on the Plex Media Server instance
itself.

## How it works

kube-plex works by replacing the Plex Transcoder program on the main PMS
instance with our own little shim. This shim intercepts calls to Plex
Transcoder, and creates Kubernetes pods to perform the work instead. These
pods use shared persistent volumes to store the results of the transcode (and
read your media!).

## Prerequisites

* A persistent volume type that supports ReadWriteMany volumes (e.g. NFS,
Amazon EFS)
* Your Plex Media Server *must* be configured to allow connections from
unauthorized users for your pod network, else the transcode job is unable to
report information back to Plex about the state of the transcode job. At some
point in the future this may change, but it is a required step in order to make
transcodes work right now.

## Setup

This guide will go through setting up a Plex Media Server instance on a
Kubernetes cluster, configured to launch transcode jobs on the same cluster
in pods created in the same 'plex' namespace.

1) Obtain a Plex Claim Token by visiting [plex.tv/claim](https://plex.tv/claim).
This will be used to bind your new PMS instance to your own user account
automatically.

2) Deploy the Helm chart included in this repository using the claim token
obtained in step 1.

If you have pre-existing persistent volume claims for your
media, you can specify its name with `--set persistence.data.claimName`. If not
specified, a persistent volume will be automatically provisioned for you.

In order for the transcoding to work, a shared transcode persistent volume claim needs to be defined with `--set persistence.transcode.claimName` or by defining the relevant parameters separately.

```bash
➜  helm install plex ./charts/kube-plex \
    --namespace plex \
    --set claimToken=[insert claim token here] \
    --set persistence.data.claimName=existing-pms-data-pvc \
    --set persistence.transcode.claimName=shared-pms-transcode-pvc \
    --set ingress.enabled=true
```

This will deploy a scalable Plex Media Server instance that uses Kubernetes as
a backend for executing transcode jobs.

3) Access the Plex dashboard. If you used claim token above, the plex instance
should be visible in [Plex Web App](https://app.plex.tv). If the token
registration failed, access the instance using using `kubectl port-forward`.

The instance can be accessed via the load balancer IP (via `kubectl get
service`) or the ingress (if provisioned with `--set ingress.enabled=true`) once
the registration has been completed.

4) Visit Settings->Server->Network and add your pod network subnet to the
`List of IP addresses and networks that are allowed without auth` (near the
bottom). For example, `10.100.0.0/16` is the subnet that pods in my cluster are
assigned IPs from, so I enter `10.100.0.0/16` in the box.

You should now be able to play media from your PMS instance

## Internal operations

Kube-plex will automatically create transcoding jobs within the Kubernetes
instance. The jobs have shared transcode and data mounts with the main kube-plex
pod. Kube-plex replaces the `Plex Transcoder` binary with a launcher on Plex
startup. Kube-plex launcher processes the arguments from Plex and creates a
transcoding job to handle the final transcoding.

```bash
ubuntu@wanted-wolf:~$ kubectl get pod,job
NAME                                        READY   STATUS    RESTARTS   AGE
pod/kube-plex-694d659b64-7wg2b   1/1     Running   0          6d23h
pod/pms-elastic-transcoder-tqw5s-8w2bc      1/1     Running   0          4s

NAME                                     COMPLETIONS   DURATION   AGE
job.batch/pms-elastic-transcoder-tqw5s   0/1           4s         5s
```

Transcoder pod will run a shim which will
* Download codecs from main kube-plex pod
* Relay transcoder callbacks from `Plex Transcoder` to main kube-plex

Logging from kube-plex processes is written to Plex process and can be viewed in `Settings->Manage->Console`.