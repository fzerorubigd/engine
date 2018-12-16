FROM phusion/baseimage:latest

COPY ./install.sh /bd_build/install.sh

RUN chmod a+x /bd_build/install.sh && \
    /bd_build/install.sh && \
    /bd_build/cleanup.sh


