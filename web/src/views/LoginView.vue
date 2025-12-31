<template>
  <div>
    <h2>Login</h2>
    <form @submit.prevent="login">
      <div>
        <label>Username</label>
        <input v-model="username" />
      </div>
      <div>
        <label>Password</label>
        <input type="password" v-model="password" />
      </div>
      <button type="submit">Login</button>
    </form>
    <div v-if="error" style="color:red">{{ error }}</div>
  </div>
</template>

<script>
export default {
  data() { return { username: '', password: '', error: '' } },
  methods: {
    async login() {
      this.error = ''
      try {
        const resp = await fetch('/api/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username: this.username, password: this.password }),
          credentials: 'include'
        })
        if (!resp.ok) {
          this.error = 'Login failed'
          return
        }
        this.$emit('login-success')
      } catch (e) {
        this.error = e.message
      }
    }
  }
}
</script>
