<template>
    <div>
        <t-textarea resize="none" :autosize="false" v-model="value" placeholder="请输入描述文案" name="description" @change="onChange" />
    </div>
</template>
<script setup>
import { onMounted, watch, computed, ref, reactive } from 'vue';

const value = ref('');
const onChange = (value,e) => {
    console.log(value)
}
</script>
<style lang="less">
.chat {
    width: 800px;
    font-size: 20px;
    margin: 0px auto
}
</style>