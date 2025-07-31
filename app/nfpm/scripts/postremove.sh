#!/bin/sh

set -e

case "$1" in
  remove|purge|upgrade|disappear|failed-upgrade|abort-install|abort-upgrade|0)
  ;;

  *)
    echo "postremove.sh called with unknown argument '$1'" >&2
    exit 1
  ;;
esac

systemctl --system daemon-reload >/dev/null || true

exit 0
