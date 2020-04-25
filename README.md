## Release-Request-Operator

Acts on release request custom resource object to initiate deployment

### Create New Project
```
operator-sdk new request-operator --repo github.com/kube-sailmaker/request-operator
```

### Create New Api/Controller

```
operator-sdk add api --api-version=deploy.kubesailmaker.io/v1alpha1 --kind=ReleaseRequest
operator-sdk add controller --api-version=deploy.kubesailmaker.io/v1alpha1 --kind=ReleaseRequest
operator-sdk generate k8s
go mod vendor
```

#### Build Operator Image and move to worker node

```
operator-sdk build $USER/release-request-operator:0.1
docker save -o operator.gz $USER/release-request-operator:0.1
scp operator.gz $USER@worker-1:/tmp/
ssh $USER@worker-1 "sudo docker load -i /tmp/operator.gz -q"
```

sed "s|REPLACE_IMAGE|$USER/release-request-operator:0.1|g" deploy/operator.yaml > deploy/_operator.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml

### Create custom resource definition
kubectl create -f deploy/crds/deploy.kubesailmaker.io_releaserequests_crd.yaml

### deploy operator
kubectl create -f deploy/_operator.yaml

### create custom resource
kubectl create -f deploy/crds/deploy.kubesailmaker.io_v1alpha1_releaserequest_cr.yaml
kubectl create -f deploy/crds/release-1.yaml
kubectl create -f deploy/crds/release-2.yaml

### Run locally
```
export OPERATOR_NAME=release-request-operator
operator-sdk run --local
```

### Verify deployment
```
kubectl get deployment
```


### View Pods/Jobs
```
kubectl get jobs,pods
```

