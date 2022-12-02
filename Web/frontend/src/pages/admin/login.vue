<template>
    <div class="container">
        <div class="screen">
            <div class="screen__content">
                <el-form ref="ruleFormRef" :model="Form" :rules="rules" class="login">
                    <el-form-item class="login__field" prop="username">
                        <div class="login__field">
                            <el-icon>
                                <User />
                            </el-icon>
                            <input type="text" class="login__input" placeholder="请输入用户名" v-model="Form.username">
                        </div>
                    </el-form-item>
                    <el-form-item class="login__field" prop="password">
                        <div class="login__field">
                            <el-icon>
                                <Lock />
                            </el-icon>
                            <input type="password" class="login__input" placeholder="请输入密码" v-model="Form.password"
                                autocomplete="true">
                        </div>
                    </el-form-item>

                    <el-form-item>
                        <button class="login__submit" @click.prevent="submitForm(ruleFormRef)"><span
                                class="button__text">登录</span>
                            <el-icon>
                                <Position />
                            </el-icon>
                        </button>
                    </el-form-item>
                </el-form>
            </div>
            <div class="screen__background">
                <span class="screen__background__shape screen__background__shape4"></span>
                <span class="screen__background__shape screen__background__shape3"></span>
                <span class="screen__background__shape screen__background__shape2"></span>
                <span class="screen__background__shape screen__background__shape1"></span>
            </div>
        </div>
    </div>
</template>

<script setup lang='ts'>
import { reactive, ref } from 'vue'
import { Position, User, Lock } from '@element-plus/icons-vue'
import { ElLoading, ElMessage, FormInstance, FormRules } from 'element-plus'
import { useRouter } from 'vue-router'
import { login } from '@/request/apis/admin'

const router = useRouter()
const ruleFormRef = ref<FormInstance>()
const Form = reactive({
    username: '',
    password: ''
})

const rules = reactive<FormRules>({
    username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
    ]
})

const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid) => {
        if (valid) {
            const loading = ElLoading.service({
                lock: true,
                text: '登录中',
                background: 'rgba(0, 0, 0, 0.7)',
            })
            login(Form.username, Form.password).then(token => {
                localStorage.setItem('adminToken', token)
                router.replace('/admin/home')
            }, err => {
                ElMessage.error(err)
            }).finally(() => {
                loading.close()
            })
        }
    })
}

</script>

<style lang='scss' scoped>
.container {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100vh;
    background: linear-gradient(90deg, #f3fdff, #c9f2ff);
    overflow: hidden;
}

.screen {
    background: linear-gradient(90deg, #aceaff, #8fe5ff);
    position: relative;
    height: 460px;
    width: 324px;
    box-shadow: 0px 0px 24px #9fe9ff;
    transition: all .3s;

    &:hover {
        scale: 1.05;
        box-shadow: 0px 0px 50px #71deff;
    }

    &__content {
        z-index: 1;
        position: relative;
        height: 100%;
    }

    &__background {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        z-index: 0;
        -webkit-clip-path: inset(0 0 0 0);
        clip-path: inset(0 0 0 0);

        &__shape {
            transform: rotate(45deg);
            position: absolute;
        }

        &__shape1 {
            height: 468px;
            width: 468px;
            background: #FFF;
            top: -45px;
            right: 108px;
            border-radius: 0 72px 0 0;
        }

        &__shape2 {
            height: 198px;
            width: 198px;
            background: #00c3ff;
            top: -154.8px;
            right: 0;
            border-radius: 32px;
        }

        &__shape3 {
            height: 486px;
            width: 171px;
            background: linear-gradient(270deg, #60ddff, #00c3ff);
            top: -21.6px;
            right: 0;
            border-radius: 32px;
        }

        &__shape4 {
            height: 360px;
            width: 180px;
            background: #60ddff;
            top: 378px;
            right: 45px;
            border-radius: 60px;
        }
    }

}

.login {
    width: 280px;
    padding: 30px;
    padding-top: 80px;

    &__field {
        padding: 10px 0px;
        position: relative;
    }

    &__input {
        border: none;
        background: none;
        border-bottom: 2px solid #D1D1D4;
        padding: 5px;
        padding-left: 10px;
        font-weight: 700;
        width: 85%;
        transition: .2s;

        &:active,
        &:focus,
        &:hover {
            outline: none;
            border-bottom-color: #6c99cd;

            &::-webkit-input-placeholder {
                color: #79bbff;
                opacity: .5;
            }
        }

        &:active,
        &:hover {
            filter: drop-shadow(0 0 20px #79bbff) drop-shadow(0 0 20px #79bbff);
        }

        &:focus {
            filter: drop-shadow(0 0 20px #abc6e2);
        }
    }


    &__submit {
        background: #fff;
        font-size: 20px;
        margin-top: 20px;
        padding: 8px 16px;
        border-radius: 20px;
        border: 1px solid #b5daff;
        font-weight: 700;
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
        height: 45px;
        color: #79bbff;
        box-shadow: 0px 0px 2px #3073ee;
        cursor: pointer;
        transition: .2s;

        &:hover {
            border-color: #67b3ff;
            outline: none;
            transform: scale(0.98);
        }
    }
}
</style>