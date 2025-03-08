#!/bin/bash

# Scripts to be executed when Raspberry Pi boots
# In order to execute script on boot, /etc/rc.local should be modified to include the line:
# chmod +x /dir_where_script_is_stored/rpi-boot.sh && ./dir_where_script_is_stored/rpi-boot.sh

cd /home/workbench-rpi-admin/repos/pify-player && make start-dev-bg

until [ "`docker inspect -f {{.State.Running}} pify-player-api`"=="true" ]; do
    sleep 0.5;
done;

until [ "`docker inspect -f {{.State.Running}} pify-player-frontend`"=="true" ]; do
    sleep 0.5;
done;

DISPLAY=:0 chromium-browser --no-user-gesture-required https://workbench-rpi.local:5173/player > /home/workbench-rpi-admin/repos/pify-player/scripts/logs.chromium.txt 2>&1
