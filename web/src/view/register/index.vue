<template>
  <div id="userLayout" class="w-full h-full bg-[#f7f9fc]">
    <div class="flex items-center justify-center min-h-screen">
      <div class="w-[400px] bg-white rounded-lg shadow-sm p-10">
        <div>
          <div class="flex items-center justify-center mb-8">
            <img class="w-16" :src="$GIN_VUE_ADMIN.appLogo" alt="Logo" />
          </div>
          <div class="mb-8">
            <p class="text-center text-2xl font-medium text-[#1f2329]">
              注册新用户
            </p>
            <p class="text-center text-sm text-[#646a73] mt-1">
              创建您的账号
            </p>
          </div>
          <el-form
            ref="registerForm"
            :model="registerFormData"
            :rules="rules"
            @keyup.enter="submitForm"
          >
            <el-form-item prop="username" class="mb-6">
              <el-input
                v-model="registerFormData.username"
                size="large"
                class="h-10 rounded-lg"
                placeholder="请输入用户名"
                suffix-icon="user"
              />
            </el-form-item>
            <el-form-item prop="password" class="mb-6">
              <el-input
                v-model="registerFormData.password"
                show-password
                size="large"
                type="password"
                class="h-10 rounded-lg"
                placeholder="请输入密码"
              />
            </el-form-item>
            <el-form-item prop="confirmPassword" class="mb-6">
              <el-input
                v-model="registerFormData.confirmPassword"
                show-password
                size="large"
                type="password"
                class="h-10 rounded-lg"
                placeholder="请确认密码"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                class="w-full h-10 bg-[#3370ff] hover:bg-[#2860df] border-none rounded-lg"
                type="primary"
                size="large"
                @click="submitForm"
              >注册</el-button>
            </el-form-item>
            <div class="text-center mt-4">
              <span class="text-[#646a73] text-sm">已有账号？</span>
              <el-button
                class="text-[#3370ff] hover:text-[#2860df] ml-1"
                type="text"
                @click="goLogin"
              >立即登录</el-button>
            </div>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { register } from '@/api/user'

defineOptions({
  name: 'Register'
})

const router = useRouter()

// 验证函数
const checkUsername = (rule, value, callback) => {
  if (value.length < 5) {
    return callback(new Error('用户名不能少于5个字符'))
  } else {
    callback()
  }
}

const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('密码不能少于6个字符'))
  } else {
    callback()
  }
}

const checkConfirmPassword = (rule, value, callback) => {
  if (value !== registerFormData.password) {
    return callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const registerForm = ref(null)
const registerFormData = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  confirmPassword: [{ validator: checkConfirmPassword, trigger: 'blur' }]
})

const submitForm = () => {
  registerForm.value.validate(async (valid) => {
    if (!valid) {
      ElMessage({
        type: 'error',
        message: '请正确填写注册信息',
        showClose: true
      })
      return false
    }

    try {
      const { username, password } = registerFormData
      await register({ username, password })
      ElMessage({
        type: 'success',
        message: '注册成功！',
        showClose: true
      })
      router.push({ name: 'Login' })
    } catch (error) {
      ElMessage({
        type: 'error',
        message: error.response?.data?.msg || '注册失败，请重试',
        showClose: true
      })
    }
  })
}

const goLogin = () => {
  router.push({ name: 'Login' })
}
</script>
