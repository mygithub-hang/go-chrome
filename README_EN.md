{card-default label="前言" width="100%"}
1. 程序基于zserge/lorca感兴趣的大佬可以直接看 [zserge/lorca](https://github.com/zserge/lorca)
2. 基于个人渣渣程序员编码习惯做简单封装，不喜勿喷
3. zserge/lorca简单功能介绍：建立了js和go之间的通讯
4. 项目 [码云地址](https://gitee.com/distance_l/go-gui-mini)
5. 带浏览器一起打包 [看这里](https://blog.ovrvo.com/archives/133.html)
6. Golang交叉编译跨平台打包  [看这里](https://blog.ovrvo.com/archives/4.html)
   {/card-default}

>目录结构

```html
go-gui-mini              项目根目录
├─app                    你写代码的目录
│  ├─bindHandle          go方法声明javaScript方法名目录
│  │  ├─func.go          go方法声明javaScript方法名文件
│
├─boot                   这.....
│  ├─boot.go             额.....
│
├─frame                  做了简单的封装，符合自己的使用习惯
│
├─resources              静态文件目录
│  ├─view                模板文件目录，界面文件
│  │   ├─js
│  │   │  ├─function.js  提供给go的方法 屏蔽右键等普通js文件
│  │   │
│  │   ├─index.html      界面文件
│  │   │
│  │   ├─favicon.png     这东西,就是哪个东西
│  │
│  ├─iconns              win,mac,linux端的图标
│
├─route                  暴露方法名到JS
│
├─main.go                南天门
│
├─build-linux.sh         打包linux程序脚本，感兴趣自己看
│
├─build-macos.sh         打包linux程序脚本，感兴趣自己看
│
├─build-windows.bat      打包linux程序脚本，感兴趣自己看
│
├─main.manifest          自己百度
│
├─main.syso              自己百度
│
```

> 创建主页面：./resources/view/index.html

```html
<!doctype html>
<html lang="en">
<head>
    <title>你没有看错</title>
    <link rel="shortcut icon" href="favicon.png">
</head>
<body>
    <p>爱咋写就咋写</p>
</body>
</html>
```

> 绑定主页
```go
// main.go 文件
// 里面的注释也是有用的别乱改
// 改宽度就行了，主页默认就是index.html
// 窗口宽度
var width int = 480
// 窗口高度
var height int = 320
// 你问我怎么新建窗口?
// <a target="_blank" href="http://www.baidu.com">百度</a>
```

> js调用go方法

```go
// 先创建方法 ./app/bindHandle/func.go
// 不带返回值 参数自己定义
func (c *BindJs) JsToGo(n, b, d int) {
	c.Lock()
	defer c.Unlock()
	fmt.Println("JS调用go：", n, b, d)
}
// 带返回值 参数自己定义
func (c *BindJs) JsToGoGetRet() string {
	c.Lock()
	defer c.Unlock()
	return "golang的返回"
}
// 暴露方法名到javascript ./route/route.go
// Init 暴露方法名到JS
func Init() {
	route.Bind("JsToGo", bindHandle.BindObj.JsToGo)
	route.Bind("JsToGoGetRet", bindHandle.BindObj.JsToGoGetRet)
}
```

```javascript
<script>
    // javascript 异步调用 golang 方法
    JsToGo(1, 2, 3);
    // javascript 调用 golang 方法并获取返回值
    const syncRun = async () => {
        var res = await JsToGoGetRet()
        alert(res)
    };
    syncRun()
    // 贼简单，和使用js自带函数那样
</script>
```

> go调用js方法
```javascript
// javascript 先声明方法
function GoToJs(aa, bb, cc, dd) {
    alert("go调用js:" + aa + ' ' + bb + ' ' + cc + ' ' + dd)
}
```

```go
// 然后go调用
// go 调用 js 方法
// js.Fun("方法名", 参数1, 参数2, ...)
js.Fun("GoToJs", "111", "2222", "333", "444")
```

> 基础用法：日志输出 - go端
```go
// 输出日志到浏览器控制台
log.Log("调用go方法成功")
log.Error()
```

> 启动预览

```cmd
go run main.go
```

> 制作Windows图标

```go
// 1.下载工具包
go get github.com/akavel/rsrc
// 2.生成带图标 syso 文件
rsrc -manifest main.manifest -ico resources/icons/icon.ico -o main.syso
// 3.打包
go build -ldflags '-w -s -H=windowsgui' -o main.exe
```