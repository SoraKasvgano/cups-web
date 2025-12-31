<template>
  <div>
    <h2>Print</h2>
    <form @submit.prevent="uploadAndPrint">
      <div>
        <label>Printer</label>
        <select v-model="printer">
          <option v-for="p in printers" :key="p.uri" :value="p.uri">{{ p.name }} — {{ p.uri }}</option>
        </select>
      </div>
      <div>
        <label>File</label>
        <input type="file" ref="file" />
      </div>
      <button type="submit">Print</button>
    </form>
    <div v-if="msg">{{ msg }}</div>
  </div>
</template>

<script>
export default {
    data() { return { printer: '', msg: '', printers: [] } },
  async mounted() {
    try {
      const resp = await fetch('/api/printers', { credentials: 'include' })
      if (resp.ok) {
        this.printers = await resp.json()
        const last = localStorage.getItem('last_printer')
        if (last) this.printer = last
        else if (this.printers.length > 0) this.printer = this.printers[0].uri
      } else {
        this.msg = 'failed to load printers'
      }
    } catch (e) {
      this.msg = 'failed to load printers: ' + e.message
    }
  },
  methods: {
    async uploadAndPrint() {
      const f = this.$refs.file.files[0]
      if (!f) { this.msg = 'select a file'; return }
      const form = new FormData()
      form.append('file', f)
      form.append('printer', this.printer)
      try {
        const resp = await fetch('/api/print', {
          method: 'POST',
          body: form,
          credentials: 'include',
          headers: { 'X-CSRF-Token': this.getCSRF() }
        })
        if (!resp.ok) throw new Error('print failed')
        const j = await resp.json()
        this.msg = 'Job queued: ' + (j.jobId || '')
        localStorage.setItem('last_printer', this.printer)
      } catch (e) {
        this.msg = e.message
      }
    },
    getCSRF() {
      const m = document.cookie.match('(^|;)\\s*csrf_token\\s*=\\s*([^;]+)')
      return m ? m.pop() : ''
    }
  }
}
</script>
