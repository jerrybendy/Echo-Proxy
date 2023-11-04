<template>
  <NForm
      :model="model"
      :rules="rules"
      label-placement="left"
      label-width="auto"
      require-mark-placement="right-hanging"
      size="small"
  >
    <NCard title="Basic" size="small" :bordered="false" header-style="--n-title-font-size: 14px">
      <NFormItem label="Name" path="name">
        <NInput v-model="model.name" placeholder="Domain name"/>
      </NFormItem>
      <NFormItem path="applyHosts">
        <NCheckbox v-model="model.applyHosts">Apply to /etc/hosts file</NCheckbox>
      </NFormItem>

    </NCard>
  </NForm>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import {FormRules, NCard, NCheckbox, NForm, NFormItem, NInput} from 'naive-ui'

const model = ref<HostConfig>({
  name: '',
  applyHosts: false,
})

const rules = ref<FormRules>({
  name: {required: true}
})

</script>

<style scoped>

</style>