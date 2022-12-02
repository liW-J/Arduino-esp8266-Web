<template>
    <div class="content">
        <Fingerprint title="点击开始验证指纹" @click="confirm" />
    </div>
</template>
  
<script setup lang='ts'>
import { useRouter } from "vue-router";
import Fingerprint from "@/components/Fingerprint.vue";
import { userLogin } from '@/request/apis/client'
import { ElLoading, ElMessage } from 'element-plus'

const router = useRouter();

const confirm = () => {
    const loading = ElLoading.service({
        lock: true,
        text: '验证中',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    userLogin().then(token => {
        localStorage.setItem('token', token)
        router.replace('/home')
    }, err => {
        ElMessage.error(err)

    }).finally(() => {
        loading.close()
    })
}
</script>
  
<style lang='scss' scoped>
.content {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>