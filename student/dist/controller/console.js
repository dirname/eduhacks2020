;layui.define(function (e) {
    layui.use(["admin", "carousel"], function () {
        var e = layui.$, t = (layui.admin, layui.carousel), a = layui.element, i = layui.device();
        e(".layadmin-carousel").each(function () {
            var a = e(this);
            t.render({
                elem: this,
                width: "100%",
                arrow: "none",
                interval: a.data("interval"),
                autoplay: a.data("autoplay") === !0,
                trigger: i.ios || i.android ? "click" : "hover",
                anim: a.data("anim")
            })
        }), a.render("progress")
    }), layui.use("table", function () {
        var e = (layui.$, layui.table);
        e.render({
            elem: "#LAY-index-topSearch",
            url: "./json/console/top-search.js",
            page: true,
            cols: [[{type: "numbers", fixed: "left"}
                , {
                    field: "keywords",
                    title: "关键词",
                    minWidth: 300,
                    templet: '<div><a href="https://www.baidu.com/s?wd={{ d.keywords }}" target="_blank" class="layui-table-link">{{ d.keywords }}</div>'
                }
                , {field: "frequency", title: "搜索次数", minWidth: 120, sort: true}
                , {
                    field: "userNums",
                    title: "用户数",
                    minWidth: 50,
                    sort: true
                }]],
            text: {
                none: '暂无数据！'
            },
            limit: 10,
            limits: [10,],
            skin: "line",
            height: 158
        }), e.render({
            elem: "#LAY-index-topCard",
            url: "./json/console/top-card.js",
            page: true,
            cellMinWidth: 120,
            cols: [[{type: "numbers", fixed: "left"}, {
                field: "title",
                title: "标题",
                minWidth: 300,
                templet: '<div><a href="https://www.baidu.com/s?wd={{ d.title }}" target="_blank" class="layui-table-link">{{ d.title }}</div>'
                // templet: '<div><a href="{{ d.href }}" target="_blank" class="layui-table-link">{{ d.title }}</div>'
            }, {field: "username", title: "发文者"}, {field: "channel", title: "类别"}
            , {field: "crt", title: "点击率", sort: !0, fixed: "right", width: 197
            }]],
            skin: "line",
            limit: 10,
            limits: [10,],
            height: 158
        })
    }), e("console", {})
});