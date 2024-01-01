<template>
  <NSpace style="padding: 60px 40px">
    <img class="logo" :src="appIcon" alt="Icon">
    <div style="flex: 1">
      <h1>{{ wailsConfig.name }}</h1>
      <h4>{{ wailsConfig.info.comments }}</h4>
      <p>{{ wailsConfig.info.copyright }}</p>

      <div style="margin-top: 60px;">
        <NButton type="primary" @click="onGithubClick">
          <template #icon>
            <NIcon><LogoGithub /></NIcon>
          </template>
          GitHub
        </NButton>
      </div>
    </div>
  </NSpace>
</template>

<script setup lang="ts">
  import {NSpace, NButton, NIcon} from "naive-ui"
  import appIcon from "../../../build/appicon.png"
  import wailsConfig from "../../../wails.json"
  import {LogoGithub} from "@vicons/ionicons5"
  import {BrowserOpenURL} from "../../wailsjs/runtime";

  function onGithubClick() {
    BrowserOpenURL("https://github.com/jerrybendy/echo-proxy")
  }

</script>

<style scoped>
  .logo {
    width: 100px;
    height: 100px;
    margin-right: 40px;
  }
</style>