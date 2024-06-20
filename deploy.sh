#!/bin/bash
SHELL_DIR=$(
  cd "$(dirname "$0")" || exit
  pwd
)
cd "${SHELL_DIR}" || exit

set -eu

DEST_HOST=43.130.28.11
MODULE=hbase
deploy_pwd=Isd@cloud123

#bash build.sh ${MODULE}

: "${deploy_ssh_port:=22}"

set +u
if [ "${deploy_pwd:=?}" != "?" ]; then
  if [ -e "/usr/local/Cellar/sshpass/1.05/bin/sshpass" ]; then
    SSHPASS_CMD="/usr/local/Cellar/sshpass/1.05/bin/sshpass -p ${deploy_pwd} "
  else
    SSHPASS_CMD="sshpass -p ${deploy_pwd} "
  fi
fi

${SSHPASS_CMD} ssh -oStrictHostKeyChecking=no -p ${deploy_ssh_port} root@"${DEST_HOST}" "mkdir -p /root/go_demos/"
${SSHPASS_CMD} scp -oStrictHostKeyChecking=no -P ${deploy_ssh_port} -r ./pkg/hbase/main.go root@"${DEST_HOST}":/root/go_demos/pkg/hbase/
${SSHPASS_CMD} ssh -oStrictHostKeyChecking=no -p ${deploy_ssh_port} root@"${DEST_HOST}" "chmod 777 -R /root/go_demos/ && chown -R hadoop:hadoop /root/go_demos/"
