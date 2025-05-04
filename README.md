# ZafkielTools

一些使用[Zafkiel](https://github.com/Aues6uen11Z/Zafkiel)开发项目时可用的工具

## CopyTemplate

将`CopyTemplate.exe`放在项目根目录，点击运行将：

1. 清除`tasks`目录下所有形如`tplxxxxxx.png`的、AirtestIDE截图生成的临时文件。
2. 将`tasks`目录下存在，但`templates`目录下不存在的图片复制到后者。

若要更改这两个目录名，需要添加参数启动，例如

```shell
CopyTemplate.exe -src "C:\path\to\source" -tgt "C:\path\to\target"
```

## ImageWindow

将图片变成窗口，便于调试Zafkiel程序。