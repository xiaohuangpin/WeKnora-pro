<template>
    <div class="user_msg">
        {{ content }}
    </div>
</template>
<script setup>
import { ref, defineProps } from "vue";

const props = defineProps({
    // 必填项
    content: {
        type: String,
        required: false
    }
});
</script>
<style scoped lang="less">
.user_msg {
    width: max-content;
    max-width: 776px;
    display: flex;
    padding: 10px 12px;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 4px;
    flex: 1 0 0;
    border-radius: 4px;
    background: #8CE97F;
    margin-left: auto;
    color: #000000e6;
    font-size: 16px;
    text-align: justify;
    word-break: break-all;
    max-width: 100%;
    box-sizing: border-box;
}
</style>