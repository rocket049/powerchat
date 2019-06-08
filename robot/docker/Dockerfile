FROM frolvlad/alpine-glibc
MAINTAINER fuhz<fuhuizn@163.com>

RUN mkdir /powerchatrobot
ADD robot.tar.gz /powerchatrobot/
RUN chmod +x /powerchatrobot/robot
WORKDIR /powerchatrobot/

VOLUME /powerchatrobot/data

ENTRYPOINT ["/powerchatrobot/robot", "-u", "服务号Robot", "-p", "xxxxxx"]

