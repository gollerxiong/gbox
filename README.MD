gbox -d=./app/models -dsn="goller_admin:123456@tcp(127.0.0.1:3306)/app_api?charset=utf8mb4" -j=true -lp=./app/library/User -t=user_first -pj=appbackend

gmax
介绍
gmax是一个从mysql表中生成go的模型结构体的工具，同时，它自动将模型的操作细分成各个独立的文件，对数据的操作进行收拢，提高代码的封装，统一业务代码规范，提升项目开发效率

软件架构
软件架构说明
-d=./app/models -dsn="goller_admin:123456@tcp(127.0.0.1:3306)/app_api?charset=utf8mb4" -j=true -lp=./app/library/Dao -t=user_info -p=appbackend
安装教程
1. 先下载
    go get github.com/gollerxiong/gbox

2. 安装可执行文件 
   go install github.com/gollerxiong/gbox@latest

   验证是否安装成功：
   which gbox
   如果出现 gbox命令路径则表示安装成功

3. 执行命令生成指定表的脚手架
   3.1 假如你的代码仓库的包名：appbackend (go.mod里面声明的)
   3.2 假如你想将表生成的文件放到你当前目录的：./app/library
   3.3 假如你想把dao层目录放在当前目录的：./app/dao
   3.4 假如你的db链接如右所示：goller_admin:123456@tcp(127.0.0.1:3306)/app_api?charset=utf8mb4

4. 执行命令
   gbox -d=./app/models -dsn="goller_admin:123456@tcp(127.0.0.1:3306)/app_api?charset=utf8mb4" -j=true -lp=./app/library/Dao -t=user_info -p=appbackend

使用说明
-d: 指定模型目录生成路径
-dsn: 数据库连接信息
-j: 生成的表是否包含json tag
-lp: dao文件生成目录
-t: 表名
-p: 个人项目的包名


参与贡献
Fork 本仓库
新建 Feat_xxx 分支
提交代码
新建 Pull Request
特技
使用 Readme_XXX.md 来支持不同的语言，例如 Readme_en.md, Readme_zh.md
Gitee 官方博客 blog.gitee.com
你可以 https://gitee.com/explore 这个地址来了解 Gitee 上的优秀开源项目
GVP 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
Gitee 官方提供的使用手册 https://gitee.com/help
Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 https://gitee.com/gitee-stars/