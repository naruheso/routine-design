import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'

import App from './App.vue'
import router from './router'

import './assets/styles/variables.css'
import './assets/styles/global.css'

const app = createApp(App)

// Pinia（状態管理）
app.use(createPinia())

// Vue Router
app.use(router)

// PrimeVue（UIフレームワーク）
app.use(PrimeVue, {
  theme: {
    preset: Aura,
  },
})

app.mount('#app')
