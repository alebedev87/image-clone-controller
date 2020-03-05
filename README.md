# Image clone controller
This controller watches the applications and “caches” the images by re-uploading to your own
registry repository and reconfiguring the applications to use these copies.

## Notice
This controller uses [skopeo](https://github.com/containers/skopeo) utility as a means to mirror images. Which means that it needs to be pre installed to the controller's image.   
The default image from Helm chart already has `skopeo` installed.

## Build the binary
```bash
go build -o image/image-clone-controller cmd/main.go
```

## Run tests
```bash
go test -v ./...
```

## All flags
```bash
Usage of ./image-clone-controller:
      --additional-namespace-blacklist strings   List of namespace(s) which should NOT be watched.
      --backup-registry string                   Backup image registry.
      --img-copy-timeout int                     Timeout for the copy of a single image to the backup registry (in seconds). (default 3600)
      --kubeconfig string                        Paths to a kubeconfig. Only required if out-of-cluster.
      --master --kubeconfig                      (Deprecated: switch to --kubeconfig) The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.
      --registry-org string                      Backup image registry's organization.
      --registry-password string                 Password to access the backup image registry.
      --registry-username string                 Username to access the backup image registry.
```

## Build the image
```bash
VERSION="0.0.1"
REGISTRY="quay.io"
REGISTRY_USER="registry-user"
docker build -t image-clone-controller:${VERSION} -f image/Dockerfile image/
docker tag image-clone-controller:${VERSION} ${REGISTRY}/${REGISTRY_USER}/image-clone-controller:${VERSION}
docker login ${REGISTRY}
docker push ${REGISTRY}/${REGISTRY_USER}/image-clone-controller:${VERSION}
```

## Deploy the controller
```bash
REGISTRY="registry-1.docker.io"
REGISTRY_USERSPACE="myaccount"
REGISTRY_USER="registry-user"
REGISTRY_PWD="registry-password"
helm --kubeconfig ${KUBE_CONFIG} -n ${TARGET_NAMESPACE} \
     --set backupRegistry.name=${REGISTRY} \
     --set backupRegistry.organization=${REGISTRY_USERSPACE} \
     --set backupRegistry.username=${REGISTRY_USER} \
     --set backupRegistry.password=${REGISTRY_PWD} \
     install image-clone-controller charts/image-clone-controller
```

## Deploy the controller without Helm
```bash
REGISTRY="registry-1.docker.io"
REGISTRY_USERSPACE="myaccount"
REGISTRY_USER="registry-user"
REGISTRY_PWD="registry-password"
helm -n ${TARGET_NAMESPACE} \
     --set backupRegistry.name=${REGISTRY} \
     --set backupRegistry.organization=${REGISTRY_USERSPACE} \
     --set backupRegistry.username=${REGISTRY_USER} \
     --set backupRegistry.password=${REGISTRY_PWD} \
     template image-clone-controller charts/image-clone-controller | kubectl apply -f -
```

## Smoke test
```bash
kubectl run nginx --image=nginx --replicas=2
# nginx PODs are supposed to be recreated with the new images in a matter of seconds after they become ready for the first time
kubectl create -f test/ds.yaml
# busybox POD is supposed to be recreated with the new image in a matter of seconds after it becomes ready for the first time
```

## Things to improve
- Use a golang library to talk to image registries
- Manage concurrent image copy from different registed controllers, currently they may turn out to be copying the same image at the same time
- Some sort of integration testing using `envtest` or even a real Kubernetes cluster
- Better test coverage
