const { createApp, ref } = Vue;
const common = new Utils;
const app = createApp({
    data() {
        return {
            mail: Mail,
            webSite: AddWebSite,
            alert: {
                api: "/api/alert/list",
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
    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');