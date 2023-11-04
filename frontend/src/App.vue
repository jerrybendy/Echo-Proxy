
<template>
  <NConfigProvider :theme="theme" style="height: 100%">
    <NGlobalStyle />

    <NCard :bordered="false" style="height: 100%">
      <NTabs type="line" animated style="height: 100%" pane-wrapper-style="flex: 1" pane-style="height: calc(100% - var(--n-pane-padding-top))">
        <NTabPane name="hosts" tab="Hosts">
          <Hosts />
        </NTabPane>
        <NTabPane name="settings" tab="Settings">
          Settings
        </NTabPane>
      </NTabs>
    </NCard>


  </NConfigProvider>
</template>

<script setup lang="ts">
  import {computed} from "vue"
  import {NConfigProvider, NGlobalStyle, useOsTheme, darkTheme, NCard, NTabs, NTabPane} from "naive-ui"
  import Hosts from "./tabHosts/Hosts.vue";

  const osThemeRef = useOsTheme()

  const theme = computed(() => {
    return osThemeRef.value === 'dark' ? darkTheme : null
  })

</script>
