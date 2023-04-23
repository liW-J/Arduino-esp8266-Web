<template>
    <router-view></router-view>
</template>

<script setup lang='ts'>
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/loading/style/css'
import 'element-plus/es/components/message-box/style/css'
</script>

<style>
@import "open-props/normalize";
@import "open-props/style";


:root {
    font-family: var(--font-mono);
}

#app {
    width: 100vw;
    height: 100vh;
}

::-webkit-scrollbar {
    display: none;
}

:where(p) {
    max-inline-size: none;
}
</style>