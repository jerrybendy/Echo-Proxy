<template>
  <NDataTable
      :columns="columns"
      :data="data"
      size="small"
      striped
      table-layout="fixed"
      :bordered="false"
  />
</template>

<script setup lang="ts">
import {h} from 'vue'
import {DataTableColumn, NCheckbox, NDataTable} from 'naive-ui'

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

const data: HostConfig[] = [
  // {
  //   id: 1,
  //   name: 'localhost',
  //   applyHosts: true,
  // }
]

</script>

<style scoped>

</style>