import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import 'primeicons/primeicons.css'
import './assets/main.css'
import Aura from '@primeuix/themes/aura';
import App from './App.vue'
import router from './router'
import Button from 'primevue/Button'
import Textarea from 'primevue/textarea'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'

import Card from 'primevue/card'
import Chart from 'primevue/chart'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.use(PrimeVue, {
  ripple: true,
  theme: {
        preset: Aura
    }
  
})
app.component('DataTable',DataTable)
app.component('Column',Column)
app.component('Tag',Tag)

app.component('Card',Card)
app.component('Chart',Chart)

app.component('Button',Button)
app.component('Textarea',Textarea)
app.mount('#app')
