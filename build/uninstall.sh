#!/bin/bash
echo 'Begin uninstall ...'
systemctl stop asteroid
systemctl disable asteroid

echo 'Delete files ...'
rm -f /usr/lib/systemd/system/asteroid.service
rm -f /usr/local/bin/asteroid
rm -f /usr/local/etc/asteroid/config.yaml

echo 'end'
