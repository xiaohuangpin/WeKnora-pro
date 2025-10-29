<script setup lang="ts">
</script>
<template>
    <div class="empty">
        <img class="empty-img" src="@/assets/img/upload.svg" alt="">
        <span class="empty-txt">知识为空，拖放上传</span>
        <span class="empty-type-txt">pdf、doc 格式文件，不超过10M</span>
        <span class="empty-type-txt">text、markdown格式文件，不超过200K</span>
    </div>
</template>
<style scoped lang="less">
.empty {
    flex: 1;
    display: flex;
    flex-flow: column;
    justify-content: center;
    align-items: center;
}

.empty-txt {
    color: #00000099;
    font-family: "PingFang SC";
    font-size: 16px;
    font-weight: 600;
    line-height: 26px;
    margin: 12px 0 16px 0;
}

.empty-type-txt {
    color: #00000066;
    text-align: center;
    font-family: "PingFang SC";
    font-size: 12px;
    font-weight: 400;
    width: 217px;
}

.empty-img {
    width: 162px;
    height: 162px;
}
</style>
