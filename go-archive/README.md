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

## random (un)related code snippets

```
#cloud-config (hcloud) (user data)
chpasswd:
  list: |
    root:chooseAlongAndComplicatedPasswd
    user2:password2
  expire: False
```

```
#cloud-config (user data)
chpasswd:
  list: |
    root:$1$SaltSalt$EZlGtD8tNS8u/ZCr7qjKK.
  expire: False
```

```
$ openssl passwd -1 -salt SaltSalt 123456
$1$SaltSalt$EZlGtD8tNS8u/ZCr7qjKK.
```

```
#cloud-config (user data)

users:
  - default
  - name: root
    lock_passwd: false
    passwd: $1$SaltSalt$EZlGtD8tNS8u/ZCr7qjKK.
  - name: admin
    groups: root
    lock_passwd: false
    passwd: $1$SaltSalt$EZlGtD8tNS8u/ZCr7qjKK.
```

## hetzner
### resources:
- https://gist.github.com/manveru/a0cf4ac0d6ed111b2d41d61c15c90355
- https://github.com/hetznercloud/cli/blob/50a7de3c37155fd71ff1e6cb0f9206eeb6c37eb8/cli/server_create.go


## setup flow ?!

### setup toolsServer
- check if toolsServer is wanted (means is described in cluster.yaml)
  - exit if yes (bec. it's optinal)
- check provider connection + authentication
- gather_facts
- check if toolsServer already exists
  - exit if yes
- provision toolsServer
- install (cli) tools
  - 100%
    - nano
    - vi(m)
    - ssh
    - openssl
    - kubectl
  - optional
    - git
    - helm

###
gather_facts
test_if_proxing_kubectl_cmd_is_possible / toolsNode is available
  if fail
    try to determine new node (max 3 tries if amount of nodes to try is available)

### provider setup / provision
ssh_pub_keys
  local/check
    create_if_not_exist
  remote/check
    create_if_not_exist / add
(create_network)
vms
  remote/check
  remote/create_if_not_exist
  remote/trigger_vm_on_if_not_already
  remote/wait_for_vms
(ensure_vm_ssh_access)

### os setup basic and provider hooks
before_os
os (ubuntu expected)
  swap
  docker
  k8s_bins
after_os

## todo's
- config - switch to use pointer for empty fields
- config
  - BaseName for groups master/worker can might be just omitted and we can use the cluster name + role + rnd string
  - staticServer.Name might should be smth like "nameMiddlePart" or "nameSuffix"
- derived config?
- sanity check
- https://cyruslab.net/2020/10/23/golang-how-to-write-ssh-hostkeycallback/

- docker run -it --rm -v $(pwd):/workdir -w /workdir --entrypoint=ash golang:1.17.6-alpine3.15
