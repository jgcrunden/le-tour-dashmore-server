#!/bin/bash

export CONF_DIR=/etc/le-tour-dashmore-server
export CONF_FILE=${CONF_DIR}/server.conf
export SQL_SETUP_FILE=${CONF_DIR}/database-setup.sql
export PASS_FIELD="db.password"

export password=$(grep -Po "${PASS_FIELD}=\K.*$" ${CONF_FILE})
if [ "${password}" == "" ]
then
    export password=$(tr -dc 'A-Za-z0-9' < /dev/urandom | head -c 20)
fi

sudo -u postgres sh -c "psql --set LETOUR_PASS=${password} -f ${SQL_SETUP_FILE}"

sed -i -E "s/${PASS_FIELD}=.*/${PASS_FIELD}=${password}/" ${CONF_FILE}
