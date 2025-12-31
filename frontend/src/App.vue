<template>
  <div class="grid grid-rows-[auto_1fr] min-h-screen w-full bg-base-200">
    <div class="navbar bg-base-100 shadow z-10 px-4 w-full">
      <div class="flex-1">
        <h1 class="text-xl font-bold text-base-content">CUPS 打印</h1>
      </div>
      <div class="flex-none">
        <button v-if="view === 'PrintView'" class="btn btn-sm btn-outline" @click="logout">登出</button>
      </div>
    </div>
    <div class="overflow-auto relative">
      <component :is="view" @login-success="onLogin" @logout="onLogout" />
    </div>
  </div>
</template>

<script>
import LoginView from './views/LoginView.vue'
import PrintView from './views/PrintView.vue'

export default {
  data() {
    return { view: 'LoginView' }
  },
  async mounted() {
    // check existing session on page load; if session present, go straight to PrintView
    try {
      const resp = await fetch('/api/session', { credentials: 'include' })
      if (resp.ok) {
        this.view = 'PrintView'
      }
    } catch (e) {
      // ignore network errors
    }
  },
  components: { LoginView, PrintView },
  methods: {
    onLogin() {
      this.view = 'PrintView'
    },
    onLogout() {
      this.view = 'LoginView'
    },
    async logout() {
      try {
        await fetch('/api/logout', { method: 'POST', credentials: 'include' })
      } catch (e) {
        // ignore errors
      }
      this.onLogout()
    }
  }
}
</script>
