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
          <NSelect v-model:value="proxy.matchType" :options="targetType" style="width: 140px" />
          <NInput v-model:value="proxy.matchRule" />
        </NInputGroup>
<!--        <NCheckbox v-if="proxy.matchType === 'REGEXP'" v-model:checked="proxy.matchParams.caseInsensitive" style="margin: 8px 0 0 4px;">Case-insensitive</NCheckbox>-->
      </div>
    </NFormItem>

    <template #footer>
      <h3>Rule</h3>
      <NFormItem label="Target" :path="`proxies.${index}.target`" :show-require-mark="false" :rule="formRules.target">
        <NInput v-model:value="proxy.target" placeholder="http://127.0.0.1:8080"/>
      </NFormItem>
      <NSpace vertical style="padding-left: 4px">
        <NCheckbox v-model:checked="proxy.changeOrigin">Change origin</NCheckbox>
        <!--              <NCheckbox v-model:checked="model.applyHosts">X-Forwarded-*</NCheckbox>-->
        <!--              <NCheckbox v-model:checked="model.applyHosts">Change cookie domain</NCheckbox>-->
      </NSpace>

    </template>
  </NCard>
</template>

<script setup lang="ts">
  import {ref} from "vue"
  import {NCard, NCheckbox, NFormItem, NInput, NInputGroup, NSelect, NSpace, FormItemRule} from "naive-ui";
  import ProxyMoveButtons from "./ProxyMoveButtons.vue";
  import {service} from "../../wailsjs/go/models";

  defineProps<{
    proxy: service.HostProxy,
    index: number,
    count: number,
  }>()

  const emit = defineEmits<{
    moveUp: [],
    moveDown: [],
    remove: [],
  }>()

  const targetType = ref([
    {label: 'PREFIX', value: 'PREFIX'},
    // {label: 'REGEXP', value: 'REGEXP'},
    // {label: 'GLOB', value: 'GLOB'},
    // {label: 'STATIC', value: 'STATIC'},
  ])

  const formRules: { [key: string]: FormItemRule | FormItemRule[] } = {
    matchRule: {required: true, message: 'Match rule is required'},
    target: [
      {required: true, message: 'Target is required'},
      {
        pattern: /^https?:\/\/(([a-zA-Z0-9_-])+(\.)?)*(:\d+)?(\/((\.)?(\?)?=?&?[a-zA-Z0-9_-](\?)?)*)*$/i,
        message: 'Target must be a valid URL'
      },
    ],
  }
</script>
