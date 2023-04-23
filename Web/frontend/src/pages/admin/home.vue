<template>
    <el-tabs type="border-card" @tab-click="getData">
        <el-tab-pane class="content" label="体温日志">
            <el-input class="filterInput" v-model="search.temp" placeholder="请输入 工号/姓名" />
            <el-table :data="filterTempLogs" border max-height="80vh" style="width: 100%"
                :row-class-name="tableRowClassName">
                <el-table-column prop="ID" label="ID" sortable width="120" />
                <el-table-column prop="workNum" label="工号" width="240" />
                <el-table-column prop="userName" label="用户名" width="180" />
                <el-table-column prop="fingerprintData" :filters="[
                    { text: '体温正常', value: 'normal' },
                    { text: '体温低', value: 'lower' },
                    { text: '体温异常', value: 'abnormal' },
                ]" :filter-method="filterTemperature" label="体温" width="180" />
                <el-table-column prop="CreatedAt" sortable label="测量时间" />
<!--                <el-table-column prop="UpdatedAt" sortable label="更新时间" />-->
            </el-table>
        </el-tab-pane>
        <el-tab-pane class="content" label="用户管理">
            <el-button class="addUserBtn" type="primary" plain @click="dialogFormVisible = true">添加用户</el-button>
            <el-table :data="filterUserList" border max-height="80vh" stripe style="width: 100%">
                <el-table-column prop="ID" label="ID" sortable width="120" />
                <el-table-column prop="workNum" label="工号" width="180" />
                <el-table-column prop="userName" label="姓名" width="180" />
                <el-table-column prop="finger1" label="指纹1" width="120" />
                <el-table-column prop="finger2" label="指纹2" width="120" />
                <el-table-column prop="CreatedAt" sortable label="创建时间" />
                <el-table-column prop="UpdatedAt" sortable label="更新时间" />
                <el-table-column align="right">
                    <template #header>
                        <el-input v-model="search.user" size="small" placeholder="请输入 工号/姓名" />
                    </template>
                    <template #default="scope">
                        <el-popconfirm title="确认删除?" @confirm="handleDelete(scope.$index, scope.row)">
                            <template #reference>
                                <el-button size="small" type="danger">删除</el-button>
                            </template>
                        </el-popconfirm>
                    </template>
                </el-table-column>
            </el-table>
        </el-tab-pane>
        <el-tab-pane class="content" style="padding-top:0" label="系统监控">
            <el-tabs v-model="sysActiveName" @tab-click="getSysData" style="height: fit-content;">
                <el-tab-pane label="操作日志" name="operate">
                    <el-input class="filterInput" v-model="search.operate" placeholder="请输入 操作人/操作对象" />
                    <el-table :data="filterOperateLogs" stripe border max-height="70vh" style="width: 100%">
                        <el-table-column prop="ID" label="ID" sortable width="120" />
                        <el-table-column prop="role" label="操作人" width="150" />
                        <el-table-column prop="operate" label="操作类型" width="180" />
                        <el-table-column prop="object" label="操作对象" />
                        <el-table-column prop="CreatedAt" sortable label="操作时间" />
<!--                        <el-table-column prop="UpdatedAt" sortable label="更新时间" />-->
                    </el-table>
                </el-tab-pane>
                <el-tab-pane label="指纹识别日志" name="finger">
                    <el-input class="filterInput" v-model="search.finger" placeholder="请输入 工号/指纹ID" />
                    <el-table :data="filterFingerLogs" stripe border max-height="70vh" style="width: 100%">
                        <el-table-column prop="ID" label="ID" sortable width="120" />
                        <el-table-column prop="workNum" label="工号" width="180" />
                        <el-table-column prop="fingerId" label="指纹ID" width="180" />
                        <el-table-column prop="status" sortable label="验证状态" />
                        <el-table-column prop="CreatedAt" sortable label="创建时间" />
<!--                        <el-table-column prop="UpdatedAt" sortable label="更新时间" />-->
                    </el-table>
                </el-tab-pane>
            </el-tabs>
        </el-tab-pane>
    </el-tabs>
    <el-button class="quit-btn" type="primary" @click="quit">退出</el-button>

    <el-dialog style="width:420px" v-model="dialogFormVisible" title="请输入用户信息">
        <el-form ref="ruleFormRef" style="width:360px" :model="form" :rules="form.rules" label-width="100px">
            <el-form-item label="用户名：" prop="userName">
                <el-input v-model="form.userName" placeholder="请输入用户名" autocomplete="off" />
            </el-form-item>
            <el-form-item label="工号：" label-width="100px" prop="workNum">
                <el-input v-model="form.workNum" placeholder="请输入工号" />
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogFormVisible = false">取消</el-button>
                <el-button type="primary" @click.prevent="addUser(ruleFormRef)">添加</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<script setup lang='ts'>
import { reactive, ref, watchEffect } from 'vue'
import { useRouter } from 'vue-router';
import { getTempLogList, getUserList, userSignUp, deleteUser, getFingerLogList, getOperateLogList } from '@/request/apis/admin'
import type { tempLog, userInfo, fingerLog, operateLog } from '@/type/admin'
import { useFormateDate } from '@/hooks/useFormateDate';
import { ElLoading, ElMessage, FormInstance } from 'element-plus';

const router = useRouter()
let tempLogs = ref<Array<tempLog>>()
let filterTempLogs = ref<Array<tempLog>>()
let userList = ref<Array<userInfo>>()
let filterUserList = ref<Array<userInfo>>()
let fingerLogs = ref<Array<fingerLog>>()
let filterFingerLogs = ref<Array<fingerLog>>()
let operateLogs = ref<Array<operateLog>>()
let filterOperateLogs = ref<Array<operateLog>>()
let dialogFormVisible = ref(false)
let sysActiveName = ref('operate')

const ruleFormRef = ref<FormInstance>()
let form = reactive({
    userName: '',
    workNum: '',
    rules: {
        userName: [
            { required: true, message: '请输入用户名', trigger: 'blur' },
        ],
        workNum: [
            { required: true, message: '请输入密码', trigger: 'blur' },
        ]
    }
})
const search = reactive({
    temp: '',
    user: '',
    finger: '',
    operate: ''
})


const quit = () => {
    localStorage.removeItem('adminToken')
    router.replace('/admin')
}

const showLoading = (text = '加载中') => {
    return ElLoading.service({
        lock: true,
        text: text,
        background: 'rgba(0, 0, 0, 0.7)',
    })
}

const tableRowClassName = ({
    row,
}: {
    row: tempLog
}) => {
    if (+row.fingerprintData < 36) {
        return 'warning-row'
    } else if (+row.fingerprintData > 37.2) {
        return 'danger-row'
    } else {
        return 'success-row'
    }
}

const getTempLogs = () => {
    const loading = showLoading()
    getTempLogList().then(data => {
        tempLogs.value = data
        for (let log of tempLogs.value) {
            log.CreatedAt = useFormateDate(new Date(log.CreatedAt))
            if (log.UpdatedAt) {
                log.UpdatedAt = useFormateDate(new Date(log.UpdatedAt))
            } else {
                log.UpdatedAt = '暂无更新'
            }
        }
        loading.close()
    })
}
getTempLogs()

const getUsers = () => {
    const loading = showLoading()
    getUserList().then(data => {
        userList.value = data
        for (let user of userList.value) {
            user.CreatedAt = useFormateDate(new Date(user.CreatedAt))
            if (user.UpdatedAt) {
                user.UpdatedAt = useFormateDate(new Date(user.UpdatedAt))
            } else {
                user.UpdatedAt = '暂无更新'
            }
            if (!user.finger1) user.finger1 = '暂无指纹'
            if (!user.finger2) user.finger2 = '暂无指纹'
        }
        loading.close()
    })
}

const getFingerLogs = () => {
    const loading = showLoading()
    getFingerLogList().then(data => {
        fingerLogs.value = data
        for (let log of fingerLogs.value) {
            log.CreatedAt = useFormateDate(new Date(log.CreatedAt))
            if (log.UpdatedAt) {
                log.UpdatedAt = useFormateDate(new Date(log.UpdatedAt))
            } else {
                log.UpdatedAt = '暂无更新'
            }
        }
        loading.close()
    })
}

const getOperateLogs = () => {
    const loading = showLoading()
    getOperateLogList().then(data => {
        operateLogs.value = data
        for (let log of operateLogs.value) {
            log.CreatedAt = useFormateDate(new Date(log.CreatedAt))
            if (log.UpdatedAt) {
                log.UpdatedAt = useFormateDate(new Date(log.UpdatedAt))
            } else {
                log.UpdatedAt = '暂无更新'
            }
        }
        loading.close()
    })
}


const getData = (pane: any) => {
    const label = pane.props.label
    if (label === '体温日志') {
        getTempLogs()
    } else if (label === '用户管理') {
        getUsers()
    } else if (label === '系统监控') {
        getOperateLogs()
    }
}

const getSysData = (pane: any) => {
    const label = pane.props.label
    if (label === '操作日志') {
        getOperateLogs()
    } else if (label === '指纹识别日志') {
        getFingerLogs()
    }
}

watchEffect(() => {
    filterTempLogs.value = tempLogs.value?.filter((data) =>
        !search.temp ||
        data.userName.toLowerCase().includes(search.temp.toLowerCase()) ||
        data.workNum.includes(search.temp)
    )
})


watchEffect(() => {
    filterUserList.value = userList.value?.filter((data) =>
        !search.user ||
        data.userName.toLowerCase().includes(search.user.toLowerCase()) ||
        data.workNum.includes(search.user)
    )
})

watchEffect(() => {
    filterFingerLogs.value = fingerLogs.value?.filter((data) =>
        !search.finger ||
        data.workNum.includes(search.finger) ||
        data.fingerId.includes(search.finger)
    )
})

watchEffect(() => {
    filterOperateLogs.value = operateLogs.value?.filter((data) =>
        !search.operate ||
        data.role.toLowerCase().includes(search.operate.toLowerCase()) ||
        data.object.toLowerCase().includes(search.operate.toLowerCase())
    )
})

const addUser = async (formEl: FormInstance | undefined) => {
    if (!formEl) return

    await formEl.validate((valid) => {
        if (valid) {
            const loading = showLoading('等待验证指纹')
            userSignUp(form.userName, form.workNum).then(res => {
                dialogFormVisible.value = false
                getUsers()
                ElMessage.success(res)
            }, err => {
                ElMessage.error(err)
            }).finally(() => {
                loading.close()
            })
        }
    })
}

const handleDelete = (index: number, row: userInfo) => {
    const loading = showLoading('删除中')
    deleteUser(row.workNum).then(res => {
        filterUserList.value?.splice(index, 1)
        ElMessage.success(res)
    }, err => {
        ElMessage.error(err)
    }).finally(() => {
        loading.close()
    })
}

const filterTemperature = (
    value: string,
    row: tempLog,
) => {
    if (value === 'lower') return +row.fingerprintData < 36
    if (value === 'normal') return (+row.fingerprintData >= 36 && +row.fingerprintData <= 37.2)
    if (value === 'abnormal') return +row.fingerprintData > 37.2
}
</script>

<style lang='scss' scoped>
:deep(.el-tabs__item) {
    height: 55px;
    line-height: 55px;
}

.quit-btn {
    position: absolute;
    height: 34px;
    margin: 13px 20px;
    right: 0;
    top: 0;
}

.content {
    flex: 1;
    padding: 3vh 5vw;
}

.filterInput,
.addUserBtn {
    margin-bottom: 15px;
    width: 200px;
}
</style>

<style>
.el-tabs {
    display: flex;
    flex-direction: column;
    height: 100vh;
}

.el-tabs__content {
    overflow: auto;
}

.el-table .warning-row {
    --el-table-tr-bg-color: var(--el-color-warning-light-9);
}

.el-table .success-row {
    --el-table-tr-bg-color: var(--el-color-success-light-8);
}

.el-table .danger-row {
    --el-table-tr-bg-color: var(--el-color-danger-light-7);
}
</style>