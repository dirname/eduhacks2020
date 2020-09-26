/**

 @Name：layuiAdmin 公共业务
 @Author：贤心
 @Site：http://www.layui.com/admin/
 @License：LPPL

 */

layui.define(function (exports) {
    let $ = layui.$
        , layer = layui.layer
        , laytpl = layui.laytpl
        , setter = layui.setter
        , view = layui.view
        , admin = layui.admin

    //公共业务的逻辑处理可以写在此处，切换任何页面都会执行
    //……


    //退出
    admin.events.logout = function () {
        let obj = {}
        obj.token = layui.data(setter.tableName)[setter.request.tokenName];
        obj.systemId = layui.data(setter.tableName)["deviceId"];
        let t = getRequest("/api/logout", 1, window.location.pathname, obj)
        if (!ws.isOpened) {
            layer.open({
                type: 0
                , title: "无法执行请求"
                , content: "无法执行请求, the communication with the server is disconnected, please refresh the page and try again"
                , icon: 5
            })
            return
        }
        ws.sendRequest(t).then(response => {
            admin.exit();
        });
        //执行退出接口
        // admin.req({
        //   url: './json/user/logout.js'
        //   ,type: 'get'
        //   ,data: {}
        //   ,done: function(res){ //这里要说明一下：done 是只有 response 的 code 正常才会执行。而 succese 则是只要 http 为 200 就会执行
        //
        //     //清空本地记录的 token，并跳转到登入页
        //     admin.exit();
        //   }
        // });
    };


    //对外暴露的接口
    exports('common', {});
});