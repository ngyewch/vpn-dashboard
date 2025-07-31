#!/bin/sh

set -e

case "$1" in
  install|upgrade|abort-upgrade|1|2)
  ;;

  *)
    echo "preinstall.sh called with unknown argument '$1'" >&2
    exit 1
  ;;
esac

exit 0
