import {createApp} from 'vue'
import router from "./router"
import App from './App.vue'
// import 'vfonts/Lato.css'
// import 'vfonts/FiraCode.css'
// import 'vfonts/Roboto.css'
import './style.css';

const app = createApp(App)
app.use(router)

app.mount('#app')
