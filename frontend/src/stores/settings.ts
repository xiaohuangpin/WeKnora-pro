import { defineStore } from "pinia";

// 定义设置接口
interface Settings {
  endpoint: string;
  apiKey: string;
  knowledgeBaseId: string;
}

// 默认设置
const defaultSettings: Settings = {
  endpoint: import.meta.env.VITE_IS_DOCKER ? "" : "http://localhost:8080",
  apiKey: "",
  knowledgeBaseId: "",
};

export const useSettingsStore = defineStore("settings", {
  state: () => ({
    // 从本地存储加载设置，如果没有则使用默认设置
    settings: JSON.parse(localStorage.getItem("WeKnora_settings") || JSON.stringify(defaultSettings)),
  }),

  actions: {
    // 保存设置
    saveSettings(settings: Settings) {
      this.settings = { ...settings };
      // 保存到localStorage
      localStorage.setItem("WeKnora_settings", JSON.stringify(this.settings));
    },

    // 获取设置
    getSettings(): Settings {
      return this.settings;
    },

    // 获取API端点
    getEndpoint(): string {
      return this.settings.endpoint || defaultSettings.endpoint;
    },

    // 获取API Key
    getApiKey(): string {
      return this.settings.apiKey;
    },

    // 获取知识库ID
    getKnowledgeBaseId(): string {
      return this.settings.knowledgeBaseId;
    },
  },
}); 