/**

 @Name：layuiAdmin 内容系统
 @Author：star1029
 @Site：http://www.layui.com/admin/
 @License：LPPL

 */


layui.define(['table', 'form'], function (exports) {
    let $ = layui.$
        , admin = layui.admin
        , view = layui.view
        , table = layui.table
        , form = layui.form
        , setter = layui.setter;

    //学院管理
    table.render({
        elem: '#LAY-app-content-college'
        , url: '/api/manager/college/get' //模拟接口
        , page: true
        , cols: [[
            {type: 'numbers', fixed: 'left'}
            , {field: 'id', minWidth: 100, title: 'College ID', sort: true}
            , {field: 'name', title: 'College Name', minWidth: 100}
            , {field: 'create', title: 'Create At', minWidth: 100}
            , {
                title: 'Options',
                width: 150,
                align: 'center',
                fixed: 'right',
                toolbar: '#layuiadmin-app-cont-collegebar'
            }
        ]]
        , text: {
            none: 'No college'
        }
    });

    //专业管理
    table.render({
        elem: '#LAY-app-content-major'
        , url: '/api/manager/major/get' //模拟接口
        , page: true
        , cols: [[
            {type: 'numbers', fixed: 'left'}
            , {field: 'id', minWidth: 100, title: 'Major ID', sort: true}
            , {field: 'name', title: 'Major Name', minWidth: 100}
            , {
                field: 'collegeID', title: 'College ID', hide: true
            }
            , {
                field: 'collegeName', title: "College Name", minWidth: 100
            }
            , {field: 'create', title: 'Create At', minWidth: 100}
            , {
                title: 'Options',
                width: 150,
                align: 'center',
                fixed: 'right',
                toolbar: '#layuiadmin-app-cont-majorbar'
            }
        ]]
        , text: {
            none: 'No major'
        }
    });

    //班级管理
    table.render({
        elem: '#LAY-app-content-class'
        , url: '/api/manager/class/get' //模拟接口
        , page: true
        , cols: [[
            {type: 'numbers', fixed: 'left'}
            , {field: 'id', minWidth: 100, title: 'Class ID', sort: true}
            , {field: 'name', title: 'Class Name', minWidth: 100}
            , {
                field: 'collegeID', title: 'College ID', hide: true
            },
            {
                field: 'majorID', title: 'Major ID', hide: true
            },
            {
                field: 'majorName', title: 'Major Name'
            }
            , {
                field: 'collegeName', title: "College Name", minWidth: 100
            }
            , {field: 'create', title: 'Create At', minWidth: 100}
            , {
                title: 'Options',
                width: 150,
                align: 'center',
                fixed: 'right',
                toolbar: '#layuiadmin-app-cont-classbar'
            }
        ]]
        , text: {
            none: 'No class'
        }
    });

    //监听学院工具条
    table.on('tool(LAY-app-content-college)', function (obj) {
        let data = obj.data;
        if (obj.event === 'del') {
            layer.confirm('Are you sure to delete this college?', {
                title: "Delete ?",
                btn: ['Yes', 'No'],
                yes: function (index) {
                    data.token = layui.data(setter.tableName)[setter.request.tokenName];
                    let t = getRequest("/api/manager/college/del", 1, window.location.pathname, data);
                    ws.sendRequest(t);
                    layui.table.reload('LAY-app-content-college'); //重载表格
                    layer.close(index);
                }
            });
        } else if (obj.event === 'edit') {
            admin.popup({
                title: 'Edit selected college'
                , area: ['450px', '200px']
                , id: 'LAY-popup-content-college'
                , success: function (layero, index) {
                    view(this.id).render('settings/editCollege', data).done(function () {
                        form.render(null, 'layuiadmin-form-college');

                        //监听提交
                        let id = data.id;
                        form.on('submit(layuiadmin-app-college-submit)', function (data) {
                            let field = data.field; //获取提交的字段

                            field.id = id;
                            field.token = layui.data(setter.tableName)[setter.request.tokenName];
                            let t = getRequest("/api/manager/college/edit", 1, window.location.pathname, field);
                            ws.sendRequest(t);
                            layui.table.reload('LAY-app-content-college'); //重载表格
                            layer.close(index);
                        });
                    });
                }
            });
        }
    });


    //监听专业工具条
    table.on('tool(LAY-app-content-major)', function (obj) {
        let data = obj.data;
        if (obj.event === 'del') {
            layer.confirm('Are you sure to delete this major?', {
                title: "Delete ?",
                btn: ['Yes', 'No'],
                yes: function (index) {
                    data.token = layui.data(setter.tableName)[setter.request.tokenName];
                    let t = getRequest("/api/manager/major/del", 1, window.location.pathname, data);
                    ws.sendRequest(t);
                    layui.table.reload('LAY-app-content-major'); //重载表格
                    layer.close(index);
                }
            });
        } else if (obj.event === 'edit') {
            admin.popup({
                title: 'Edit selected major'
                , area: ['450px', '250px']
                , id: 'LAY-popup-content-major'
                , success: function (layero, index) {
                    view(this.id).render('settings/editMajor', data).done(function () {
                        form.render(null, 'layuiadmin-form-major');

                        //监听提交
                        let id = data.id;
                        form.on('submit(layuiadmin-app-major-submit)', function (data) {
                            let field = data.field; //获取提交的字段
                            field.id = id;
                            field.collegeID = parseInt(field.collegeID);

                            // field.id = id;
                            field.token = layui.data(setter.tableName)[setter.request.tokenName];
                            let t = getRequest("/api/manager/major/edit", 1, window.location.pathname, field);
                            ws.sendRequest(t);
                            layui.table.reload('LAY-app-content-major'); //重载表格
                            layer.close(index);
                        });
                    });
                }
            });
        }
    });

    //监听班级工具条
    table.on('tool(LAY-app-content-class)', function (obj) {
        let data = obj.data;
        if (obj.event === 'del') {
            layer.confirm('Are you sure to delete this class?', {
                title: "Delete ?",
                btn: ['Yes', 'No'],
                yes: function (index) {
                    data.token = layui.data(setter.tableName)[setter.request.tokenName];
                    let t = getRequest("/api/manager/class/del", 1, window.location.pathname, data);
                    ws.sendRequest(t);
                    layui.table.reload('LAY-app-content-class'); //重载表格
                    layer.close(index);
                }
            });
        } else if (obj.event === 'edit') {
            admin.popup({
                title: 'Edit selected class'
                , area: ['450px', '310px']
                , id: 'LAY-popup-content-class'
                , success: function (layero, index) {
                    view(this.id).render('settings/editClass', data).done(function () {
                        form.render(null, 'layuiadmin-form-class');

                        //监听提交
                        let id = data.id;
                        form.on('submit(layuiadmin-app-class-submit)', function (data) {
                            let field = data.field; //获取提交的字段
                            field.id = id;
                            field.majorID = parseInt(field.major);

                            // field.id = id;
                            field.token = layui.data(setter.tableName)[setter.request.tokenName];
                            let t = getRequest("/api/manager/class/edit", 1, window.location.pathname, field);
                            ws.sendRequest(t);
                            layui.table.reload('LAY-app-content-class'); //重载表格
                            layer.close(index);
                        });
                    });
                }
            });
        }
    });

    exports('contlist', {})
});