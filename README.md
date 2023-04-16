### Usage

### I. Set config
#### A windows config example:
edit `misc/config.json`

```json
{
  "siteDir": "D:\\site1",
  "asteroidDataDir": "misc\\",
  "includeExt": ".php|.html",
  "excludeDir": "D:\\site1\\exclude1|D:\\site1\\exclude2",
  "monitorInterval": 10
}
```
      
#### Example Description:

siteDir: D:\\site1, `//Specify the site to watch`   
asteroidDataDir: "misc\\",` //asteroid default data dir,it'll generated 3 dirs, repository,backup and isolution after init `  
includeExt: ".php|.html",`//Watch php files and html files, use "|" as separator, leave empty to watch all files`  
excludeDir: "D:\\site1\\exclude1|D:\\site1\\exclude2",`//Ignore watch dirs, use "|" as separator`  
monitorInterval: 10 `//Watch every 10 seconds`  

**_Note:Windows dir separator '\' need double in config file._**

### II. How to run:
1. Put `asteroid.exe` file and `misc` dir to the parent level of `D:\site1` like `D:\`
2. Open a terminal like cmd or powershell  
#### . Init data

```cmd 
asteroid.exe --act init
```

#### . Monitor site  
```cmd 
asteroid.exe --act watch
```
#### . Uninstall  
```cmd 
asteroid.exe --act uninstall
```
#### . Version  
```cmd 
asteroid.exe version
```

### III. Install as service
In windows, you can use [NSSM](https://nssm.cc/download) to install `asteroid` as system service.  

```cmd
nssm install asteroid "D:\asteroid.exe" "--act watch"
```
**_Please run after init_**  

Start service
```cmd
nssm start asteroid
```
**_While service running check logs/date.log to ensure service actually run._**  

Start,stop,restart,remove service
```cmd
nssm stop asteroid
nssm restart asteroid
nssm remove asteroid confirm
```

### Linux
Expect for the installation service part, the others are almost the same.