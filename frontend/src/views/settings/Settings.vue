<template>
    <div class="settings-container">
        <div class="settings-header">
            <h2>系统配置</h2>
        </div>
        <div class="settings-form">
            <t-form ref="form" :data="formData" :rules="rules" @submit="onSubmit">
                <t-form-item label="API 服务端点" name="endpoint">
                    <t-input v-model="formData.endpoint" placeholder="请输入API服务端点，例如：http://localhost" />
                </t-form-item>
                <t-form-item label="API Key" name="apiKey">
                    <t-input v-model="formData.apiKey" placeholder="请输入API Key" />
                </t-form-item>
                <t-form-item label="知识库ID" name="knowledgeBaseId">
                    <t-input v-model="formData.knowledgeBaseId" placeholder="请输入知识库ID" />
                </t-form-item>
                <t-form-item>
                    <t-space>
                        <t-button theme="primary" type="submit">保存配置</t-button>
                        <t-button theme="default" @click="resetForm">重置</t-button>
                    </t-space>
                </t-form-item>
            </t-form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { useSettingsStore } from '@/stores/settings';

const settingsStore = useSettingsStore();
const form = ref(null);

const formData = reactive({
    endpoint: '',
    apiKey: '',
    knowledgeBaseId: ''
});

const rules = {
    endpoint: [{ required: true, message: '请输入API服务端点', trigger: 'blur' }],
    apiKey: [{ required: true, message: '请输入API Key', trigger: 'blur' }],
    knowledgeBaseId: [{ required: true, message: '请输入知识库ID', trigger: 'blur' }]
};

onMounted(() => {
    // 初始化表单数据
    const settings = settingsStore.getSettings();
    formData.endpoint = settings.endpoint;
    formData.apiKey = settings.apiKey;
    formData.knowledgeBaseId = settings.knowledgeBaseId;
});

const onSubmit = ({ validateResult }) => {
    if (validateResult === true) {
        settingsStore.saveSettings({
            endpoint: formData.endpoint,
            apiKey: formData.apiKey,
            knowledgeBaseId: formData.knowledgeBaseId
        });
        MessagePlugin.success('配置保存成功');
    }
};

const resetForm = () => {
    const settings = settingsStore.getSettings();
    formData.endpoint = settings.endpoint;
    formData.apiKey = settings.apiKey;
    formData.knowledgeBaseId = settings.knowledgeBaseId;
};
</script>

<style lang="less" scoped>
.settings-container {
    padding: 20px;
    background-color: #fff;
    border-radius: 8px;
    margin: 20px;
    min-height: 80vh;

    .settings-header {
        margin-bottom: 20px;
        border-bottom: 1px solid #f0f0f0;
        padding-bottom: 16px;

        h2 {
            font-size: 20px;
            font-weight: 600;
            color: #000000;
            margin: 0;
        }
    }

    .settings-form {
        max-width: 600px;
    }
}
</style> 