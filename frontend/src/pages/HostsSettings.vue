<template>
  <NDrawerContent title="Edit" closable :native-scrollbar="false" footer-style="justify-content: flex-start">
<!--    <NTabs type="line" animated default-value="basic">-->
<!--      <NTab name="basic" tab="Basic" />-->
<!--      <NTab name="proxy" tab="Proxy" />-->
<!--      <NTab name="ssl" tab="SSL" />-->
<!--    </NTabs>-->

    <NForm
        ref="formRef"
        :model="model"
        :rules="rules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
        size="small"
    >
      <NCard id="basic" title="Basic" size="small" :bordered="false">
        <NFormItem label="Name" path="name">
          <NInput v-model:value="model.name" placeholder="www.example.com"/>
        </NFormItem>
        <NFormItem path="applyHosts">
          <NCheckbox v-model:checked="model.applyHosts">Apply to /etc/hosts file</NCheckbox>
        </NFormItem>
      </NCard>

      <NCard id="proxy" title="Proxies" size="small" :bordered="false">
        <template #header-extra>
          <NTooltip content-style="padding: 4px 6px; background: #262626; border-radius: 4px" raw>
            <template #trigger>
              <NButton quaternary size="tiny" @click="addProxy">
                <template #icon><NIcon><Add /></NIcon></template>
              </NButton>
            </template>
            <span style="font-size: 12px">Add proxy</span>
          </NTooltip>
        </template>
        <ProxySettingCard
            v-for="(proxy, index) in model.proxies" :key="proxy.id"
            :proxy="proxy" :count="model.proxies.length" :index="index"
            @moveUp="proxyMove(index, -1)" @moveDown="proxyMove(index, 1)" @remove="proxyMove(index, 0)"
        />
      </NCard>

      <NCard title="SSL" size="small" :bordered="false">
        <NFormItem path="enableTLS">
          <NCheckbox v-model:checked="model.enableTLS">Enable SSL (HTTPS)</NCheckbox>
        </NFormItem>
        <NFormItem label="Certificate" path="TLSCertFile" :rule="{required: model.enableTLS, message: 'Certificate file is required'}">
          <NButton :disabled="!model.enableTLS" @click="selectCertificateFile">Select file</NButton>
          <NText style="margin-left: 8px">{{ filename(model.TLSCertFile) }}</NText>
        </NFormItem>
        <NFormItem label="Key" path="TLSKeyFile" :rule="{required: model.enableTLS, message: 'Key file is required'}">
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
  import {ref, watch} from 'vue'
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
    NIcon,
    NTooltip,
  } from 'naive-ui'
  import {service} from "../../wailsjs/go/models";
  import {OpenFileDialog, RemoveHost, SaveSetting} from "../../wailsjs/go/service/Hosts";
  import {Add} from "@vicons/ionicons5";
  import deepmerge from "deepmerge";
  import ProxySettingCard from "../components/ProxySettingCard.vue";

  const props = defineProps<{
    host: service.HostConfig | null,
  }>()

  const emit = defineEmits<{
    change: [],
  }>()

  const formRef = ref<FormInst | null>(null)
  const defaultConfig: Readonly<service.HostConfig> = service.HostConfig.createFrom({
    id: 0,
    name: '',
    applyHosts: false,
    enableTLS: false,
    TLSCertFile: '',
    TLSKeyFile: '',
    proxies: [] as service.HostProxy[],
  })
  const defaultProxy: Readonly<service.HostProxy> = service.HostProxy.createFrom({
    id: Date.now(),
    matchType: "PREFIX",
    matchRule: "/",
    // matchParams: {} as service.HostMatchParams,
    target: "",
    changeOrigin: false,
  })

  const model = ref<service.HostConfig>(deepmerge({}, defaultConfig))

  watch<service.HostConfig | null, true>(() => props.host, (newValue) => {
    if (newValue) {
      model.value = deepmerge(defaultConfig, newValue)
    } else {
      model.value = deepmerge({}, defaultConfig)
    }
    if (!model.value.proxies || !model.value.proxies.length) {
      model.value.proxies = [deepmerge<service.HostProxy>(defaultProxy, {id: Date.now()})]
    }
  }, {immediate: true})

  const rules = ref<FormRules>({
    name: [
      {required: true, message: 'Domain name is required'},
      {pattern: /^[-a-z0-9]+(\.[-a-z0-9]+)*\.[a-zA-Z]{2,}$/i, message: 'Invalid domain name'},
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

  function addProxy() {
    model.value.proxies.unshift(deepmerge<service.HostProxy>(defaultProxy, {id: Date.now()}))
  }

  function proxyMove(index: number, delta: number) {
    const items = model.value.proxies.splice(index, 1)
    delta !== 0 && model.value.proxies.splice(index + delta, 0, items[0])
  }

</script>
