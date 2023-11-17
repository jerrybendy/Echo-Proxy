<template>
  <div>
    <NForm
        ref="formRef"
        :model="settings"
        :rules="rules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
        size="small"
    >
      <NCard title="Ports" size="small" :bordered="false">
        <NFormItem label="HTTP Port" path="httpPort">
          <NInputNumber v-model:value="settings.httpPort" :min="1" :max="65535" :show-button="false" />
          <NText depth="3" style="margin-left: 8px;">Default is 80</NText>
        </NFormItem>
        <NFormItem label="HTTPS Port" path="httpsPort">
          <NInputNumber v-model:value="settings.httpsPort" :min="1" :max="65535" :show-button="false" />
          <NText depth="3" style="margin-left: 8px;">Default is 443</NText>
        </NFormItem>
      </NCard>

      <NCard title="Advanced" size="small" :bordered="false">
        <NFormItem label="Config">
          <NButton @click="openConfigFolder">Open config folder</NButton>
        </NFormItem>
        <NFormItem label="Reload">
          <NButton @click="reloadFrontend">Reload frontend</NButton>
        </NFormItem>
      </NCard>

    </NForm>

    <NLayoutFooter position="absolute" bordered style="padding: 12px 16px">
      <NSpace justify="end">
        <NButton type="primary" @click="saveSettings">Save settings</NButton>
      </NSpace>
    </NLayoutFooter>
  </div>
</template>

<script setup lang="ts">
  import {
    NForm,
    NCard,
    NFormItem,
    NInputNumber,
    NButton,
    NText,
    NLayoutFooter,
    NSpace,
    FormRules,
    FormInst, useMessage
  } from "naive-ui"
  import {GetSettings, OpenConfigFolder, SaveSettings} from "../../wailsjs/go/service/Setting";
  import {WindowReload} from "../../wailsjs/runtime";
  import {onMounted, ref} from "vue";
  import {service} from "../../wailsjs/go/models";

  const message = useMessage()
  const formRef = ref<FormInst | null>(null)
  const settings = ref<service.Setting>({} as service.Setting)
  const rules = ref<FormRules>({
    httpPort: [
      {required: true, message: 'HTTP port is required'},
      {type: "number", min: 1, max: 65535, message: "Invalid port range"},
    ],
    httpsPort: [
      {required: true, message: 'HTTP port is required'},
      {type: "number", min: 1, max: 65535, message: "Invalid port range"},
    ],
  })

  onMounted(async () => {
    settings.value = await GetSettings()
  })

  function openConfigFolder() {
    OpenConfigFolder()
  }

  function reloadFrontend() {
    WindowReload()
  }

  function saveSettings(e: MouseEvent) {
    e.preventDefault()
    formRef.value?.validate(async (errors) => {
      if (errors) {
        console.log(errors)
        return
      }

      const result = await SaveSettings(settings.value)
      if (result) {
        message.success("Save settings successful")
      }
    })
  }

</script>
