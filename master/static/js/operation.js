const { createApp, ref } = Vue;
const common = new Utils;
const app = createApp({
    data() {
        return {
            mail: Mail,
            webSite: AddWebSite,
        }
    },
    created:function(){
        let t = this;
        t.mail.getInfo();
    },
    destroyed:function () {
    },
    methods: {
    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');