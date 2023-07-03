#!/bin/bash

echo 'Upgrade main program.'

echo 'Grant executable permission ...'
chmod +x asteroid

cd ../../

echo 'Clean file ...'
rm -f /usr/local/bin/asteroid

echo 'Copy file ...'
cp asteroid /usr/local/bin/

echo 'Daemon reload ...'
systemctl daemon-reload

echo 'Start app ...'
systemctl restart asteroid

echo 'Check status ...'
systemctl status asteroid

echo -e "\033[33;35m ------------------------------------------------- \033[0m"
echo -e "\033[45;30m        upgrade Info        \033[0m"
echo -e "\033[33;35m       data dir：\033[0m" /usr/local/etc/asteroid/
echo -e "\033[33;35m    config file：\033[0m" /usr/local/etc/asteroid/misc/config.yaml
echo -e "\033[33;35m ------------------------------------------------- \033[0m"
