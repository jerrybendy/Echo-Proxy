<template>
  <div></div>
</template>

<script setup lang="ts">
import {EventsOn} from "../../wailsjs/runtime";
import {useNotification} from "naive-ui"

const notification = useNotification()

EventsOn('error', function (errorMessage) {
  console.error(errorMessage);
  notification.error({
    content: 'Error',
    meta: errorMessage,
    duration: 3000,
    keepAliveOnHover: true,
  })
})

EventsOn('warning', function (errorMessage) {
  console.warn(errorMessage);
  notification.warning({
    content: 'Warning',
    meta: errorMessage,
    duration: 3000,
    keepAliveOnHover: true,
  })
})

</script>
