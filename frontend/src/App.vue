
<template>
  <NConfigProvider :theme="theme" style="height: 100%">
    <NGlobalStyle />
    <NNotificationProvider placement="top-left">
      <GlobalEvents />
    </NNotificationProvider>

    <NLayout has-sider style="height: 100%">
      <NLayoutSider
          bordered
          collapse-mode="width"
          :collapsed-width="64"
          :width="240"
          show-trigger
          v-model:collapsed="menuCollapsed"
          @update:collapsed="onMenuCollapsedChange"
      >
        <NMenu
            :value="currentMenu"
            :collapsed-width="64"
            :collapsed-icon-size="22"
            :collapsed="menuCollapsed"
            :options="menuOptions"
        />
      </NLayoutSider>

      <NLayout :native-scrollbar="false" style="height: 100%">
        <div style="padding: 12px; min-height: calc(100% - 24px)">
          <NMessageProvider>
            <router-view />
          </NMessageProvider>
        </div>
      </NLayout>
    </NLayout>
  </NConfigProvider>
</template>

<script setup lang="ts">
  import {computed, Component, h, ref, watch, onBeforeMount} from "vue"
  import {RouterLink, useRoute} from "vue-router";
  import {
    NConfigProvider,
    NGlobalStyle,
    useOsTheme,
    darkTheme,
    NIcon,
    NNotificationProvider,
    NLayout,
    NLayoutSider,
    NMenu,
    MenuOption,
    NMessageProvider
  } from "naive-ui"
  import Hosts from "./pages/Hosts.vue";
  import {Home, SettingsSharp} from "@vicons/ionicons5"
  import GlobalEvents from "./components/GlobalEvents.vue";

  const osThemeRef = useOsTheme()
  const route = useRoute()
  const currentMenu = ref("")
  const menuCollapsed = ref(false)

  onBeforeMount(() => {
    const defaultMenuCollapsed = +(window.localStorage.getItem('defaultMenuCollapsed') || '0')
    menuCollapsed.value = !!defaultMenuCollapsed
  })

  const theme = computed(() => {
    return osThemeRef.value === 'dark' ? darkTheme : null
  })

  watch(route, (newRoute) => {
    currentMenu.value = newRoute.name as string || ''
  }, {immediate: true})

  const menuOptions: MenuOption[] = [
    {
      label: () => h(RouterLink, {to: {name: 'Hosts'}}, {default: () => 'Hosts'}),
      key: 'Hosts',
      icon: renderIcon(Home)
    },
    {
      label: () => h(RouterLink, {to: {name: 'Settings'}}, {default: () => 'Settings'}),
      key: 'Settings',
      icon: renderIcon(SettingsSharp)
    },
  ]

  function onMenuCollapsedChange(collapsed: boolean) {
    window.localStorage.setItem('defaultMenuCollapsed', collapsed ? '1' : '0')
    console.log(collapsed)
  }

  function renderIcon (icon: Component) {
    return () => h(NIcon, null, { default: () => h(icon) })
  }
</script>
