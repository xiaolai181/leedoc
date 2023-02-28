<!DOCTYPE html>
<html lang="zh-cn">
<head>
    <meta charset="utf-8">
    <link rel="shortcut icon" href="{{cdnimg "/favicon.ico"}}">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
    <meta name="renderer" content="webkit" />
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="author" content="MinDoc" />
    <title>{{i18n .Lang "common.new_account"}} - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/css/main.css" "version"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
</head>
<body class="manual-container">
<header class="navbar navbar-static-top smart-nav navbar-fixed-top manual-header" role="banner">
    <div class="container">
        <div class="navbar-header col-sm-12 col-md-6 col-lg-5">
            <a href="{{.BaseUrl}}" class="navbar-brand">{{.SITE_NAME}}</a>
        </div>
    </div>
</header>
<div class="container manual-body">
    <div class="row login">
        <div class="login-body">
            <form role="form" method="post" id="registerForm">
            <input type="hidden" name="_xsrf" value={{ .xsrfdata }}>
                <h3 class="text-center">{{i18n .Lang "common.new_account"}}</h3>
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-addon">
                            <i class="fa fa-user"></i>
                        </div>
                        <input type="text" class="form-control" placeholder="{{i18n .Lang "common.username"}}" name="account" id="account" autocomplete="off">
                    </div>
                </div>
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-addon">
                            <i class="fa fa-lock"></i>
                        </div>
                        <input type="password" class="form-control" placeholder="{{i18n .Lang "common.password"}}" name="password1" id="password1" autocomplete="off">
                    </div>
                </div>
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-addon">
                            <i class="fa fa-lock"></i>
                        </div>
                        <input type="password" class="form-control" placeholder="{{i18n .Lang "common.confirm_password"}}" name="password2" id="password2" autocomplete="off">
                    </div>
                </div>
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-addon" style="padding: 6px 9px;"><i class="fa fa-envelope"></i></div>
                        <input type="email" class="form-control" placeholder="{{i18n .Lang "common.email"}}" name="email" id="email" autocomplete="off">
                    </div>
                </div>

                <div class="form-group">
                    <button type="submit" id="btnRegister" class="btn btn-success" style="width: 100%"  data-loading-text="{{i18n .Lang "message.processing"}}" autocomplete="off">{{i18n .Lang "common.register"}}</button>
                    <br>
                    <br>
                     
                </div>
                
            </form>
           <a href="/login" class="btn btn-success " style="width: 100%" role="button" aria-pressed="true">登陆</a>
        </div>
    </div>
    <div class="clearfix"></div>
</div>
{{template "widgets/footer.tpl" .}}
<!-- Include all compiled plugins (below), or include individual files as needed -->
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/layer/layer.js"}}" type="text/javascript"></script>
<script src="{{cdnjs "/static/js/jquery.form.js"}}" type="text/javascript"></script>
<script type="text/javascript">
    $(function () {
        $("#account,#password1,#password2").on('focus',function () {
            $(this).tooltip('destroy').parents('.form-group').removeClass('has-error');
        });

        $(document).keyup(function (e) {
            var event = document.all ? window.event : e;
            if(event.keyCode === 13){
                $("#btnRegister").trigger("click");
            }
        });
        $("#registerForm").ajaxForm({
            beforeSubmit : function () {
                var account = $.trim($("#account").val());
                var password1 = $.trim($("#password1").val());
                var password2 = $.trim($("#password2").val());
               
                var email = $.trim($("#email").val());

                if(account === ""){
                    $("#account").focus().tooltip({placement:"auto",title : "{{i18n .Lang "message.account_empty"}}",trigger : 'manual'})
                        .tooltip('show')
                        .parents('.form-group').addClass('has-error');
                    return false;

                }
                if(password1 === ""){
                    layer.msg("{{i18n .Lang "message.password_empty"}}");
                    return false;
                }
                if(password2 !== password1){
                    layer.msg("{{i18n .Lang "message.password_not_match"}}");
                    return false;
                }if(email === ""){
                    $("#email").focus().tooltip({title : '{{i18n .Lang "message.email_empty"}}',trigger : 'manual'})
                        .tooltip('show')
                        .parents('.form-group').addClass('has-error');
                    return false;
                }else {

                    $("button[type='submit']").button('loading');
                }
            },
            success : function (res) {
                $("button[type='submit']").button('reset');
                if(res.code === 200){
                    layer.msg(res.message);
                    window.location = "/login";
                }else{
                    layer.msg(res.message);
                    $("button[type='submit']").button('reset');
                }
            },
            error: function () {
                layer.msg('{{i18n .Lang "message.system_error"}}');
                        $("button[type='submit']").button('reset');
                    }
        });
    });
</script>
</body>
</html>