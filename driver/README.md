Deployment

0. yum install jq-1.5-1 on all master and nodes

1. Create vendor directory, do this on all master and nodes.
    mkdir -p /usr/libexec/kubernetes/kubelet-plugins/volume/exec/yr~yrfs
2. Copy yrfs binary to vendor directory, do this on all master and nodes.
    cp ${YRFS_DRIVER}/yrfs /usr/libexec/kubernetes/kubelet-plugins/volume/exec/yr~yrfs
  
    chmod +x ...

3. Restart all kubelet service.
    systemctl restart kubelet
4. Setup yrfs client
    systemctl start yrfs-client

Sample to use yrfs driver

1. Create a pod

    kubectl create -f nginx-yrfs.yaml


    apiVersion: v1
    kind: Pod
    metadata:
      name: nginx-yrfs
      namespace: default
    spec:
      containers:
      - name: nginx-yrfs
        image: nginx
        volumeMounts:
        - name: yrtest
          mountPath: /data
        ports:
        - containerPort: 80
      volumes:
      - name: yrtest
        flexVolume:
          driver: "yr/yrfs"             #vendor is yr
          fsType: "yrfs"                #our fsType is yrfs
          options:
            path: "k8s-vol-1"           #must be set explict, directory will be created at /mnt/yrfs/$path
            accessMode: "ReadWriteMany" #default is ReadWriteMany, value can be set ReadWriteMany/ReadWriteOnce/ReadOnlyMany

2. Create file in mountpoint
    kubectl exec nginx-yrfs -- touch /data/test.log
3. Check the result
    ls /mnt/yrfs/${path}/test.log

Use PV and PVC

1. Create PersistentVolume
    kubectl create -f yrfs-pv.yaml
    kubectl get pv  #check the result
2. Create PersistentVolumeClaim
    kubectl create -f yrfs-pvc.yaml
    kubectl get pvc #check the result, pv will be bound to pvc
3. Create busybox
    kubectl create -f yrfs-busybox-rc.yaml
    kubectl exec pod <pod_name> -- cat /mnt/index.html
