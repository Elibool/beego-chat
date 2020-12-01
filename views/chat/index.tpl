<div>
    <div class="clearfix"></div>
    <div class="row">
        <div class="col-md-8 col-sm-12 col-md-offset-2">
            <div class="x_panel">
                <div class="x_title">
                    <h2><a href="/"> beego <i class="fa fa-bullhorn"></i> chat </a></h2>
                    <div class="clearfix"></div>
                </div>
                <div class="x_content" style="height: 650px">
                    <div class="col-md-4" style="height:100%;overflow-y: auto ">
                        <ul class="list-unstyled top_profiles scroll-view member_box">
                            <li class="media event online-li li-box-{{.chatInfo.id}} active">
                                <a class="pull-left border-aero profile_thumb">
                                    <i class="fa fa-github-alt blue"></i>
                                </a>
                                <div class="media-body">
                                    <a class="title" href="#">{{.chatInfo.name}}  </a>
                                    <p></p>
                                    <p><small> 在线</small></p>
                                </div>
                            </li>
                            <li class="media event online-li li-box-0">
                                <a class="pull-left border-aero profile_thumb">
                                    <i class="fa fa-comments purple"></i>
                                </a>
                                <div class="media-body">
                                    <a class="title check_member" data-id="0" href="#"> 公共聊天室 </a>
                                    <p></p>
                                    <p><small> 在线</small></p>
                                </div>
                            </li>
                            {{range $index , $item := .lists}}
                                <li class="media event {{if $item.IsOnline}} online-li {{else}} un-online-li {{end}}  li-box-{{$item.Id}}">
                                    <a href="#" data-id="{{$item.Id}}"
                                       class="pull-left border-aero profile_thumb">
                                        <i class="fa fa-github-alt {{$item.color}}"></i>
                                    </a>
                                    <div class="media-body">
                                        <div class="media-body">
                                            <a class="title check_member" href="#"
                                               data-id="{{$item.Id}}">{{$item.Name}}
                                            </a>

                                            <p> {{date $item.LastLoginAt "Y-m-d H:i:s"}} </p>
                                            <p><small> {{if $item.IsOnline}} 在线 {{else}} 已下线 {{end}} </small></p>
                                        </div>
                                    </div>
                                </li>
                            {{end}}
                        </ul>
                    </div>
                    <div class="col-md-8 mail_view" style="height:100%;">
                        <div class="col-md-12">
                            <h4> 聊天室 <i class="fa fa-volume-up"></i> </h4>
                        </div>
                        <div class="col-md-12" id="messages-box-div" style="height: 450px;overflow-y: auto">
                            <ul class="messages" id="messages-box">
                            </ul>
                        </div>
                        <div class="col-md-12">
                            <div class="ln_solid"></div>
                            <textarea id="message" required="required" class="form-control" name="message"></textarea>
                            <br>
                            <button type="button" class="btn msg-send btn-success btn-sm"><i class="fa fa-send"></i> 发 送
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<script>
    WS = window.WS || {};
    MB = window.MB || {};
    let dayTime = 5000;
    $(function () {
        let socket = new WebSocket('ws://' + window.location.host + '/chat/ws?id=' + WS.params.chatId);
        // 当webSocket连接成功的回调函数
        socket.onopen = function () {
            console.log("webSocket open");
            WS.params.connected = true;
        };
        // 断开webSocket连接的回调函数
        socket.onclose = function () {
            console.log("webSocket close");
            WS.params.connected = false;
        };

        //被动接收消息并处理
        socket.onmessage = function (event) {
            let data = JSON.parse(event.data);
            // console.log("receive : ", data);
            if (data.status) {
                WS.handleHsMessage(data);
            }
            if (data.type) {
                MB.handelChatOnline(data)
            }
        };

        //建立聊天请求
        $('.member_box').on('click', '.check_member', function () {
            if ($(this).hasClass('active')) {
                return
            }
            NProgress.start();
            WS.params.acceptId = parseInt($(this).attr('data-id'));
            $('.member_box').children('li').removeClass('active');
            $('.member_box').find(".li-box-" + WS.params.acceptId).addClass('active');
            $('#messages-box').empty();
            let data = {
                accept_id: WS.params.acceptId,
                type: 1,
                msg: ""
            };
            $('.member_box').find(".li-box-" + WS.params.acceptId).find('.badge').remove();

            socket.send(JSON.stringify(data));
            NProgress.done();
        });

        //发送消息
        $('.msg-send').on('click', function () {
            let data = {
                accept_id: WS.params.acceptId,
                type: 2,
                msg: $('#message').val()
            };
            if (WS.params.acceptId < 0) {
                new PNotify({
                    title: "注意",
                    text: "请选择聊天对象",
                    type: "warning",
                    styling: 'bootstrap3',
                    delay: dayTime
                });
                return;
                return;
            }

            if (data.msg == "" || data.msg == " ") {
                new PNotify({
                    title: "注意",
                    text: "信息不能为空",
                    type: "warning",
                    styling: 'bootstrap3',
                    delay: dayTime
                });
                return;
            }
            NProgress.start();
            //公共聊天室
            data.type = data.accept_id > 0 ? 2 : 3;
            socket.send(JSON.stringify(data));
            $('#message').val("");
            NProgress.done();
        })
    });


    // --------------------------------------- 自定义通知类  -----------------------------------------
    MB = {
        handelChatOnline: function (data) {
            //上线或下线通知
            if (data.id == WS.params.chatId) {
                return
            }
            let html = MB.getChatHtml(data);
            if ($('.member_box').find('.li-box-' + data.id) && data.type < 3) {
                //判断是否已在列表
                $('.member_box').find('.li-box-' + data.id).remove();
            }
            if (1 == data.type) {
                //上线
                $('.member_box').find('.li-box-0').after(html);
            } else if (2 == data.type) {
                //下线
                $('.member_box').find('.online-li').last().after(html);
            } else if (3 == data.type) {
                //聊天消息通知
                MB.newMessageNT(data.id);
            }
            MB.onlineNotice(data);
        },
        newMessageNT: function (accept_id) {
            //消息未读提示
            let obj = $('.member_box').find('.li-box-' + accept_id).find('.badge').text();
            if (obj) {
                let num = $('.member_box').find('.li-box-' + accept_id).find('.badge').text();
                num = parseInt(num) + 1;
                $('.member_box').find('.li-box-' + accept_id).find('.badge').text(num);
            } else {
                let html = '<span class="badge bg-orange" style="float: right" >1</span>';
                $('.member_box').find('.li-box-' + accept_id).find('.title').after(html);
            }
        },
        onlineNotice: function (data) {
            //上线提示
            let pType = 2 == data.type ? "warning" : "success";
            let title = data.member.name + (2 == data.type ? "下线了" : "上线了");
            title = 3 == data.type ? "你有一条新消息" : title;

            new PNotify({
                title: title,
                text: "",
                type: pType,
                styling: 'bootstrap3',
                delay: dayTime
            });
        },
        getChatHtml: function (data) {
            if (data.type > 2) {
                return '';
            }
            let online_text = data.type == 1 ? '上线' : '下线';
            let online_class = data.type == 1 ? 'online-li' : 'un-online-li';
            let member_color = data.type == 1 ? 'green' : 'aero';

            let html = '<li class="media event  ' + online_class + ' li-box-' + data.member.id + '">\n' +
                '                                    <a href="#" data-id="14" class="pull-left border-aero profile_thumb" >\n' +
                '                                        <i class="fa fa-github-alt ' + member_color + '"></i>\n' +
                '                                    </a>\n' +
                '                                    <div class="media-body">\n' +
                '                                        <div class="media-body">\n' +
                '                                            <a class="title check_member" href="#" data-id="' + data.member.id + '">' + data.member.name + '</a>\n' +
                '                                            <p> ' + data.member.last_login_at + ' </p>\n' +
                '                                            <p><small>  ' + online_text + '  </small></p>\n' +
                '                                        </div>\n' +
                '                                    </div>\n' +
                '                                </li>';
            return html;
        },
    };


    //-------------------------- 自定义聊天消息类 ------------------------
    WS = {
        params: {
            acceptId: -1,
            chatId: {{.chatInfo.id}},
            connected: false
        },
        handleHsMessage: function (data) {
            if (202 == data.status && data.member_id != WS.chatId && WS.params.acceptId != 0) {
                MB.newMessageNT(0);
                return
            }

            //处理历史数据
            let liHtml = WS.getMessageHtml(data);
            if (200 == data.status) {
                $('#messages-box').prepend(liHtml);
            } else {
                $('#messages-box').append(liHtml);
            }

            //偏移至最底部
            let height = $('#messages-box').height();
            $('#messages-box-div').scrollTop(height + 100);
            console.log(height);
        },
        getMessageHtml: function (data) {
            let headHtml = WS.getChatHeadImg(data.member_id);
            let chatName = WS.getChatName(data.member_id);

            //生成 html
            let html = '<li>\n' +
                '                                    <div class="avatar">' + headHtml + '</div>\n' +
                '                                    <div class="message_wrapper">\n' +
                '                                        <h5 class="heading">' + chatName + '</h5>\n' +
                '                                        <blockquote class="message" style="font-size: 14px">\n' + data.messages + ' </blockquote>\n' +
                '                                        <br>\n' +
                '                                        <p class="url byline"> ' + data.create_at + ' </p>\n' +
                '                                    </div>\n' +
                '                                </li>';

            return html;
        },
        getChatHeadImg: function (chatId) {
            //获取聊天对象头像
            let fa = $('.li-box-' + chatId).find('.fa').attr('class');
            return '<i class="' + fa + '" style="font-size: 26px"></i>';
        },
        getChatName: function (chatId) {
            if (chatId == WS.params.chatId) {
                return '我'
            }
            return $('.li-box-' + chatId).find('.check_member').text()
        }
    }
</script>