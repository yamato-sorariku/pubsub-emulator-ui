<template>
  <div class="bg-white shadow overflow-hidden sm:rounded-lg">
    <div class="px-4 py-5 sm:px-6">
      <h3 class="text-lg leading-6 font-medium text-gray-900">ID: {{ id }}</h3>
    </div>
    <div class="border-t border-gray-200">
      <dl>
        <div
          class="bg-white px-4 py-5 sm:grid sm:grid-cols-12 sm:gap-4 sm:px-6"
        >
          <button @click="copyToClipboard">COPY</button>
          <dt class="text-sm font-medium text-gray-500">Data</dt>
          <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-10">
            <span v-html="data" />
          </dd>
        </div>
      </dl>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

@Component
export default class PubSubMessage extends Vue {
  @Prop(String) readonly id!: String
  @Prop(String) readonly data!: String

  copyToClipboard() {
    // @ts-ignore
    this.$copyText(this.data)
  }
}
</script>
