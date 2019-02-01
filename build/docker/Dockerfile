FROM alpine:3.7

EXPOSE 8046
EXPOSE 7948

RUN set -xe; \
    apk update && apk upgrade; \
    # add some usefull tools
    apk add --no-cache bash util-linux coreutils findutils grep \
    openssh curl bind-tools wget; \
    apk add --no-cache linux-headers; \
    apk add --update ca-certificates; 


RUN mkdir  /sample-raft
RUN mkdir  /sample-raft/data/

COPY ./entrypoint.sh /sample-raft/entrypoint.sh
COPY ./node /sample-raft/node

RUN adduser -D -u 1001 application;
RUN chmod +s /bin/ping

RUN chmod +x /sample-raft/entrypoint.sh
RUN chmod +x /sample-raft/node
RUN chown -R 1001 /sample-raft




USER 1001

CMD ["/sample-raft/entrypoint.sh"]
