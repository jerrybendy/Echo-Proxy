<template>
  <div class="hosts-container">
    <div class="table-container" @click="onTableClick">
      <div class="toolbar">
        <NButton type="primary" size="small" @click="createEmptyHost"><NIcon><Add/></NIcon> Create</NButton>
        <NButton style="margin-left: auto;" @click="startServer">
          <template #icon>
            <n-icon><Power /></n-icon>
          </template>
          Start server
        </NButton>
      </div>
      <NDataTable
          :columns="columns"
          :data="data"
          size="small"
          striped
          table-layout="fixed"
          :bordered="false"
          :row-class-name="rowClassName"
          style="flex: 1"
      />
    </div>

    <NDrawer v-model:show="isShowEditDrawer" :width="500" placement="right">
      <HostsSettings
          :host="selectedHost"
          @change="isShowEditDrawer = false"
      />
    </NDrawer>
  </div>
</template>

<script setup lang="ts">
import {computed, h, onBeforeUnmount, onMounted, ref} from 'vue'
import {userData} from "../../wailsjs/go/models";
import {DataTableColumn, NButton, NCheckbox, NDataTable, NIcon, NDrawer, NDrawerContent} from "naive-ui";
import {Add, Power} from "@vicons/ionicons5";
import {GetHosts} from "../../wailsjs/go/userData/Hosts";
import HostsSettings from "./HostsSettings.vue";
import {EventsOff, EventsOn} from "../../wailsjs/runtime";
import {StartServer} from "../../wailsjs/go/userData/Service";

const data = ref<userData.HostConfig[]>([])
const selectedHostId = ref<number>(0)
const columns: DataTableColumn[] = [
  {
    title: 'Name',
    key: 'name',
  },
  {
    title: 'hosts file',
    key: 'applyHosts',
    align: "center",
    render(row) {
      return h(NCheckbox, {
        checked: row.applyHosts as boolean,
        size: "small",
      })
    }
  },
]
const isShowEditDrawer = ref(false)

onMounted(() => {
  updateHostsList()

  EventsOn('hostsChange', updateHostsList)
})

onBeforeUnmount(() => {
  EventsOff('hostsChange')
})

const selectedHost = computed<userData.HostConfig | null>(() => {
  return data.value.find(item => item.id === selectedHostId.value) || null
})

async function updateHostsList() {
  data.value = await GetHosts()
}

async function createEmptyHost() {
  selectedHostId.value = 0
  isShowEditDrawer.value = true
}

function rowClassName(row :userData.HostConfig) :string {
  const baseName = `host-id-${row.id} `
  return baseName + (row.id === selectedHostId.value ? 'selected' : '')
}

function onTableClick(event :MouseEvent) {
  let currentNode :HTMLElement = event.target as HTMLElement
  while (currentNode !== event.currentTarget && currentNode !== document.body) {
    if (currentNode?.tagName === 'TR') {
      const matches = currentNode.className.match(/host-id-(\d+)/)
      if (matches) {
        selectedHostId.value = +matches[1]
        isShowEditDrawer.value = true
        break
      }
    }
    currentNode = currentNode.parentNode as HTMLElement
  }
}

async function startServer() {
  await StartServer()
}
</script>

<style scoped>
.hosts-container {
  display: flex;
  height: 100%;
}

.table-container {
  flex: 1;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.toolbar {
  display: flex;
  align-items: center;
  padding: 0 4px 12px;
}
</style>