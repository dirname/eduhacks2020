# 科技驱动在线教育的新工具、新方法、新技术

## Bingo 是什么:grey_question:

**Bingo 是一个基于区块链、protobuf、websocket通信、渲染，ETCD、RPC实现的分布式高并发教学平台**

疫情期间线上教学显示出不少问题，比如卡顿、掉线，技术上缺乏支持；使用的工具过于分散，平台集成性差，操作过于复杂，平台功能不符合教学要求等。Bingo用科技赋能教育，引领智能教学新时代。用前沿的理念引入科技新技术，实现在线创新教学模式、教学手段和教学结构，打造赋有智能的在线课堂和学习环境，为教师提供方便智能的教学环境，为学生提供高效个性的教育服务

# 学生端

### 在线演示地址 :earth_asia:

> **学生端** :man_student:
> 学生端演示地址：[点击这里](https://www.youtube.com/watch?v=wMm_pOy1feQ&feature=youtu.be)
>
> 项目源码地址：[点击这里](https://github.com/dirname/eduhacks2020)
>
> 项目网页地址：[点击这里](https://htdocs.net/)

### 项目源码地址

* [后端](https://github.com/dirname/eduhacks2020)
* [学生端](https://github.com/dirname/eduhacks2020/tree/web/student)
* [管理端](https://github.com/dirname/eduhacks2020/tree/web/manager)

## 学生端功能 :label:

### 学生端具体功能 :rocket:

- 主页
- 工具栏
- 所有课程
  + 院校课程
  + 优质资源
- 搜索
  + 站内搜索
  + 百度一下
- 数字资源
  + 课程资料
    + 资料下载
    + 资料上传
  + 积分兑换
- 个人中心
  + 学期课程
    + 我的课表
    + 我的作业
    + 答疑中心
    + 成绩查询
    + 重修申请
    + 积分兑换
  + 个人资料
    + 基本资料
    + 修改密码

### 主页

#### 学生端主页功能

+ *学生端功能快捷入口*

__*__ 已报名课程/猜你喜欢课程推送

__*__ 主要内容快捷导航

+ *平台推荐*

__*__ 根据经常搜索和报名的专业等推荐平台

+ *热点内容推送*

__*__ 根据搜索次数和点击率进行推送热点内容

### 工具栏

#### 学生端工具栏

- *消息中心*

__*__ 接收上课、老师回复等消息

+ *便签功能*

__*__ 便签可以随时记录，帮助你记忆/按计划完成需要完成的事

### 所有课程

#### 学生端选课系统

+ *学院开课*

+ *优质慕课*

__*__ 根据专业推荐慕课，让学习不只是单调学习上课知识，增加获取知识的途径

__*__ 

### 搜索

#### 学生端搜索

__*__ 快速百度，搜索想得到的答案

### 数字资源

#### 学生端数字资料

__*__ 资料下载（按需下载）

__*__ 上课所需的资料（如课本pdf、ppt等）

__*__ 作业试题等其他

__*__ 文件上传、编辑和删除

__*__ 

### 学生端积分兑换

__*__ 学生端积分兑换实用的实物，是对学生的一种学习上的鼓励，积分由完成特定的作业获得

### 学生端个人中心

**学生端学期课程**

+ *我的课表*

__*__ 汇总所有报名和收藏的课程

__*__ 进入教室上课/互动，使用webrtc技术实现web端在线屏幕共享、推流播放、音视频会议等，让师生上课更加方便，教师只需在公告发布上课房间即可，websocket实现在线聊天，签到和答题系统（~~待开发~~），让课堂更加活跃，增加师生互动

+ *我的作业*

__*__ 汇总所有课表中的课程发布的作业，完成教师指定作业有积分获取，积分可以换取书籍等实物（可以在**我的课表**中该课程中**作业栏**进行提交，可以在**我的作业**提交/重新提交指定课程的作业）

+ *答疑中心*

__*__ 师生互动版块，让师生交流更加方便

+ *成绩查询*

__*__ 将成绩存储在区块链上，实现期末成绩上链查询，公示各科成绩，防止数据被篡改，提高数据安全性等

+ *重修申请*

__*__ 系统将收录个人的挂科课程的详细信息，学生可以对自己挂科的学科早知道并快速申请重修

## 打包压缩 :hammer:

NodeJS 10.15

```node
npm install
gulp
```

## 参考API :zap:

### API 具体参考

[蚂蚁链API 文档](https://antchain.antgroup.com/docs/11/171401)实现了成绩的存证上链，防止被篡改

[zegoAPI文档](https://doc-zh.zego.im/zh/5416.html)实现了在线推拉流和webrtc音视频在线通话，~~websocket实现在线聊天（待实现）~~

[~~阿里云视频直播API文档~~](https://help.aliyun.com/document_detail/29951.html?spm=5176.7991389.1295213.9.687e1547QjRSPR)~~实现了教师端直播的推流和学生端在线拉流播放~~

## License:page_facing_up:

 [MIT](LICENSE)

## 特别感谢:heart:

所有杰出的 贡献者

[![dirname](https://avatars1.githubusercontent.com/u/32116910?s=64&v=4)](https://github.com/dirname)[![czpei](https://avatars1.githubusercontent.com/u/46366798?s=64&v=4)](https://github.com/czpei)[![github-wander](https://avatars1.githubusercontent.com/u/46366854?s=64&v=4)](https://github.com/github-wander)
