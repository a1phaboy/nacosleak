[](banner.png)
## usage
nacosleak.exe -url http://nacosurl [-u]username [-p]password

由于nacos默认部署是不开启auth认证，因此可以尝试不指定-u和-p参数，默认会在当前文件夹生成以下文件：
results/
=> all_config.txt  所有的配置文件信息
=> passwd.txt      配置文件中所有的密码，可配合密码喷射