Compile:

glide install -v
CGO_ENABLED=0 GOOS=linux go build -a  -ldflags '-extldflags "-static"' -o yrfs-provisioner


ref : https://github.com/kubernetes-incubator/external-storage/tree/master/flex
