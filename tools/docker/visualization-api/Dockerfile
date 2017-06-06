# visualization-api docker
FROM ubuntu:xenial

COPY build/deb/*.deb /deb/
RUN ls -la /deb
RUN dpkg -i /deb/*.deb

ENTRYPOINT ["visualizationapi"]
