<template>
    <div class="header">{{ date + second }}</div>
    <el-button class="quit-btn" type="primary" @click="quit">退出</el-button>
    <div class="tabs">
        <TabItem class="tabs-clockIn" title="体温打卡" @click="clockIn"></TabItem>
        <TabItem class="tabs-manage" title="指纹管理" @click="router.push('/manage')"></TabItem>
    </div>
</template>

<script setup lang='ts'>
import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import TabItem from '@/components/TabItem.vue'
import { getFormatTime, getSecond } from '@/hooks/useGetTime'
import { ElLoading, ElMessage } from 'element-plus'
import { recordTemp } from '@/request/apis/client'

const router = useRouter()
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

const quit = () => {
    localStorage.removeItem('token')
    router.replace('/')
}

onBeforeUnmount(() => {
    clearInterval(timer)
})

const clockIn = () => {
    const loading = ElLoading.service({
        lock: true,
        text: '测温中',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    recordTemp().then(temperature => {
        ElMessage.success(temperature+"℃\n测温成功，体温正常！")
    }, err => {
        ElMessage.error(err)
    }).finally(() => {
        loading.close()
    })
}
</script>

<style lang='scss' scoped>
.header {
    height: 30vh;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 7vh;
    font-weight: bold;
}

.tabs {
    width: 100vw;
    display: flex;
    justify-content: space-around;
    padding: 0 10vw;
    flex-wrap: wrap;


    &-clockIn {
        background: var(--gradient-9);
    }

    &-manage {
        background: var(--gradient-28);
    }
}

.quit-btn {
    position: absolute;
    height: 34px;
    margin: 13px 20px;
    right: 0;
    top: 0;
}
</style>