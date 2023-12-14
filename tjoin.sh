#!/usr/bin/env bash

# tile-join all month .mbtiles files into one yearly .mbtiles files
# eg.
# 2010-1.mbtiles 2010-2.mbtiles ... 2010-12.mbtiles -> 2010.mbtiles
# 2011-1.mbtiles 2011-2.mbtiles ... 2011-12.mbtiles -> 2011.mbtiles
# 2012-1.mbtiles 2012-2.mbtiles ... 2012-12.mbtiles -> 2012.mbtiles

for y in {2010..2023}; do
    set -x
    time tile-join \
        --force \
        --no-tile-size-limit \
        -o ./output/${y}.mbtiles \
        ./output/${y}-*.mbtiles
    { set +x; } 2>&-
done
