<template>
    <el-page-header class="header" @back="emit('back')">
        <template #content>
            <span class="header-title"> {{ title }} </span>
        </template>
        <template #extra>
            <div class="time">
                {{ date + second }}
            </div>
        </template>
    </el-page-header>
    <el-divider />
</template>

<script setup lang='ts'>
import { ref, onBeforeUnmount, Prop } from "vue";
import { getFormatTime, getSecond } from '@/hooks/useGetTime'

type Props = {
    title: string
}

defineProps<Props>()
const emit = defineEmits(['back'])
const date = ref()
const second = ref<string>()

date.value = getFormatTime()
second.value = getSecond()

const timer = setInterval(() => {
    second.value = getSecond()
    if (second.value === '00') {
        date.value = getFormatTime()
    }
}, 1000)

onBeforeUnmount(() => {
    clearInterval(timer)
})
</script>

<style lang='scss' scoped>
.header {
    padding: 24px 20px 0 20px;

    &-title {
        font-size: var(--font-size-4);
        font-weight: bold;
    }

    .time {
        font-size: var(--font-size-2);
    }
}
</style>