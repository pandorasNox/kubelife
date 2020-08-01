# kubelife
Kubelife is a CLI tool (and kubeadm wrapper) for setting up a single node (and therefore not production ready) kubernetes seed clusters. 

Currently it is intended to be used with cloud vm providers such as digitalocean and hetzner (cloud).

## how it operates
Kubelife expects a `cluster.yaml` file, which describes the infrastructure (servers/vm's) and some meta information.
See (here)[] for a full example of a `cluster.yaml` file.

Kubelife expects a working ssh connection via unix ssh agent to all servers described in the `cluster.yaml`

When everything is configuted kubelife after executing will:
- run commands on all nodes via ssh
- ensures docker installation
- ensures kubeadm, kubelet, kubectl installation
- ensures further needed kubeadm componends
- installs some basic cli tools
- setup k8s master node on the in the `cluster.yaml` described `primary_master` node
- waits for `primary_master` setup
- joins worker nodes (as described in the `cluster.yaml`)

## future
- currently the tool is for my private setups
- in the future this might evolve into a full fledged k8s lifecycle management tool for seed clusters
