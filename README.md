# Kubernetes v0.1 (First Commit)

## Overview

✨ This repository was created for historical and archival purposes. It provides static binary builds in addition to the initial Kubernetes code commit made by Joe Beda ([@jbeda](https://github.com/jbeda)). The earliest binaries officially released by Kubernetes start from v0.4.0; this project extends that archive back to v0.1.0.

Also provided is a convenient tutorial for trying this out yourselves! 🚀

## Available Binaries

You can download binaries built for the following platforms:
- Linux: amd64, arm, arm64
- Darwin: amd64, arm64 (Note: Darwin binaries are provided for experimental purposes only)

All binaries are available in the [bin](bin) folder.

## Build Process

The following steps were taken to build these binaries:

1. **Repository Cloning and Checkout**: Cloned the original repository and performed a hard checkout to the first commit
2. **Code Reference Update**: Updated YAML references in the code from `gopkg.in/v1/yaml` to `gopkg.in/yaml.v1`
3. **Dependency Management**: Added a `replace` entry in `go.mod` to resolve circular dependencies by redirecting `github.com/GoogleCloudPlatform/kubernetes` to the current directory
4. **Additional Dependencies**: Executed `go get github.com/fsouza/go-dockerclient@3b6f84ca70bed7ba0353656a48390f840bc55a0e` to lock a specific commit from the same time period

## Tweaks

1. **Enabled Skip Pull**: Modified pkg/kubelet/kubelet.go with a simple check for /etc/kubernetes_skip_pull - Docker pulls will fail from this era owing to legacy registry standards, we will need to preload images, creating this file tells the kubelet to skip the pull process

## Archives

Some of the files in the [examples/guestbook](examples/guestbook) folder makes reference to two container images by Brendan Burns ([@brendandburns](https://github.com/brendandburns)) -

```
REPOSITORY                 TAG       IMAGE ID       CREATED       SIZE
brendanburns/php-redis     latest    2ba4a84f8382   9 years ago   377MB
brendanburns/redis-slave   latest    3baf645dfdde   9 years ago   433MB
```

At the time of writing, these are available on Docker Hub. For archival purposes, these have also been saved and uploaded to the [Internet Archive](https://archive.org/download/brendanburns_php-redis.tar). Should you wish to experiment with these in the future and they're no longer available on Docker Hub, you should hopefully be able to download these and manually load them.

## Scripts

- **Environment Setup**: [01_env_setup.sh](01_env_setup.sh) (for reference only, not necessary if `go.mod` and `go.sum` are used)
- **Build Script**: [02_build.sh](02_build.sh) (customizable to meet different build requirements)

## Trying it out yourself!

Easiest way to explore this, is through [Lima VM](https://github.com/lima-vm/lima). A Virtual Machine image was created for the [Kubernetes v1.0.0 - 10 year project](https://github.com/spurin/kubernetes-v1.0-lab) and this can be re-used to explore this first commit.

Start the instance with  -

```
limactl start --tty=false --cpus 2 --memory 8 https://raw.githubusercontent.com/spurin/kubernetes-v1.0-lab/main/lima.yaml
```

Example usage -

```
% limactl start --tty=false --cpus 2 --memory 8 https://raw.githubusercontent.com/spurin/kubernetes-v1.0-lab/main/lima.yaml
INFO[0000] Terminal is not available, proceeding without opening an editor
INFO[0000] Starting the instance "lima" with VM driver "qemu"
INFO[0000] Attempting to download the image              arch=x86_64 digest="sha256:a348500d04de3352af3944e9dae99a08d60091221e1370941b93bd7a27971568" location="http://cloud-images-archive.ubuntu.com/releases/vivid/release-20160203/ubuntu-15.04-server-cloudimg-amd64-disk1.img"
INFO[0000] Using cache "/Users/james/Library/Caches/lima/download/by-url-sha256/77c902f99f70b354e56b996968a3838834ff4616a650c2a3b45490ba26927ac8/data"
INFO[0000] [hostagent] hostagent socket created at /Users/james/.lima/lima/ha.sock
INFO[0000] [hostagent] Starting QEMU (hint: to watch the boot progress, see "/Users/james/.lima/lima/serial*.log")
INFO[0000] SSH Local Port: 51278
INFO[0000] [hostagent] Waiting for the essential requirement 1 of 2: "ssh"
INFO[0028] [hostagent] Waiting for the essential requirement 1 of 2: "ssh"
INFO[0038] [hostagent] Waiting for the essential requirement 1 of 2: "ssh"
INFO[0051] [hostagent] The essential requirement 1 of 2 is satisfied
INFO[0051] [hostagent] Waiting for the essential requirement 2 of 2: "user session is ready for ssh"
INFO[0051] [hostagent] The essential requirement 2 of 2 is satisfied
INFO[0051] [hostagent] Waiting for the guest agent to be running
INFO[0051] [hostagent] Forwarding "/run/lima-guestagent.sock" (guest) to "/Users/james/.lima/lima/ga.sock" (host)
INFO[0051] [hostagent] Guest agent is running
INFO[0051] [hostagent] Waiting for the final requirement 1 of 1: "boot scripts must have finished"
INFO[0051] [hostagent] Not forwarding TCP 0.0.0.0:22
INFO[0051] [hostagent] Not forwarding TCP [::]:22
INFO[0091] [hostagent] Waiting for the final requirement 1 of 1: "boot scripts must have finished"
INFO[0097] [hostagent] The final requirement 1 of 1 is satisfied
INFO[0098] READY. Run `limactl shell lima` to open the shell.

% limactl shell lima
james@lima-lima:~$
```
 
Once you're in the instance, we'll start by installing docker using the ubuntu pkg version available at the time to match the release period around this time -

```bash
cd; sudo apt install -y docker.io
```

Download and install etcd, version 2, to also match the release period of the time, install to /usr/local/bin -

```bash
curl -L https://github.com/coreos/etcd/releases/download/v2.0.12/etcd-v2.0.12-linux-amd64.tar.gz -o etcd-v2.0.12-linux-amd64.tar.gz
tar xzvf etcd-v2.0.12-linux-amd64.tar.gz
sudo install etcd-v2.0.12-linux-amd64/etcd /usr/local/bin
```

Run etcd on port 4001 (legacy, current is 2379) in the background as root and follow the logs, press `Ctrl-C` when you're ready, this will continue to run in background -

```bash
sudo bash -c 'etcd --listen-client-urls http://0.0.0.0:4001 --advertise-client-urls http://localhost:4001 &> /var/log/etcd.log &'; tail -f /var/log/etcd.log
```

etcdctl is handy for this tutorial in understanding what Kubernetes v0.1.0 is doing behind the scenes, install this -

```bash
ETCD_VERSION=${ETCD_VERSION:-v3.3.1}
curl -L https://github.com/coreos/etcd/releases/download/$ETCD_VERSION/etcd-$ETCD_VERSION-linux-amd64.tar.gz -o etcd-$ETCD_VERSION-linux-amd64.tar.gz
tar xzvf etcd-$ETCD_VERSION-linux-amd64.tar.gz
rm etcd-$ETCD_VERSION-linux-amd64.tar.gz
sudo cp etcd-$ETCD_VERSION-linux-amd64/etcdctl /usr/local/bin/
rm -rf etcd-$ETCD_VERSION-linux-amd64
```

If we check etcdctl, currently we have no keys -

```bash
etcdctl ls --recursive
```

Clone this project repository and move to the kubernetes-v0.1 folder -

```bash
git clone https://github.com/spurin/kubernetes-v0.1-lab.git
cd kubernetes-v0.1
```

Run the apiserver in the background as root and follow the logs, press `Ctrl-C` when you're ready, this will continue to run in the background -

```bash
sudo bash -c 'bin/linux/amd64/apiserver -machines $HOSTNAME -etcd_servers=http://127.0.0.1:4001 &> /var/log/apiserver.log &'; tail -f /var/log/apiserver.log
```

With this running, you can query the apiserver if you wish, example output as follows -

```bash
curl http://127.0.0.1:8080/
<html><body>Welcome to Kubernetes</body></html>

curl http://127.0.0.1:8080/api/v1beta1/replicationControllers
{}

curl http://127.0.0.1:8080/api/v1beta1/tasks
{
    "items": []
}

curl http://127.0.0.1:8080/api/v1beta1/services
{
    "items": null
}
```

Next run the controller-manager and monitor the logs, at this point it will complain, that the /registry/controllers key in etcd does not exist, press `Ctrl-C` when you're ready, this will continue to run in the background -

```bash
sudo bash -c 'bin/linux/amd64/controller-manager -etcd_servers http://127.0.0.1:4001 -master 127.0.0.1:8080 &> /var/log/controller-manager.log &'; tail -f /var/log/controller-manager.log
```

Run the proxy, similar to the controller-manager, this will be poling /registry/services in etcd and will complain, as this key does not exist, press `Ctrl-C` when you're ready, this will continue to run in the background -

```bash
sudo bash -c 'bin/linux/amd64/proxy -etcd_servers http://localhost:4001 &> /var/log/proxy.log &'; tail -f /var/log/proxy.log
```

We'll now run the kubelet, this will monitor /registry/hosts/$HOSTNAME in etcd where the hostname in question is lima-lima if you're using the recommended instance to follow this tutorial, press `Ctrl-C` when you're ready, this will continue to run in the background -

```bash
sudo bash -c 'bin/linux/amd64/kubelet --etcd_servers=http://127.0.0.1:4001 -address=$HOSTNAME &> /var/log/kubelet.log &'; tail -f /var/log/kubelet.log
```

Pulling container images from Docker Hub may fail due to changes in registry standards. Originally, Docker Hub utilized v1 registry standards, which have been deprecated in favor of the newer v2 standards recognized today. As a result, attempts to pull container images from Kubernetes to Docker Hub using the old standards will not succeed. However, there is a workaround to this issue.

To bypass the v1/v2 compatibility problem, preload the necessary images. This involves saving the images into a tar file and then manually loading them using the docker load command. For a more convenient method, Skopeo can be used to directly save images from Docker Hub into a tar file. Start by downloading and preparing Skopeo for this purpose.

```bash
sudo curl -L https://github.com/lework/skopeo-binary/releases/download/v1.14.4/skopeo-linux-amd64 -o /usr/bin/skopeo && sudo chmod 755 /usr/bin/skopeo; sudo mkdir -p /etc/containers; sudo echo '{ "default": [ { "type": "insecureAcceptAnything" } ] }' | sudo tee /etc/containers/policy.json
```

Create a convenient shell function that uses skopeo and docker load to preload images, you can use this function for any other images that you wish to use -

```bash
skopeo-save-load() { local safe_image_name=$(echo "$1" | tr '/:' '_'); local tar_path="/tmp/${safe_image_name}.tar"; skopeo copy "docker://$1" "docker-archive:${tar_path}:$1"; sudo docker load -i "$tar_path"; rm -f "$tar_path"; }
```

Preload nginx:latest -

```bash
skopeo-save-load nginx:latest
```

Disable docker image pulls via the kubelet (custom code modification) -

```bash
sudo touch /etc/kubernetes_skip_pull
```

We now have all of our core components running and respectfully, each of these are polling etcd in different locations. Next we're going to use cloudcfg which operates in a similar capacity to kubectl in a modern kubernetes cluster. It is essentially a convenient tool for making api requests to the apiserver.

Save the following as a template file. This file was based on the json that was created via the use of cloudcfg but, with one addition, `containers":[{"name": "mynginx"]` - in my initial testing, containers were failing as this value wasn't set.

I spent quite a bit of time looking through the code base to try and understand if/why this wasn't being set. Possibly, a dependency difference or maybe even how this tool was being used at the time (it's the first commit, possibly this didn't work). By adding this entry, it allows us to continue -

```bash
echo '{"id":"mynginx","desiredState":{"replicas":1,"replicasInSet":{"name":"mynginx"},"taskTemplate":{"desiredState":{"manifest":{"version":"","volumes":null,"containers":[{"name": "mynginx", "image":"nginx","ports":[{"hostPort":1234,"containerPort":80}]}]}},"labels":{"name":"mynginx"}}},"labels":{"name":"mynginx"}}' > /tmp/mynginx_controller.json
```

cloudcfg expects a file called ~/.kubernetes_auth to exist, we've not setup authentication so a file structured accordingly will suffice -

```bash
cat <<EOF > ~/.kubernetes_auth
{
  "User": "user",
  "Password": ""
}
EOF
```

Then using cloudcfg, we will pass our custom replicationController json manifest to the /replicationControllers endpoint -

```bash
bin/linux/amd64/cloudcfg -h http://localhost:8080 -c /tmp/mynginx_controller.json create /replicationControllers
```

If you now look with etcdctl, you'll see that this registered a controller, which will in turn, be processed by the controller-manager - individual tasks will be created and these will be put in the monitoring path of the kubelet (/registry/hosts/lima-lima/tasks)

```
etcdctl ls --recursive

/events
/events/mynginx
/events/mynginx/119
/registry
/registry/controllers
/registry/controllers/mynginx
/registry/hosts
/registry/hosts/lima-lima
/registry/hosts/lima-lima/kubelet
/registry/hosts/lima-lima/tasks
/registry/hosts/lima-lima/tasks/12df41e3b26ec4c3
```

And if you now check with docker, you'll see that your container is running! If you wish you can also check the kubelet log file to see what was actioned `cat /var/log/kubelet.log` -

```
sudo docker ps -a

CONTAINER ID        IMAGE               COMMAND                CREATED             STATUS                      PORTS                  NAMES
904594d40b22        nginx:latest        "/docker-entrypoint.   2 minutes ago       Up 2 minutes                0.0.0.0:1234->80/tcp   mynginx--12df41e3b26ec4c3--c5ac7d34
```

Congratulations, you've successfully used the first ever commit of Kubernetes (0.1.0)!

The remaining contents of this README mirror those of the original commit to preserve historical accuracy.

***

# Kubernetes

Kubernetes is an open source reference implementation of container cluster management.

## Getting started on Google Compute Engine

### Prerequisites

1. You need a Google Cloud Platform account with billing enabled.  Visit http://cloud.google.com/console for more details
2. You must have Go installed: [www.golang.org](http://www.golang.org)
3. Ensure that your `gcloud` components are up-to-date by running `gcloud components update`.
4. Get the Kubernetes source:  `git clone https://github.com/GoogleCloudPlatform/kubernetes.git`

### Setup
```
cd kubernetes
./src/scripts/dev-build-and-up.sh
```

### Running a container (simple version)
```
cd kubernetes
./src/scripts/build-go.sh
./src/scripts/cloudcfg.sh -p 8080:80 run dockerfile/nginx 2 myNginx
```

This will spin up two containers running Nginx mapping port 80 to 8080.

To stop the container:
```
./src/scripts/cloudcfg.sh stop myNginx
```

To delete the container:
```
./src/scripts/cloudcfg.sh rm myNginx
```

### Running a container (more complete version)
```
cd kubernetes
./src/scripts/cloudcfg.sh -c examples/task.json create /tasks
```

Where task.json contains something like:
```
{
  "ID": "nginx",
  "desiredState": {
    "image": "dockerfile/nginx",
    "networkPorts": [{
      "containerPort": 80,
      "hostPort": 8080
    }]
  },
  "labels": {
    "name": "foo"
  }
}
```

Look in the ```examples/``` for more examples

### Tearing down the cluster
```
cd kubernetes
./src/scripts/kube-down.sh
```

## Development

### Hooks
```
# Before committing any changes, please link/copy these hooks into your .git
# directory. This will keep you from accidentally committing non-gofmt'd
# go code.
cd kubernetes
ln -s "../../hooks/prepare-commit-msg" .git/hooks/prepare-commit-msg
ln -s "../../hooks/commit-msg" .git/hooks/commit-msg
```

### Unit tests
```
cd kubernetes
./src/scripts/test-go.sh
```

### Coverage
```
cd kubernetes
go tool cover -html=target/c.out
```

### Integration tests
```
# You need an etcd somewhere in your path.
# To get from head:
go get github.com/coreos/etcd
go install github.com/coreos/etcd
sudo ln -s "$GOPATH/bin/etcd" /usr/bin/etcd
# Or just use the packaged one:
sudo ln -s "$REPO_ROOT/target/bin/etcd" /usr/bin/etcd
```

```
cd kubernetes
./src/scripts/integration-test.sh
```

### Keeping your development fork in sync
One time after cloning your forked repo:
```
git remote add upstream https://github.com/GoogleCloudPlatform/kubernetes.git
```

Then each time you want to sync to upstream:
```
git fetch upstream
git rebase upstream/master
```

### Regenerating the documentation
Install [nodejs](http://nodejs.org/download/), [npm](https://www.npmjs.org/), and
[raml2html](https://github.com/kevinrenskers/raml2html), then run:
```
cd kubernetes/api
raml2html kubernetes.raml > kubernetes.html
```
