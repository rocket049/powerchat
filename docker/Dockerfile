FROM frolvlad/alpine-glibc
MAINTAINER fuhz<fuhuizn@163.com>
        
ADD powerchatserver.tar.gz /
RUN chmod +x /powerchatserver/powerchatserver
WORKDIR /powerchatserver

VOLUME /powerchatserver/dbstore
EXPOSE 7889
ENTRYPOINT ["/powerchatserver/powerchatserver"]

