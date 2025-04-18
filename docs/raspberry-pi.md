# Preparing Raspberry Pi device

## Scripts, configurations

1. chmod +x rpi-boot.sh
2. Edit `~/.config/wayfire.ini` and add following:

```
[autostart]
pify_player=/home/workbench-rpi-admin/repos/pify-player/scripts/rpi-boot.sh
```

3. Amend `/boot/firmware/config.txt` and make the following changes:
   - Set `display_auto_detect=0` to disable auto detection of displays
   - Add line `dtoverlay=vc4-kms-dsi-7inch` to load device tree overlay needed for RPi official 7-inch touchscreen
   - Set `usb_max_current_enable=1` to allow max current for USB peripherals
4. Amend `~/.profile` and add the following block of code at end of file:

```
# only invoke wayfire-pi to start desktop if not logged in thru SSH
if [ -z "$SSH_CLIENT" ]; then
    wayfire-pi
fi
```

5. Make sure Raspberry Pi device boots to console (change via Raspberry Pi Configuration or `sudo raspi-config`)

## (Re)building of apps

1. If there are code changes to either `api` or `player`, `make build` command must be executed.
2. Additionally, if there are changes to `api/cmd/host_handler/main.go` and/or its dependencies, `make install-host-handler` command must be executed.

## Additional changes

### Keep bluetooth always on and discoverable

Open `/etc/rc.local` and add the following lines before `exit 0`

```
# leave bluetooth discoverable on
sudo bluetoothctl <<EOF
power on
discoverable on
EOF
```

### Add virtual keyboard

https://itsfoss.com/raspberry-pi-os-onscreen-keyboard/

```
sudo apt update
sudo apt install squeekboard
```

After installing, enable keyboard through "Display" tab in Raspberry Pi Configuration.

## References

https://www.raspberrypi.com/tutorials/how-to-use-a-raspberry-pi-in-kiosk-mode/
https://www.raspberrypi.com/documentation/accessories/display.html
https://www.raspberrypi.com/documentation/computers/config_txt.html
https://www.reddit.com/r/raspberry_pi/comments/18iwpay/starting_desktop_from_cli_on_rpi_5_startx_command/
