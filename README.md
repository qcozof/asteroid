### Usage

### I. Set config
#### A windows config example:
edit `misc/config.yaml`

```yaml
asteroid-data-dir: misc #asteroid default data dir,it'll generated 3 dirs, repository,backup and isolution after init
monitor-interval: 10 #Watch every 10 seconds
site-list:
  - site-name: site1 #sitename cannot duplicate
    site-dir: D:\site1 #Specify the site dir to watch
    include-ext: ['.php','.asp'] #Watch file extensions, leave [] to watch all files
    exclude-dir: #Ignore watch dirs, leave [] to watch all dirs
      - D:\site1\exclude1
      - D:\site1\exclude2
  - site-name: site2
    site-dir: D:\site2
    include-ext:  #Watch file *.php and *.config, leave [] to watch all files
      - .php
      - .config
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
1. Put `asteroid.exe` file and `misc` dir to the parent level of `D:\site1` like `D:\`
2. Open a terminal like cmd or powershell  
- #### Init data

`all` can be replace as specify site name, such as `site1`.  
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

```cmd
nssm install asteroid "D:\asteroid.exe" "--act watch" "--site all"
```
**_Please run after init_**  

Start service
```cmd
nssm start asteroid
```
**_While service running check logs/{date}.log to ensure service actually run._**  

Start,stop,restart,remove service
```cmd
nssm stop asteroid
nssm restart asteroid
nssm remove asteroid confirm
```

### Linux
>` Expect for the installation service part, the others are almost the same.`

### I. Set config
Refer to windows configuration.

### II. Install as service and start
```shell
chmod +x install.sh
./install.sh
```
### III. Uninstall service
```shell
chmod +x uninstall.sh
./uninstall.sh
```



