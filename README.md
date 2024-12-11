# dplist

dplist 是一个拆图工具. 游戏发布的时候通常会采用和图来提高游戏运行效率, dplist 可以根据和图集描述文件拆分出子图片, 并且还原图片真实大小.

- 支持 TexturePacker 各种版本的 plist 文件导出
- 支持 TexturePacker 部分 json 文件导出
- 支持 fnt 位图字体文件导出
- 支持 spine 的 atlas 文件导出

## 安装

- 首先安装 golang 环境
- 执行 go install github.com/ajenpan/dplist

## 使用说明

```
$ dplist inputdir [-e json, plist, fnt, atlas] [--trimdir]
```

## todolist

- [ ] 支持自定义输出目录
