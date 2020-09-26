{
    "code":0
    ,"msg":""
    , "data":[{
        "title": "主页"
        , "icon": "layui-icon-home"
        ,"jump": "index"
    }, {
        "name": "courses"
        , "title": "所有课程"
        , "icon": "layui-icon-flag"
        , "list": [{
            "name": "curriculum"
            , "title": "院校课程"
        ,"list":[{
                "name": "all"
            ,"title": "全部课程"
        }]
        }, {
            "name": "resources"
            , "title": "优质资源"
            , "list": [{
                "name": "mooc"
                , "title": "慕课中心"
            }]
        }]
    }, {
        "name": "template"
        , "title": "搜索"
        , "icon": "layui-icon-search"
        , "list": [{
            "name": "search"
            , "title": "搜索结果"
        ,"jump": "template/search"
        }, {
            "name": "baidu"
            , "title": "百度一下"
            , "spread": true
            , "jump": "/template/link/baidu"
        }]
    }, {
        "name": "app"
        , "title": "数字资源"
        , "icon": "layui-icon-file"
        , "list": [{
            "name": "content"
            , "title": "课程资料"
            , "list": [{
                "name": "download"
                , "title": "资料下载"
            }, {
                "name": "upload"
                , "title": "资料上传"
            }]
        }, {
            "name": "exchange"
            , "title": "积分兑换"
            ,"jump": "app/exchange"
    }]
    }, {
        "name": "user"
        , "title": "个人中心"
        , "icon": "layui-icon-term"
        , "list": [{
            "name": "term"
            , "title": "学期课程"
            , "list": [{
                "name": "mycourses"
                , "title": "我的课表"
            ,"jump": "user/term/mycourses"
        }, {
                "name": "homework"
                , "title": "我的作业"
            ,"jump": "user/term/myhomework"
            }, {
                "name": "quesans"
                , "title": "答疑中心"
            ,"jump": "user/term/msgboard"
            }, {
            "name": "score_query"
            , "title": "成绩查询"
        }, {
            "name": "revision_apply"
            , "title": "重修申请"
        }, {
                "name": "exchange"
                , "title": "积分兑换"
            }]
        }, {
            "name": "admimistrators"
            , "title": "个人资料"
            , "list": [{
                "name": "baseInfo"
                , "title": "基本资料"
            ,"jump": "user/info"
            }, {
                "name": "replys"
                , "title": "修改密码"
                ,"jump": "user/changepwd"
            }]
        }]
    }]
}