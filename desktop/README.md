# Battle Log Desktop

Windows 打包入口在这个目录下。Electron 负责创建桌面窗口并启动 Go 后端，Go 后端负责提供 API 和 `frontend/dist` 静态页面。

## 构建安装包

```powershell
.\desktop\build-windows.ps1
```

产物会输出到：

```text
desktop\release
```

## 本地调试桌面壳

先运行构建脚本生成 `desktop\bin`，然后执行：

```powershell
cd desktop
npm start
```

## 数据文件

打包时会把 `app.ini`、`aion.db`、`storage` 和前端 `dist` 一起复制到桌面应用的后端资源目录。运行时 SQLite 仍然使用相对路径 `aion.db`，上传日志仍然写入 `storage\uploads`。
