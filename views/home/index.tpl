
<div>

    <div class="row">
        <div class="col-md-8 col-sm-12 col-md-offset-2">
            <div class="page-title">
                <div class="title_left">
                    <h3> <i class="fa fa-xing"></i> beego chat demo </h3>
                </div>
            </div>
            <div class="x_panel">
                <div class="x_content">
                    <h3>功能实现</h3>
                    <p> 1. 多人群聊天室 </p>
                    <p> 2. 单对单聊天 </p>
                    <p> 3. 用户上下线提醒 </p>
                    <p> 4. 陌生用户发送消息提醒 </p>
                    <p> 5. 消息未读提醒  </p>
                    <p> 6. 消息已读标记 </p>
                    <p> 7. 修改用户头像和暱称 ; 现为随机默认 （待完成 ）</p>
                    <p> 8. 聊天数据保存时间为 7 天，超时自动清除 </p>

                    <h3> 技术栈 </h3>
                    <p>
                        <span class="badge bg-blue"> beego </span>
                        <span class="badge bg-blue-sky"> mysql </span>
                        <span class="badge bg-red"> redis </span>
                        <span class="badge bg-blue"> jquery bootstrap </span>
                        <span class="badge bg-green"> html5 </span>
                    </p>
                </div>
            </div>
        </div>
    </div>
    <div class="row">
        <div class="col-md-8 col-sm-12 col-md-offset-2">
            <div class="x_panel">
                <div class="x_content">
                    <p> </p>
                    <a href="/chat" class="btn btn-success btn-lg">  进入聊天室 <i class="fa fa-bullhorn"></i> </a>
                </div>
            </div>
        </div>
    </div>
</div>
<script>
    $(function () {
        new PNotify({
            title: "welcome",
            text: "欢迎来到 demo",
            type: "success",
            styling: 'bootstrap3',
            delay: 1500
        });
    });
</script>