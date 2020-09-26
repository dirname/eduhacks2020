/**

 @Name：layuiAdmin 用户管理 管理员管理 角色管理
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

    //用户管理
    table.render({
        elem: '#LAY-student-manage'
        , url: '/api/manager/student/get' //模拟接口
        , cols: [[
            {type: 'numbers', fixed: 'left'}
            , {field: 'id', width: 100, title: 'ID', sort: true}
            , {field: 'username', title: 'Student ID', minWidth: 100}
            , {field: 'nickname', title: 'Name', minWidth: 100}
            , {
                field: 'college.name', title: 'College', minWidth: 100, templet: function (value) {
                    return value.college.name;
                }
            }
            , {
                field: 'major.name', title: 'Major', minWidth: 100, templet: function (value) {
                    return value.major.name;
                }
            }
            , {
                field: 'class.name', title: 'Class', minWidth: 100, templet: function (value) {
                    return value.class.name;
                }
            }
            , {
                field: 'gender', title: 'Gender', templet: function (value) {
                    if (value.gender) {
                        return "Female"
                    } else {
                        return "Male"
                    }
                }
            }
            , {
                field: 'prison', title: 'Banned', templet: function (value) {
                    if (value.prison) {
                        return "Yes"
                    } else {
                        return "No"
                    }
                }
            }
            , {
                field: 'phone', title: 'Phone', templet: function (value) {
                    if (value.phone === "") {
                        return "Not set"
                    } else {
                        return value.phone
                    }
                }
            }
            , {
                field: 'email', title: 'Email', templet: function (value) {
                    if (value.email === "") {
                        return "Not set"
                    } else {
                        return value.email
                    }
                }
            }
            , {field: 'create', title: 'Create At', sort: true}
            , {title: 'Options', minWidth: 250, align: 'center', fixed: 'right', toolbar: '#table-useradmin-student'}
        ]]
        , page: true
        , limit: 30
        , height: 'full-320'
        , text: {
            none: 'No eligible students'
        }
    });

    //学生管理工具条
    table.on('tool(LAY-student-manage)', function (obj) {
        let data = obj.data;
        if (obj.event === 'del') {
            layer.confirm('Are you sure to delete student: ' + data.nickname + " ?", function (index) {
                data.token = layui.data(setter.tableName)[setter.request.tokenName];
                data.type = 0
                let t = getRequest("/api/manager/student/status", 1, window.location.pathname, data);
                ws.sendRequest(t);
                table.reload('LAY-student-manage'); //重载表格
                layer.close(index);
            });
        } else if (obj.event === 'ban') {
            layer.confirm('Are you sure to ban student: ' + data.nickname + " ?", function (index) {
                data.token = layui.data(setter.tableName)[setter.request.tokenName];
                data.type = 1
                let t = getRequest("/api/manager/student/status", 1, window.location.pathname, data);
                ws.sendRequest(t);
                table.reload('LAY-student-manage'); //重载表格
                layer.close(index);
            });
        } else if (obj.event === 'unblock') {
            layer.confirm('Are you sure to unblock student: ' + data.nickname + " ?", function (index) {
                data.token = layui.data(setter.tableName)[setter.request.tokenName];
                data.type = 2
                let t = getRequest("/api/manager/student/status", 1, window.location.pathname, data);
                ws.sendRequest(t);
                table.reload('LAY-student-manage'); //重载表格
                layer.close(index);
            });
        }
    });

    // 教师管理
    table.render({
        elem: '#LAY-teacher-manage'
        , url: '/api/manager/teacher/get' //模拟接口
        , cols: [[
            {type: 'numbers', fixed: 'left'}
            , {field: 'id', width: 100, title: 'ID', sort: true}
            , {field: 'username', title: 'Teacher ID', minWidth: 100}
            , {field: 'nickname', title: 'Name', minWidth: 100}
            , {
                field: 'gender', title: 'Gender', templet: function (value) {
                    if (value.gender) {
                        return "Female"
                    } else {
                        return "Male"
                    }
                }
            }
            , {
                field: 'fullTime', title: 'Full-time', templet: function (value) {
                    if (value.fullTime) {
                        return "Yes"
                    } else {
                        return "No"
                    }
                }
            }
            , {
                field: 'phone', title: 'Phone', templet: function (value) {
                    if (value.phone === "") {
                        return "Not set"
                    } else {
                        return value.phone
                    }
                }
            }
            , {
                field: 'email', title: 'Email', templet: function (value) {
                    if (value.email === "") {
                        return "Not set"
                    } else {
                        return value.email
                    }
                }
            }
            , {field: 'create', title: 'Create At', sort: true}
            , {title: 'Options', minWidth: 100, align: 'center', fixed: 'right', toolbar: '#table-useradmin-teacher'}
        ]]
        , page: true
        , limit: 30
        , height: 'full-320'
        , text: {
            none: 'No eligible teachers'
        }
    });

    //学生管理工具条
    table.on('tool(LAY-teacher-manage)', function (obj) {
        let data = obj.data;
        if (obj.event === 'del') {
            layer.confirm('Are you sure to delete teacher: ' + data.nickname + " ?", function (index) {
                data.token = layui.data(setter.tableName)[setter.request.tokenName];
                data.type = 0
                let t = getRequest("/api/manager/teacher/del", 1, window.location.pathname, data);
                ws.sendRequest(t);
                table.reload('LAY-teacher-manage'); //重载表格
                layer.close(index);
            });
        }
    });

    exports('useradmin', {})
});