/** layuiAdmin.pro-v1.4.0 LPPL License By https://www.layui.com/admin/ */
;layui.define(["table", "form"], function (e) {
    var i = (layui.$, layui.admin), t = layui.view, l = layui.table, r = layui.form;
    l.render({
        elem: "#LAY-user-manage",
        url: "./json/useradmin/mycourse.js",
        cols: [[{
            type: "numbers"
        }
            ,{
            field: "id",
            width: 100,
            title: "课程ID",
            sort: !0
        }
            , {
                field: "coursename",
                title: "课程名称",
                minWidth: 100
            }, {
                field: "teacher",
                width: 150,
                title: "教师姓名"
            }, {
                field: "email",
                title: "邮箱"
            }, {
                field: "type",
                width: 90,
                title: "课程类型",
                templet: "#typeTpl",
                align: "center"
            }, {
                field: "opentime",
                title: "开课时间",
                sort: true
            }, {
                field: "status",
                width: 120,
                title: "课程状态",
                templet: "#statusTpl",
                sort: !0,
                align: "center"
            }
            , {
                title: "操作",
                width: 200,
                align: "center",
                fixed: "right",
                toolbar: "#table-useradmin-webuser"
            }]],
        limit: 10,
        height: "full-510",
        text:{
            none: '暂无相关数据'
    }
    }), e("useradmin", {})
});