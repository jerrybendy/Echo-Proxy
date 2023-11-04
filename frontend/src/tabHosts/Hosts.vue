<template>
  <div class="hosts-container">
    <div class="table-container">
      <HostsTable />
    </div>
    <div class="options-container">
      <HostsSettings />
    </div>
  </div>
</template>

<script setup lang="ts">
  import HostsTable from "./HostsTable.vue";
  import HostsSettings from "./HostsSettings.vue";
</script>

<style scoped>
  .hosts-container {
    display: flex;
    height: 100%;
  }
  .table-container {
    width: 420px;
    height: 100%;
  }
</style>