<template>
  <div style="margin-left: auto;">
    <NButton
        :loading="loading"
        :type="isRunning ? 'success' : 'default'"
        :ghost="isRunning"
        @click="toggleServer"
    >
      <template #icon>
        <n-icon><Power /></n-icon>
      </template>
      {{ isRunning ? 'Stop server' : 'Start server' }}
    </NButton>
  </div>
</template>

<script setup lang="ts">
import {onMounted, onBeforeUnmount, ref, computed} from "vue";
import {Power} from "@vicons/ionicons5";
import {NButton, NIcon, useMessage} from "naive-ui";
import {GetServerStatus, ShutdownServer, StartServer} from "../../wailsjs/go/service/Service";

interface serverStatusType {
  HTTP: boolean
  TLS: boolean
}

const message = useMessage()
const refreshStatusTimer = ref(0)
const loading = ref(false)
const serverStatus = ref<serverStatusType>({} as serverStatusType)

onMounted(() => {
  refreshStatusTimer.value = setInterval(refreshServerStatus, 1000)
})

onBeforeUnmount(() => {
  if (refreshStatusTimer.value) {
    clearInterval(refreshStatusTimer.value)
  }
})

const isRunning = computed(() => {
  return serverStatus.value.HTTP || serverStatus.value.TLS
})

function toggleServer() {
  if (isRunning.value) {
    stopServer()
  } else {
    startServer()
  }
}

async function startServer() {
  loading.value = true
  const result = await StartServer()
  loading.value = false
  refreshServerStatus().then(() => {})
  if (result) {
    message.success("Sever is started", {
      keepAliveOnHover: true,
    })
  } else {
    message.error("Start server failed")
  }
}
async function stopServer() {
  loading.value = true
  await ShutdownServer()
  loading.value = false
  refreshServerStatus().then(() => {})
  message.success("Sever is shutdown", {
    keepAliveOnHover: true,
  })
}

async function refreshServerStatus() {
  serverStatus.value = (await GetServerStatus()) as unknown as serverStatusType
}

</script>

<style scoped>

</style>