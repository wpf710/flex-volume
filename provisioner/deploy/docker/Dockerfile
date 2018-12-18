FROM launcher.gcr.io/google/debian9
COPY yrfs-provisioner /yrfs-provisioner
RUN chmod +x /yrfs-provisioner

# adapt yrfs driver
RUN touch /usr/local/bin/jq
ENTRYPOINT ["/yrfs-provisioner"]
