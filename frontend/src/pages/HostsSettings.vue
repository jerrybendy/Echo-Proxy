<template>
  <NDrawerContent title="Edit" closable :native-scrollbar="false" footer-style="justify-content: flex-start">
    <NForm
        ref="formRef"
        :model="model"
        :rules="rules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
        size="small"
    >
      <NCard title="Basic" size="small" :bordered="false">
        <NFormItem label="Name" path="name">
          <NInput v-model:value="model.name" placeholder="Domain name"/>
        </NFormItem>
        <NFormItem path="applyHosts">
          <NCheckbox v-model:checked="model.applyHosts">Apply to /etc/hosts file</NCheckbox>
        </NFormItem>
      </NCard>

      <NCard title="Proxy" size="small" :bordered="false">
        <NFormItem label="Default target" path="defaultTarget">
          <NInput v-model:value="model.defaultTarget" placeholder="http://"/>
        </NFormItem>
      </NCard>

      <NCard title="TLS" size="small" :bordered="false">
        <NFormItem path="enableTLS">
          <NCheckbox v-model:checked="model.enableTLS">Enable TLS (HTTPS)</NCheckbox>
        </NFormItem>
        <NFormItem label="Certificate file" path="TLSCertFile" :rule="{required: model.enableTLS, message: 'Certificate file is required'}">
          <NButton :disabled="!model.enableTLS" @click="selectCertificateFile">Select file</NButton>
          <NText style="margin-left: 8px">{{ filename(model.TLSCertFile) }}</NText>
        </NFormItem>
        <NFormItem label="Key file" path="TLSKeyFile" :rule="{required: model.enableTLS, message: 'Key file is required'}">
          <NButton :disabled="!model.enableTLS" @click="selectCertificateKeyFile">Select file</NButton>
          <NText style="margin-left: 8px">{{ filename(model.TLSKeyFile) }}</NText>
        </NFormItem>
      </NCard>


    </NForm>

    <template #footer>
      <NPopconfirm
          :negative-button-props="{type: 'error', ghost: true} as ButtonProps"
          :positive-button-props="{type: 'error'} as ButtonProps"
          @positiveClick="removeHost"
      >
        <template #trigger>
          <NButton type="error" secondary>Delete</NButton>
        </template>
        Do you confirm to remove this host?<br/> This operation can not be revoked!
      </NPopconfirm>
      <NButton type="primary" secondary style="margin-left: auto;" @click="saveSetting">Save</NButton>
    </template>
  </NDrawerContent>
</template>

<script setup lang="ts">
  import {ref, readonly, watch} from 'vue'
  import {
    FormRules,
    NCard,
    NCheckbox,
    NForm,
    NFormItem,
    NInput,
    NButton,
    NDrawerContent,
    FormInst,
    NPopconfirm,
    ButtonProps,
    NText,
  } from 'naive-ui'
  import {service} from "../../wailsjs/go/models";
  import {OpenFileDialog, RemoveHost, SaveSetting} from "../../wailsjs/go/service/Hosts";


  const props = defineProps<{
    host: service.HostConfig | null,
  }>()

  const emit = defineEmits<{
    change: [],
  }>()

  const formRef = ref<FormInst | null>(null)
  const defaultConfig = readonly<service.HostConfig>({
    id: 0,
    name: '',
    applyHosts: false,
    defaultTarget: '',
    enableTLS: false,
    TLSCertFile: '',
    TLSKeyFile: '',
  })
  const model = ref<service.HostConfig>(defaultConfig)

  watch<service.HostConfig | null, true>(() => props.host, (newValue) => {
    if (newValue) {
      model.value = {...newValue}
    } else {
      model.value = {...defaultConfig}
    }
  }, {immediate: true})

  const rules = ref<FormRules>({
    name: [
      {required: true, message: 'Domain name is required'},
      {pattern: /^[-a-z0-9]+(\.[-a-z0-9]+)*\.[a-zA-Z]{2,}$/i, message: 'Invalid domain name'},
    ],
    defaultTarget: [
      {
        pattern: /^https?:\/\/(([a-zA-Z0-9_-])+(\.)?)*(:\d+)?(\/((\.)?(\?)?=?&?[a-zA-Z0-9_-](\?)?)*)*$/i,
        message: 'Default target must be a valid URL'
      },
    ],
  })

  async function saveSetting(e: MouseEvent) {
    e.preventDefault()
    formRef.value?.validate(async (errors) => {
      if (!errors) {
        await SaveSetting(model.value)
        emit('change')

      } else {
        console.log(errors)
      }
    })
  }

  async function removeHost() {
    await RemoveHost(model.value.id)
    emit('change')
  }

  async function selectCertificateFile() {
    const filePath = await OpenFileDialog("Select a certificate file", "Certificate files (*.pem, *.crt)", "*.pem;*.crt")
    if (filePath) {
      model.value.TLSCertFile = filePath
    }
  }

  async function selectCertificateKeyFile() {
    const filePath = await OpenFileDialog("Select a certificate key file", "Key files (*.pem, *.key)", "*.pem;*.key")
    if (filePath) {
      model.value.TLSKeyFile = filePath
    }
  }

  function filename(path: string): string {
    return path.substring(path.replace(/\\/g, '/').lastIndexOf('/') + 1)
  }
</script>

<style scoped>

</style>