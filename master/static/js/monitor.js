const { createApp, ref } = Vue;
import common from './public.js'
const app = createApp({
    data() {
        return {
            base: {
                api: "/api/slave/info",
                info: {},
                conf: {},
                mail: {},
                saveMailAPI: "/api/mail/conf",
                sendMailAPI: "/api/mail/sendTest",
                saveConfAPI: "/api/conf/save"
            }
        }
    },
    created:function(){
        let t = this;
        t.getData();
    },
    destroyed:function () {
    },
    methods: {
        getData: function () {
            let t = this;
            common.AjaxGet(t.base.api, function (data){
                t.base.info = data.data.info;
                t.base.conf = data.data.conf;
                t.base.mail = data.data.mail;
            });
        },
        saveConf: function () {
            let t = this;
            common.AjaxPost(t.base.saveConfAPI, t.base.conf, function (data) {
                common.ToastShow(data.msg);
            });
        },
        mailParam: function () {
            let t = this;
            t.base.mail.ToList = common.MailToListJoin(t.base.mail.ToList);
            t.base.mail.Port = Number(t.base.mail.Port);
        },
        saveMail: function () {
            let t = this;
            t.mailParam();
            common.AjaxPost(t.base.saveMailAPI, t.base.mail, function (data){
                common.ToastShow(data.msg);
            });
        },
        testSendMail: function () {
            let t =this;
            t.mailParam();
            common.AjaxPost(t.base.sendMailAPI, t.base.mail, function (data){
                common.ToastShow(data.msg);
            });
        },
    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');