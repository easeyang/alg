# 最简单的霍夫曼对文件进行压缩解压实现
### 目的是熟悉哈夫曼编码解码原理，实现一个简单的文件压缩解压工具

## 简单使用
### 编译
```shell
# 切换到项目根目录执行 会在bin目录下生成可执行文件
make build app=hfm
```

### 压缩
```shell
hfm --action=zip --src=文件名 --dst=目标文件名
```
### 解缩
```shell
hfm --action=unzip --src=已经压缩的文件名 --dst=解压生成的文件名
```

### 压缩文件的头文件格式
```
#  版本  # 创建时间 # 数据位长度 # 映射表数据长度 # 映射表内容      # 数据
#  1字节 # 4字节    # 4字节     # 2字节         #  映射表数据长度  # 
#        #          #          #               #  的值个字节     # 
```