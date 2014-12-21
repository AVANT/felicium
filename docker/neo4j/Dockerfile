From ubuntu:14.04

# Get signing key
RUN apt-get update -yyq && apt-get install -yyq wget software-properties-common

#Installing oracle's java
RUN wget -O - http://debian.neo4j.org/neotechnology.gpg.key| apt-key add -
RUN echo 'deb http://debian.neo4j.org/repo stable/' > /etc/apt/sources.list.d/neo4j.list
RUN add-apt-repository ppa:webupd8team/java
RUN apt-get update -yyq
RUN echo debconf shared/accepted-oracle-license-v1-1 select true | debconf-set-selections
RUN echo debconf shared/accepted-oracle-license-v1-1 seen true | debconf-set-selections
RUN apt-get install -yyq oracle-java7-installer

# install neo4j
RUN apt-get install -yyq neo4j

# install bindfs to allow us to rebind if the the data dir is present 
RUN apt-get install -yyq bindfs
RUN sed -i "s|#node_auto_indexing|node_auto_indexing|g" /var/lib/neo4j/conf/neo4j.properties && \
    sed -i "s|#node_keys_indexable|node_keys_indexable|g" /var/lib/neo4j/conf/neo4j.properties && \ 
    echo "remote_shell_host=0.0.0.0" >> /var/lib/neo4j/conf/neo4j.properties

# add the launch
ADD launch.sh /

# expose REST and shell server ports
EXPOSE 7474
EXPOSE 1337

CMD ["/bin/bash", "-c", "/launch.sh"]