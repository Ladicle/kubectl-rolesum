FROM gcr.io/distroless/static-debian10@sha256:8dd1e64607f5037ccde630b7b20bb053c7cc1c8cc7e09c59a3f5a98956b67ed4
COPY kubectl-rolesum /opt/bin/
ENTRYPOINT ["/opt/bin/kubectl-rolesum"]
