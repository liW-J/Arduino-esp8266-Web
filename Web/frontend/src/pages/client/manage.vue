<template>
    <HeaderBar title="指纹管理" @back="router.back()" />
    <div class="tabs" v-if="!show">
        <TabItem class="tabs-add" title="添加指纹" @click="showScanFinger('add')" />
        <TabItem class="tabs-update" title="更新指纹" @click="showScanFinger('update')" />
    </div>
    <div v-else class="finger">
        <Fingerprint title="点击开始录制指纹" @click="handelEdit"></Fingerprint>
        <h3 class="cancel" @click="show = false">取消</h3>
    </div>

    <el-dialog style="width:420px" v-model="dialogFormVisible" title="请选择要更新的指纹">
        <el-form ref="ruleFormRef" style="width:360px" :model="form" label-width="120px">
            <el-form-item label="指纹编号：">
                <el-radio-group v-model="form.fingerNum">
                    <el-radio :label="1" size="large" border>指纹1</el-radio>
                    <el-radio :label="2" size="large" border>指纹2</el-radio>
                </el-radio-group>
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogFormVisible = false">取消</el-button>
                <el-button type="primary" @click.prevent="handelUpdate">更新</el-button>
            </span>
        </template>
    </el-dialog>
</template>
  
<script setup lang='ts'>
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import HeaderBar from '@/components/HeaderBar.vue'
import Fingerprint from '@/components/Fingerprint.vue'
import TabItem from '@/components/TabItem.vue'
import { ElLoading, ElMessage, FormInstance } from 'element-plus';
import { addFinger, updateFinger } from '@/request/apis/client'

const router = useRouter()
const show = ref(false)
const type = ref('add')
let dialogFormVisible = ref(false)
const ruleFormRef = ref<FormInstance>()
let form = reactive({
    fingerNum: 1
})

const showScanFinger = (t: string) => {
    show.value = true
    type.value = t
}

const showLoading = () => {
    return ElLoading.service({
        lock: true,
        text: '扫描中',
        background: 'rgba(0, 0, 0, 0.7)',
    })
}

const handelAdd = () => {
    const loading = showLoading()
    addFinger().then(msg => {
        ElMessage.success(msg)
    }, err => {
        ElMessage.error(err)
    }).finally(() => {
        loading.close()
        show.value = false
    })
}

const handelUpdate = () => {
    const loading = showLoading()
    updateFinger(form.fingerNum).then(msg => {
        ElMessage.success(msg)
    }, err => {
        ElMessage.error(err)
    }).finally(() => {
        loading.close()
    })
}

const handelEdit = () => {
    if (type.value === 'update') {
        dialogFormVisible.value = true
    } else {
        handelAdd()
    }
}
</script>
  
<style lang='scss' scoped>
.tabs {
    height: 80vh;
    width: 100vw;
    display: flex;
    justify-content: space-around;
    align-items: center;
    padding: 0 10vw;
    flex-wrap: wrap;

    &-add {
        background: linear-gradient(to left top, #faffd1, #a1ffce);
    }

    &-update {
        background: var(--gradient-15);
    }
}

.finger {
    height: 80vh;
    display: flex;
    align-items: center;
    flex-direction: column;
    justify-content: center;

    .cancel {
        text-align: center;
        margin: 0 auto;
        margin-top: 2vh;
        cursor: pointer;
        transition: all;

        &:hover {
            color: var(--blue-6);
        }
    }
}
</style>