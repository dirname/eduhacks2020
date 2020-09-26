layui.define(['tableHTTP', 'form'], function (exports) {
    let $ = layui.$
        , admin = layui.admin
        , view = layui.view
        , table = layui.tableHTTP
        , form = layui.form
        , setter = layui.setter;

    // $.ajaxSetup({
    //     xhrFields: {
    //         withCredentials: true
    //     }
    // })

    table.render({
        elem: '#LAY-app-content-Devices'
        , url: layui.setter.api + '/api/manager/client/get' //模拟接口
        , headers: {
            Authorization: layui.data(setter.tableName)[setter.request.tokenName]
        }
        , cols: [[
            {type: 'numbers', fixed: 'left'}
            , {
                field: 'systemId', minWidth: 100, title: 'systemID', templet: function (value) {
                    if (value.systemId === deviceId) {
                        return '<span class="layui-badge layui-bg-black">This device</span>'
                    } else {
                        return '<span class="layui-badge">' + value.systemId + '</span>'
                    }
                }
            }
            , {field: 'clientId', title: 'clientID'}
            , {
                field: 'user', title: 'ConnUser', templet: function (value) {
                    if (value.user == "") {
                        return "Not logged in"
                    } else {
                        return value.user
                    }
                }
            }
            , {
                field: 'name', title: 'Username', templet: function (value) {
                    if (value.name == "") {
                        return "Not logged in"
                    } else {
                        return value.name
                    }
                }
            }
            , {
                field: 'connectTime', title: 'ConnTime', templet: function (value) {
                    return layui.util.timeAgo(value.connectTime * 1000)
                }
            }
            , {
                field: 'role', title: 'User Role', templet: function (value) {
                    switch (value.role) {
                        default:
                            return '<span class="layui-badge layui-bg-gray">Not logged in</span>'
                        case -1:
                            return '<span class="layui-badge layui-bg-orange">Administrator</span>'
                        case 0:
                            return '<span class="layui-badge layui-bg-gray">Not logged in</span>'
                        case 1:
                            return '<span class="layui-badge layui-bg-green">Academic staff</span>'
                        case 2:
                            return '<span class="layui-badge layui-bg-cyan">Teacher</span>'
                        case 3:
                            return '<span class="layui-badge layui-bg-blue">Student</span>'
                    }
                }
            }
            , {
                title: 'Options',
                minWidth: 300,
                align: 'center',
                fixed: 'right',
                toolbar: '#layuiadmin-app-cont-Devicesbar'
            }
        ]]
        , page: true
        , limit: 10
        , limits: [10, 15, 20, 25, 30]
        , text: {
            none: "No eligible device connected"
        }
    });

    table.on('tool(LAY-app-content-Devices)', function (obj) {
        let data = obj.data;
        if (obj.event === 'Disconnect') {
            layer.confirm('Are you sure to disconnect this client from the server ?', {
                title: "Disconnect ?",
                btn: ['Yes', 'No'],
                yes: function (index) {
                    if (data.systemId == "This device") {
                        layer.msg("You can't disconnect this device")
                    }
                    $.ajax(setter.api + "/api/manager/client/close", {
                        type: 'post',
                        data: JSON.stringify({"clientId": data.clientId}),
                        contentType: "application/json;charset=utf-8",
                        headers: {
                            Authorization: layui.data(setter.tableName)[setter.request.tokenName],
                            SystemID: deviceId,
                        },
                        success: function (res, status) {
                            layer.msg(res.msg)
                        },
                        error: function (e) {
                            layer.open({
                                type: 0,
                                title: "Disconnect Result",
                                content: e.statusText,
                                icon: 6
                            });
                        }
                    });
                    table.reload('LAY-app-content-Devices'); //重载表格
                    layer.close(index);
                }
            });
        } else if (obj.event === 'Msg') {
            admin.popup({
                title: 'Send message to client'
                , area: ['450px', '350px']
                , id: 'LAY-popup-content-send-client'
                , success: function (layero, index) {
                    view(this.id).render('control/send', data).done(function () {
                        form.render(null, 'layuiadmin-form-sender');
                        let clientId = data.clientId;
                        form.on('submit(layuiadmin-app-sender-submit)', function (data) {
                            if (data.field.id === "") {
                                data.field.id = "layerMsgBox"
                            }
                            let sendObj = {}
                            sendObj.clientId = clientId;
                            sendObj.sendUserId = "Admin";
                            sendObj.code = 0;
                            sendObj.msg = "server msg";
                            sendObj.data = JSON.stringify(data.field);
                            $.ajax(setter.api + "/api/manager/client/send", {
                                type: 'post',
                                data: JSON.stringify(sendObj),
                                contentType: "application/json;charset=utf-8",
                                headers: {
                                    Authorization: layui.data(setter.tableName)[setter.request.tokenName],
                                    SystemID: layui.data(setter.tableName)["deviceId"],
                                },
                                success: function (res, status) {
                                    let icon
                                        , content;
                                    if (res.code == 0) {
                                        icon = 6
                                        content = res.msg
                                    } else {
                                        icon = 5
                                        content = "Failed: " + res.msg + " Status: " + status
                                    }
                                    layer.open({
                                        type: 0,
                                        title: "Send Result",
                                        content: content,
                                        icon: icon
                                    });
                                },
                                error: function (e) {
                                    layer.open({
                                        type: 0,
                                        title: "Send Result",
                                        content: e.statusText,
                                        icon: 6
                                    });
                                }
                            })
                            layer.close(index);
                        });
                    });
                }
            });
        }
    });


    exports('control',
        {}
    )
})
;