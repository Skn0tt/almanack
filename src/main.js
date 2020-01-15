import Vue from "vue";
import router from "./router.js";
import "./filters/index.js";
import "./plugins/font-awesome.js";

import App from "./components/TheApp.vue";

Vue.config.ignoredElements = ["raw-html"];

new Vue({
  router,
  render: h => h(App),
}).$mount("#app");
