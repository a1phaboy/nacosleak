![](https://socialify.git.ci/a1phaboy/nacosleak/image?font=Source%20Code%20Pro&language=1&name=1&owner=1&pattern=Circuit%20Board&stargazers=1&theme=Dark)
## usage
nacosleak.exe -url http://nacosurl [-u]username [-p]password

由于nacos默认部署是不开启auth认证，因此可以尝试不指定-u和-p参数

默认会在当前文件夹生成以下文件：


results/namespace/


=> all_config.txt  所有的配置文件信息


=> passwd.txt      配置文件中所有的密码，可配合密码喷射


=> aksk.txt        配置文件中的aksk
