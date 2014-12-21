#!/bin/bash

set -ex

NEO4J_HOME=/var/lib/neo4j

if [ -d /bindmount/data ] ; then
    bindfs --force-user=neo4j -o nonempty /bindmount/data /var/lib/neo4j/data
fi
sed -i "s|#org.neo4j.server.webserver.address=0.0.0.0|org.neo4j.server.webserver.address=$HOSTNAME|g" $NEO4J_HOME/conf/neo4j-server.properties
ulimit -n 65536 ; .$NEO4J_HOME/bin/neo4j console