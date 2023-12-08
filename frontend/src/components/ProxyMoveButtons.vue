<template>
  <NTooltip content-style="padding: 4px 6px; background: #262626; border-radius: 4px" raw>
    <template #trigger>
      <NButton quaternary size="tiny" :disabled="index === 0" @click="emit('moveUp')">
        <template #icon><NIcon><ArrowUp /></NIcon></template>
      </NButton>
    </template>
    <span style="font-size: 12px">Move up</span>
  </NTooltip>
  <NTooltip content-style="padding: 4px 6px; background: #262626; border-radius: 4px" raw>
    <template #trigger>
      <NButton quaternary size="tiny" :disabled="index === count - 1" @click="emit('moveDown')">
        <template #icon><NIcon><ArrowDown /></NIcon></template>
      </NButton>
    </template>
    <span style="font-size: 12px">Move down</span>
  </NTooltip>
  <NPopconfirm
      :negative-button-props="{type: 'error', ghost: true} as ButtonProps"
      :positive-button-props="{type: 'error'} as ButtonProps"
      @positiveClick="emit('remove')"
  >
    <template #trigger>
      <NTooltip content-style="padding: 4px 6px; background: #262626; border-radius: 4px" raw>
        <template #trigger>
          <NButton quaternary size="tiny" :disabled="count === 1" type="error">
            <template #icon><NIcon><Close /></NIcon></template>
          </NButton>
        </template>
        <span style="font-size: 12px">Delete</span>
      </NTooltip>
    </template>
    Do you confirm to remove this proxy setting?
  </NPopconfirm>
</template>

<script setup lang="ts">
  import {ArrowDown, ArrowUp, Close} from "@vicons/ionicons5";
  import {ButtonProps, NButton, NIcon, NPopconfirm, NTooltip} from "naive-ui";

  const props = defineProps<{
    index: number,
    count: number,
  }>()

  const emit = defineEmits<{
    moveUp: [],
    moveDown: [],
    remove: [],
  }>()

</script>
