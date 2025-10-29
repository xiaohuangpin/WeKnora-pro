import { ref, computed, reactive } from 'vue'

import { defineStore } from 'pinia';

export const useMenuStore = defineStore('menuStore', {
    state: () => ({
        menuArr: reactive([
            { title: '知识库', icon: 'zhishiku', path: 'knowledge-bases' },
            {
                title: '对话',
                icon: 'prefixIcon',
                path: 'creatChat',
                childrenPath: 'chat',
                children: reactive<object[]>([]),
            },
            { title: '系统信息', icon: 'tenant', path: 'tenant' },
            { title: '退出登录', icon: 'logout', path: 'logout' }
        ]),
        isFirstSession: false,
        firstQuery: ''
    }
    ),
    actions: {
        clearMenuArr() {
            this.menuArr[1].children = reactive<object[]>([]);
        },
        updatemenuArr(obj: any) {
            this.menuArr[1].children?.push(obj);
        },
        updataMenuChildren(item: object) {
            this.menuArr[1].children?.unshift(item)
        },
        updatasessionTitle(session_id: string, title: string) {
            this.menuArr[1].children?.forEach((item: any) => {
                if (item.id == session_id) {
                    item.title = title;
                    item.isNoTitle = false;
                }
            });
        },
        changeIsFirstSession(payload: boolean) {
            this.isFirstSession = payload;
        },
        changeFirstQuery(payload: string) {
            this.firstQuery = payload;
        }
    }
});


