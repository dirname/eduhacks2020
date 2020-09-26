/** layuiAdmin.pro-v1.4.0 LPPL License By https://www.layui.com/admin/ */
;layui.define(["form", "upload"], function (t) {
    var i = layui.$, e = layui.layer, n = (layui.laytpl, layui.setter, layui.view, layui.admin), a = layui.form,
        s = layui.upload;
    i("body");
    a.val("base_info-form", {
        "username": 18150090,
        "nickname": "eduHacks2020",
        "sex": 1,
        "cellphone": 13077884433,
        "email": "898577777@qq.com",
        "remarks": "梦想环游世界"
    })
    // var value1 = a.val("base_info-form");
    $('#avatar').attr('src', "https://res.shiguangkey.com/file/201912/27/20191227144115089852982.jpg");
    a.render(), a.verify({
        nickname: function (t, i) {
            return new RegExp("^[a-zA-Z0-9_一-龥\\s·]+$").test(t) ? /(^\_)|(\__)|(\_+$)/.test(t) ? "用户名首尾不能出现下划线'_'" : /^\d+\d+\d$/.test(t) ? "用户名不能全为数字" : void 0 : "用户名不能有特殊字符"
        }, pass: [/^[\S]{8,20}$/, "密码必须8到20位，且不能出现空格"], repass: function (t) {
            if (t !== i("#LAY_password").val()) return "两次密码输入不一致"
        }
    }), a.on("submit(set_website)", function (t) {
        return e.msg(JSON.stringify(t.field)), !1
    }), a.on("submit(set_system_email)", function (t) {
        return e.msg(JSON.stringify(t.field)), !1
    }), a.on("submit(setmyinfo)", function (t) {
        return e.msg(JSON.stringify(t.field)), !1
    });
    var r = i("#LAY_avatarSrc");
    s.render({
        url: "/", elem: "#LAY_avatarUpload",
        before: function (obj) {
            //预读本地文件示例，不支持ie8
            obj.preview(function (index, file, result) {
                $('#preview').append('<img src="' + result + '" alt="' + file.name + '" class="layui-upload-img" style="height: 60px;width: 60px">')
            });
        }
        , done: function (t) {
            0 == t.status ? r.val(t.url) : e.msg(t.msg, {icon: 5})
        }
    }), n.events.avatarPreview = function (t) {
        var i = r.val();
        e.photos({photos: {title: "查看头像", data: [{src: i}]}, shade: .01, closeBtn: 1, anim: 5})
    }, a.on("submit(setmypass)", function (t) {
        return e.msg(JSON.stringify(t.field)), !1
    }), t("set", {})
});