### Usage

### I. Set config
#### A Windows OS config example:
Copy `misc/config.sample.yaml` as `misc/config.yaml` and then edit it.

```yaml
asteroid-data-dir: misc #When initializing the asteroid default data directory, it will generate three directories: 'repository', 'backup', and 'isolation'."
watch-interval: 10 #Watch every 10 seconds
site-list:
  - site-name: site1 #sitename cannot be duplicated
    site-dir: D:\site1 #Specify the site directory to watch
    include-ext: ['.php','.asp'] #Specify the file extensions to watch. Leave [] to watch all files
    exclude-dir: #Ignore watch directories. Leave [] to watch all directories
      - D:\site1\exclude1
      - D:\site1\exclude2
  - site-name: site2
    site-dir: D:\site2
    include-ext: ['.php','.config']  #Watch file with the extensions *.php and *.config. Leave [] to watch all files
    exclude-dir:
      - D:\site2\exclude1
      - D:\site2\exclude2

telegram:
  api-url: https://api.telegram.org
  token: 123:xxxxx
  chat-id: -123
  enable: true

email:
  host: smtp.exmail.qq.com
  port: 465
  username: xxx@qq.cn
  password: 123456
  enable: true
  email-to:
    - aaa@gmail.com
    - bbb@gmail.com
```


### II. How to run:
1. Move the `asteroid.exe` file and the `misc` directory to the parent level of `D:\site1`, such as `D:` (assuming `D:\site1` is the desired directory).
2. Open a terminal window, such as Command Prompt (cmd) or PowerShell.
- #### Init data

`all` can be replace as the specify site name, such as `site1`.
```cmd 
asteroid.exe --act init --site all
```

- #### Monitor site
```cmd 
asteroid.exe --act watch --site all
```
- ####  Uninstall
```cmd 
asteroid.exe --act uninstall --site all
```
- ####  Version
```cmd 
asteroid.exe version
```

### III. Install as service
In windows, you can use [NSSM](https://nssm.cc/download) to install `asteroid` as system service.

>Please initialize the data first ! (Refer to `How to run` > `Init data`)

For example, asteroid was placed on D drive. Install a service as follows.
```cmd
nssm install asteroid "D:\asteroid.exe" "--act watch" "--site all"
```
**_Please run the program after it has been initialized._**

Start service.
```cmd
nssm start asteroid
```
**_While the service is running,you can check logs in the `logs/{date}.log` to ensure the service is actually running._**

Start,stop,restart,remove the service.
```cmd
nssm stop asteroid
nssm restart asteroid
nssm remove asteroid confirm
```

### Linux
```shell
cd build/linux/

chmod +x *.sh
```

### I. Set config
Please refer to the configuration of Windows OS.

### II. Install as system service and start
>Please initialize the data first ! (Refer to `How to run` > `Init data`)

```shell
./install.sh
```
### III. Uninstall the service
```shell
./uninstall.sh
```
