const { createApp, ref } = Vue;
const common = new Utils;
const app = createApp({
    data() {
        return {
            mail: Mail,
            webSite: AddWebSite,
            alert: {
                api: "/api/alert/list",
                read: function (id) { return  "/api/alert/read/"+id; },
                del: function (id) {return "/api/alert/del/"+id; },
                list: [],
            },
        }
    },
    created:function(){
        let t = this;
        t.mail.getInfo();
        t.getList();
    },
    destroyed:function () {
    },
    methods: {
        getList: function () {
            let t = this;
            common.AjaxGet(t.alert.api, function (data) {
                t.alert.list = data.data.list;
            });
        },
        readMark: function (id) {
            let t = this;
            common.AjaxGet(t.alert.read(id), function (data) {
                common.ToastShow(data.msg);
                t.getList();
            });
        },
        delAlert: function (id) {
            let t = this;
            common.AjaxGet(t.alert.del(id), function (data){
               common.ToastShow(data.msg);
               t.getList();
            });
        },
        gotoUrl: function (url) {
            console.log(url);
            window.open(url, '_blank');
        },
    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');