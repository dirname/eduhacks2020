;layui.define(["table", "form"], function (t) {
    var e = (layui.$, layui.admin), i = layui.view, n = layui.table, l = layui.form;
    n.render({
        elem: "#LAY-app-content-list",
        url: "./json/content/upload.js",
        cols: [[{field: "id", width: 100, title: "文件ID", sort: !0},
            {field: "type", width: 110, title: "文件类型", sort: true, templet: "#typeTpl"},
            {field: "title", title: "内容标题"}, {field: "class_name", title: "课程名称"},
            {field: "uploadtime", width: 120, title: "上传时间", sort: !0},
            {field: "status", title: "提交状态", templet: "#buttonTpl", minWidth: 80, align: "center", sort: true},
            {title: "操作", minWidth: 150, align: "center", fixed: "right", toolbar: "#table-content-list"
        }]],
        page: !0,
        limit: 10,
        limits: [10,],
        text: "对不起，加载出现异常！"
    }), n.on("tool(LAY-app-content-list)", function (t) {
        var n = t.data;
        "del" === t.event ? layer.confirm("确定删除该文件吗？", function (e) {
            t.del(), layer.close(e)
        }) : "edit" === t.event && e.popup({
            title: "编辑文件" + "属性",
            area: ["400px", "400px"],
            id: "LAY-popup-content-edit",
            success: function (t, e) {
                i(this.id).render("app/content/upload_file", n).done(function () {
                    l.render(null, "layuiadmin-app-form-list"), l.on("submit(layuiadmin-app-form-submit)", function (t) {
                        t.field;
                        layui.table.reload("LAY-app-content-list"), layer.close(e)
                    })
                })
            }
        })
    }), n.render({
        elem: "#LAY-app-content-tags",
        url: "./json/content/tags.js",
        cols: [[{type: "numbers", fixed: "left"}, {field: "id", width: 100, title: "ID", sort: !0}, {
            field: "tags",
            title: "分类名",
            minWidth: 100
        }, {title: "操作", width: 150, align: "center", fixed: "right", toolbar: "#layuiadmin-app-cont-tagsbar"}]],
        text: "对不起，加载出现异常！"
    }), n.on("tool(LAY-app-content-tags)", function (t) {
        var n = t.data;
        "del" === t.event ? layer.confirm("确定删除此分类？", function (e) {
            t.del(), layer.close(e)
        }) : "edit" === t.event && e.popup({
            title: "编辑分类",
            area: ["450px", "200px"],
            id: "LAY-popup-content-tags",
            success: function (t, e) {
                i(this.id).render("app/content/tagsform", n).done(function () {
                    l.render(null, "layuiadmin-form-tags"), l.on("submit(layuiadmin-app-tags-submit)", function (t) {
                        t.field;
                        layui.table.reload("LAY-app-content-tags"), layer.close(e)
                    })
                })
            }
        })
    }), n.render({
        elem: "#LAY-app-content-comm",
        url: "./json/content/download.js",
        cols: [[{type: "checkbox", fixed: "left"}, {field: "index", width: 50, templet: "#indexTpl", align: "center"}, {
            field: "id",
            width: 100,
            title: "资料ID",
            sort: !0
        }, {field: "teacher", title: "上传教师", width: 150, minWidth: 100}
            , {
                field: "course",
                title: "课程名称",
                width: 250,
                minWidth: 100
            }
            , {field: "content", title: "文件详情", minWidth: 100}, {
                field: "upload_time",
                title: "上传时间",
                minWidth: 100,
                sort: !0
            }, {title: "操作", width: 150, align: "center", fixed: "right", toolbar: "#table-content-com"}]],
        page: !0,
        limit: 10,
        limits: [10,],
        text: "对不起，加载出现异常！"
    }), n.render({
        elem: "#LAY-app-content-homework",
        url: "./json/content/homework.js",
        cols: [[{type: "numbers", fixed: "left"}, {field: "id", width: 100, title: "文件ID", sort: !0}, {
            field: "type",
            width: 110,
            title: "文件类型",
            templet: "#typeTpl",
            minWidth: 50,
            sort: true
        }, {field: "title", title: "文件标题"}, {field: "courseName", title: "课程名称"}, {
            field: "uploadtime",
            width: 120,
            title: "上传时间",
            sort: !0
        }, {field: "status", title: "完成情况", templet: "#buttonTpl", width: 120, align: "center", sort: !0}, {
            title: "操作",
            minWidth: 150,
            align: "center",
            fixed: "right",
            toolbar: "#table-content-list"
        }]],
        page: !0,
        limit: 10,
        limits: [10,],
        text: "对不起，加载出现异常！"
    }), n.on("tool(LAY-app-content-homework)", function (t) {
        var n = t.data;
        "resubmit" === t.event && e.popup({
            title: "重新提交",
            area: ["415px", "500px"],
            id: "LAY-popup-content-comm",
            success: function (t, e) {
                i(this.id).render("app/content/resubmit", n).done(function () {
                    l.render(null, "layuiadmin-form-comment"), l.on("submit(layuiadmin-app-com-submit)", function (t) {
                        t.field; //回调传值
                        layui.table.reload("LAY-app-content-homework"), layer.close(e) //结束,重载
                    })
                })
            }
        })
    }), n.render({
        elem: "#LAY-app-content-score_query",
        url: "./json/content/score.js",
        cols: [[{type: "numbers", fixed: "left"}, {field: "courseID", width: 100, title: "课程ID", align: "center"},
            {field: "studentName", align: "center", title: "学生姓名",width: 120},
            {field: "courseName", title: "课程名称"},
            {field: "credit", title: "课程学分", width: 90, align: "center"},
            {field: "semester", width: 120, title: "开课学期", align: "center", sort: !0},
            {field: "score", title: "成绩", width: 120, align: "center"},
            {field: "pass", title: "评价", templet: "#buttonTpl", width: 120, align: "center", sort: !0},
            {title: "操作",
            minWidth: 150,
            align: "center",
            fixed: "right",
            toolbar: "#table-content-list"
        }]],
        limit: 10,
        limits: [10,],
        text: "对不起，加载出现异常！"
    }), n.on("tool(LAY-app-content-score_query)", function (t) {
        var n = t.data;
        "apply_revision" === t.event && e.popup({
            title: "申请重修",
            area: ["415px", "500px"],
            id: "LAY-popup-content-score",
            success: function (t, e) {
                i(this.id).render("app/content/apply_revision", n).done(function () {
                    l.render(null, "layuiadmin-form-comment"), l.on("submit(layuiadmin-app-com-submit)", function (t) {
                        t.field; //回调传值
                        layui.table.reload("LAY-app-content-homework"), layer.close(e), layer.msg("申请成功，请到指定地点完成缴费") //结束,重载
                    })
                })
            }
        })
    }), t("contlist", {})
});