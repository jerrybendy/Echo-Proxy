<template>
  <NDrawerContent title="Edit" closable :native-scrollbar="false">
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
    </NForm>

    <template #footer>
      <NButton type="primary" secondary @click="saveSetting">Save</NButton>
    </template>
  </NDrawerContent>
</template>

<script setup lang="ts">
import {ref, watchEffect, PropType, readonly} from 'vue'
import {FormRules, NCard, NCheckbox, NForm, NFormItem, NInput, NButton, NDrawerContent, FormInst} from 'naive-ui'
import {userData} from "../../wailsjs/go/models";
import {SaveSetting} from "../../wailsjs/go/userData/Hosts";

const props = defineProps({
  host: {
    type: Object as PropType<userData.HostConfig|null>,
    required: true,
  },
})

const emit = defineEmits<{
  change: [],
}>()

const formRef = ref<FormInst | null>(null)
const defaultConfig = readonly<userData.HostConfig>({
  id: 0,
  name: '',
  applyHosts: false,
})
const model = ref<userData.HostConfig>(defaultConfig)

watchEffect(() => {
  if (props.host) {
    model.value = {...props.host}
  } else {
    model.value = {...defaultConfig}
  }
})

const rules = ref<FormRules>({
  name: [
    {required: true, message: 'Domain name is required'},
    {pattern: /^[-a-z0-9]+(\.[-a-z0-9]+)*\.[a-zA-Z]{2,}$/i, message: 'Invalid domain name'},
  ]
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

</script>

<style scoped>

</style>