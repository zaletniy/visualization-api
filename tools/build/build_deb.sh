#!/bin/bash
set -x
# Package building script
PROJECT_ROOT=/app

rm -rf platform-visibility-all* || true
fpm -s dir -t deb \
  -n platform-visibility-all \
  -v 1.0.0 \
  --after-install /app/tools/build/debian/postinstall \
  $PROJECT_ROOT/etc/platformvisibility/visualization-api/visualization-api.toml=/etc/platformvisibility/visualization-api/ \
  $PROJECT_ROOT/build/linux-amd64/visualizationapi=/usr/bin/

mkdir -p $PROJECT_ROOT/build/deb
mv platform-visibility-all* $PROJECT_ROOT/build/deb/
