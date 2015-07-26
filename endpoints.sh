#!/bin/bash

set -exu

API="battery"
URL="https://icumn7abiu.appspot.com/_ah/api/discovery/v1/apis/${API}/v1/rest"
curl -s $URL > ${API}.rest.discovery

# Optionally check the discovery doc
less ${API}.rest.discovery

endpointscfg.py gen_client_lib java -bs gradle ${API}.rest.discovery

unzip ${API}.rest.zip

rm -r android/${API}; mv -f ${API} android/

rm ${API}.rest.discovery ${API}.rest.zip


