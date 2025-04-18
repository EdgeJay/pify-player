#!/bin/bash

# Scripts to be executed when Raspberry Pi boots
# Create a unit file under /lib/systemd/system/sample.service
# Reference: https://www.dexterindustries.com/howto/run-a-program-on-your-raspberry-pi-at-startup/

cd /home/workbench-rpi-admin/repos/pify-player && make start

# host_handler is custom app built using "make install-host-handler"
# Refer to Makefile for more details
host_handler > /home/workbench-rpi-admin/repos/pify-player/scripts/host_handler.log 2>&1

until [ "`docker inspect -f {{.State.Running}} pify-player-api`"=="true" ]; do
    sleep 0.5;
done;

until [ "`docker inspect -f {{.State.Running}} pify-player-frontend`"=="true" ]; do
    sleep 0.5;
done;

sleep 5

DISPLAY=:0 chromium-browser --noerrdialogs \
  --disable-infobars --no-first-run --ozone-platform=wayland \
  --enable-features=OverlayScrollbar --start-maximized --kiosk \
  --no-user-gesture-required https://workbench-rpi.local:5173/player > /home/workbench-rpi-admin/repos/pify-player/scripts/logs.chromium.txt 2>&1