FROM centos
ADD flybook /tmp/
RUN chmod +x /tmp/flybook
ENTRYPOINT /tmp/flybook