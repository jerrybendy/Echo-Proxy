<template>
  <NCard title="Match" size="small" :bordered="true"
         style="margin-bottom: 12px;"
         content-style="padding-bottom: 0"
         footer-style="padding-top: 0"
         :segmented="{footer: 'soft'}"
  >
    <template #header-extra>
      <ProxyMoveButtons
          :index="index" :count="count"
          @moveUp="emit('moveUp')" @moveDown="emit('moveDown')" @remove="emit('remove')"
      />
    </template>
    <NFormItem :path="`proxies.${index}.matchRule`" :show-require-mark="false" :rule="formRules.matchRule">
      <div>
        <NInputGroup>
          <NSelect v-model:value="proxy.matchType" :options="matchType" style="width: 140px" />
          <NInput v-model:value="proxy.matchRule" />
        </NInputGroup>
<!--        <NCheckbox v-if="proxy.matchType === 'REGEXP'" v-model:checked="proxy.matchParams.caseInsensitive" style="margin: 8px 0 0 4px;">Case-insensitive</NCheckbox>-->
      </div>
    </NFormItem>

    <template #footer>
      <h3>Target</h3>
      <NFormItem label="Target type" :path="`proxies.${index}.targetType`" :rule="formRules.targetType">
        <NSelect v-model:value="proxy.targetType" :options="targetType"/>
      </NFormItem>

      <!-- Target type PROXY -->
      <div v-if="proxy.targetType === 'PROXY'">
        <NFormItem label="Proxy pass" :path="`proxies.${index}.targetParams.proxyPass`" :rule="formRules.proxyPass">
          <NInput v-model:value="proxy.targetParams.proxyPass" placeholder="http://127.0.0.1:8080"/>
        </NFormItem>
        <NSpace vertical style="padding-left: 4px">
          <NCheckbox v-model:checked="proxy.targetParams.changeOrigin">Change origin</NCheckbox>
<!--          <NCheckbox >X-Forwarded-*</NCheckbox>-->
<!--          <NCheckbox >Change cookie domain</NCheckbox>-->
        </NSpace>
      </div>

      <!-- Target type STATIC -->
      <div v-if="proxy.targetType === 'STATIC'">
        <NFormItem label="Document root" :path="`proxies.${index}.targetParams.documentRoot`" :rule="formRules.documentRoot">
          <div>
            <NButton @click="selectDocumentRoot">Select folder</NButton>
            <NText style="margin-left: 8px; word-break: break-all">{{ proxy.targetParams.documentRoot }}</NText>
          </div>
        </NFormItem>
        <NFormItem label="Fallback" :path="`proxies.${index}.targetParams.fallback`">
          <NInput v-model:value="proxy.targetParams.fallback" placeholder="/index.html"/>
        </NFormItem>
      </div>

      <!-- Target type PHP -->
      <div v-if="proxy.targetType === 'PHP'">
        <NFormItem label="Document root" :path="`proxies.${index}.targetParams.documentRoot`" :rule="formRules.documentRoot">
          <div>
            <NButton @click="selectDocumentRoot">Select folder</NButton>
            <NText style="margin-left: 8px; word-break: break-all">{{ proxy.targetParams.documentRoot }}</NText>
          </div>
        </NFormItem>
        <NFormItem label="FPM Address" :path="`proxies.${index}.targetParams.fpmAddress`">
          <NInput v-model:value="proxy.targetParams.fpmAddress" placeholder="127.0.0.1:9000"/>
        </NFormItem>
        <NFormItem label="Fallback" :path="`proxies.${index}.targetParams.fallback`">
          <NInput v-model:value="proxy.targetParams.fallback" placeholder="/index.html"/>
        </NFormItem>
      </div>


    </template>
  </NCard>
</template>

<script setup lang="ts">
  import {ref} from "vue"
  import {
    NCard,
    NCheckbox,
    NFormItem,
    NInput,
    NInputGroup,
    NSelect,
    NSpace,
    FormItemRule,
    NButton,
    NText
  } from "naive-ui";
  import ProxyMoveButtons from "./ProxyMoveButtons.vue";
  import {service} from "../../wailsjs/go/models";
  import {OpenFolder} from "../../wailsjs/go/service/Hosts";

  const props = defineProps<{
    proxy: service.HostProxy,
    index: number,
    count: number,
  }>()

  const emit = defineEmits<{
    moveUp: [],
    moveDown: [],
    remove: [],
  }>()

  const matchType = ref([
    {label: 'PREFIX', value: 'PREFIX'},
    // {label: 'REGEXP', value: 'REGEXP'},
    // {label: 'GLOB', value: 'GLOB'},
  ])

  const targetType = ref([
    {label: 'Reverse Proxy', value: 'PROXY'},
    {label: 'Static files', value: 'STATIC'},
    {label: 'PHP', value: 'PHP'},
  ])

  const formRules: { [key: string]: FormItemRule | FormItemRule[] } = {
    matchRule: {required: true, message: 'Match rule is required'},
    targetType: {required: true, message: 'Target type is required'},
    proxyPass: [
      {required: true, message: 'Proxy pass is required'},
      {
        pattern: /^https?:\/\/(([a-zA-Z0-9_-])+(\.)?)*(:\d+)?(\/((\.)?(\?)?=?&?[a-zA-Z0-9_-](\?)?)*)*$/i,
        message: 'Proxy pass must be a valid URL'
      },
    ],
    documentRoot: [
      {required: true, message: 'Document root is required'},
    ],
  }

  async function selectDocumentRoot() {
    const filePath = await OpenFolder("Select a folder")
    if (filePath) {
      props.proxy.targetParams.documentRoot = filePath
    }
  }
</script>
