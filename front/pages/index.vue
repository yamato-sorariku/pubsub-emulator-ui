<template>
  <div>
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold leading-tight text-gray-900">Messages</h1>
      </div>
    </header>
    <main>
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8"></div>
      <template v-for="(message, index) in messages">
        <PubSubMessage
          :id="message.messageId"
          :key="index"
          :data="message.data"
        />
      </template>
    </main>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { ICloseEvent, w3cwebsocket } from 'websocket'
import PubSubMessage from '~/components/PubSubMessage.vue'
const W3CWebSocket = w3cwebsocket

@Component({
  components: {
    PubSubMessage,
  },
})
@Component
export default class Dashboard extends Vue {
  messages: any[] = []

  mounted() {
    const url = ((location) => {
      const schema = location.protocol === 'https:' ? 'wss:' : 'ws:'
      return `${schema}//${location.host}/ws`
    })(window.location)

    const socket = new W3CWebSocket(url)
    socket.onmessage = (e) => {
      if (typeof e.data === 'string') {
        this.messages = [JSON.parse(e.data), ...this.messages]
      }
    }
    socket.onopen = () => {
      console.log('Socket opened')
    }
    socket.onclose = (e: ICloseEvent) => {
      console.log('Socket closed [' + e.code + ']')
    }
  }
}
</script>
