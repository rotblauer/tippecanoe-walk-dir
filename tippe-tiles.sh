#!/usr/bin/env bash

out=$1
in=$2
name=$3

echo "out: $out"
echo "in: $in"
echo "name: $name"

tippecanoe --maximum-tile-bytes 330000 \
           --cluster-densest-as-needed \
           --cluster-distance=1 \
           --calculate-feature-density \
           -EElevation:max \
           -ESpeed:max \
           -EAccuracy:mean \
           -EPressure:mean \
           -r1 \
           --minimum-zoom 3 \
           --maximum-zoom 18 \
           -l $name \
           -n $name \
           -o $out \
           --force \
           --read-parallel $in
