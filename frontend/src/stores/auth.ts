import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, TenantInfo, KnowledgeBaseInfo } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const user = ref<UserInfo | null>(null)
  const tenant = ref<TenantInfo | null>(null)
  const token = ref<string>('')
  const refreshToken = ref<string>('')
  const knowledgeBases = ref<KnowledgeBaseInfo[]>([])
  const currentKnowledgeBase = ref<KnowledgeBaseInfo | null>(null)

  // 计算属性
  const isLoggedIn = computed(() => {
    return !!token.value && !!user.value
  })

  const hasValidTenant = computed(() => {
    return !!tenant.value && !!tenant.value.api_key
  })

  const currentTenantId = computed(() => {
    return tenant.value?.id || ''
  })

  const currentUserId = computed(() => {
    return user.value?.id || ''
  })

  // 操作方法
  const setUser = (userData: UserInfo) => {
    user.value = userData
    // 保存到localStorage
    localStorage.setItem('weknora_user', JSON.stringify(userData))
  }

  const setTenant = (tenantData: TenantInfo) => {
    tenant.value = tenantData
    // 保存到localStorage
    localStorage.setItem('weknora_tenant', JSON.stringify(tenantData))
  }

  const setToken = (tokenValue: string) => {
    token.value = tokenValue
    localStorage.setItem('weknora_token', tokenValue)
  }

  const setRefreshToken = (refreshTokenValue: string) => {
    refreshToken.value = refreshTokenValue
    localStorage.setItem('weknora_refresh_token', refreshTokenValue)
  }

  const setKnowledgeBases = (kbList: KnowledgeBaseInfo[]) => {
    // 确保输入是数组
    knowledgeBases.value = Array.isArray(kbList) ? kbList : []
    localStorage.setItem('weknora_knowledge_bases', JSON.stringify(knowledgeBases.value))
  }

  const setCurrentKnowledgeBase = (kb: KnowledgeBaseInfo | null) => {
    currentKnowledgeBase.value = kb
    if (kb) {
      localStorage.setItem('weknora_current_kb', JSON.stringify(kb))
    } else {
      localStorage.removeItem('weknora_current_kb')
    }
  }


  const logout = () => {
    // 清空状态
    user.value = null
    tenant.value = null
    token.value = ''
    refreshToken.value = ''
    knowledgeBases.value = []
    currentKnowledgeBase.value = null

    // 清空localStorage
    localStorage.removeItem('weknora_user')
    localStorage.removeItem('weknora_tenant')
    localStorage.removeItem('weknora_token')
    localStorage.removeItem('weknora_refresh_token')
    localStorage.removeItem('weknora_knowledge_bases')
    localStorage.removeItem('weknora_current_kb')

  }

  const initFromStorage = () => {
    // 从localStorage恢复状态
    const storedUser = localStorage.getItem('weknora_user')
    const storedTenant = localStorage.getItem('weknora_tenant')
    const storedToken = localStorage.getItem('weknora_token')
    const storedRefreshToken = localStorage.getItem('weknora_refresh_token')
    const storedKnowledgeBases = localStorage.getItem('weknora_knowledge_bases')
    const storedCurrentKb = localStorage.getItem('weknora_current_kb')

    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
      } catch (e) {
        console.error('解析用户信息失败:', e)
      }
    }

    if (storedTenant) {
      try {
        tenant.value = JSON.parse(storedTenant)
      } catch (e) {
        console.error('解析租户信息失败:', e)
      }
    }

    if (storedToken) {
      token.value = storedToken
    }

    if (storedRefreshToken) {
      refreshToken.value = storedRefreshToken
    }

    if (storedKnowledgeBases) {
      try {
        const parsed = JSON.parse(storedKnowledgeBases)
        knowledgeBases.value = Array.isArray(parsed) ? parsed : []
      } catch (e) {
        console.error('解析知识库列表失败:', e)
        knowledgeBases.value = []
      }
    }

    if (storedCurrentKb) {
      try {
        currentKnowledgeBase.value = JSON.parse(storedCurrentKb)
      } catch (e) {
        console.error('解析当前知识库失败:', e)
      }
    }
  }

  // 初始化时从localStorage恢复状态
  initFromStorage()

  return {
    // 状态
    user,
    tenant,
    token,
    refreshToken,
    knowledgeBases,
    currentKnowledgeBase,
    
    // 计算属性
    isLoggedIn,
    hasValidTenant,
    currentTenantId,
    currentUserId,
    
    // 方法
    setUser,
    setTenant,
    setToken,
    setRefreshToken,
    setKnowledgeBases,
    setCurrentKnowledgeBase,
    logout,
    initFromStorage
  }
})