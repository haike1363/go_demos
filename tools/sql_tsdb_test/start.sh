#!/bin/bash
SHELL_DIR=$(
  cd "$(dirname "$0")" || exit
  pwd
)
cd "${SHELL_DIR}" || exit

set -eu


# MYAPP_OPENTSDBADDR=":4242" MYAPP_HOST=11.149.47.4 MYAPP_PORT=27009 MYAPP_DBNAME="default" ./sql_tsdb_test -alsologtostderr

MYAPP_OPENTSDBADDR=":4242" MYAPP_HOST=11.149.47.4 MYAPP_PORT=27009 MYAPP_DBNAME="default" nohup ./sql_tsdb_test -alsologtostderr 1>/dev/null 2>&1 &
MYAPP_OPENTSDBADDR=":4243" MYAPP_HOST=9.164.5.112 MYAPP_PORT=27009 MYAPP_DBNAME="default" nohup ./sql_tsdb_test -alsologtostderr 1>/dev/null 2>&1 &
MYAPP_OPENTSDBADDR=":4244" MYAPP_HOST=9.221.118.44 MYAPP_PORT=27009 MYAPP_DBNAME="default" nohup ./sql_tsdb_test -alsologtostderr 1>/dev/null 2>&1 &
MYAPP_OPENTSDBADDR=":4245" MYAPP_HOST=30.43.224.174 MYAPP_PORT=27009 MYAPP_DBNAME="default" nohup ./sql_tsdb_test -alsologtostderr 1>/dev/null 2>&1 &

