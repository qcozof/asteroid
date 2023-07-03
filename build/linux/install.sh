#!/bin/bash
echo 'Create dir ...'
mkdir -p /usr/local/etc/asteroid/

cd ../../

echo 'Copy files ...'

file="misc/config.yaml"
if [ ! -f "$file" ]; then
    echo -e "\033[31m Error: Please rename the file 'config.sample.yaml' in the 'misc' directory to 'config.yaml' first, and configure it according to your actual situation.\033[0m"
    exit 1
fi

cp -rf misc/ /usr/local/etc/asteroid/
cp build/linux/asteroid.service /usr/lib/systemd/system/

chmod +x asteroid
cp asteroid /usr/local/bin/

echo 'Reload daemon  ...'
systemctl daemon-reload

echo 'Enable service ...'
systemctl enable asteroid

echo 'Start app ...'
systemctl start asteroid

echo 'Check status ...'
systemctl status asteroid

echo -e "\033[33;35m ------------------------------------------------- \033[0m"
echo -e "\033[45;30m        Install Info        \033[0m"
echo -e "\033[33;35m       data dir：\033[0m" /usr/local/etc/asteroid/
echo -e "\033[33;35m    config file：\033[0m" /usr/local/etc/asteroid/misc/config.yaml
echo -e "\033[33;35m ------------------------------------------------- \033[0m"
