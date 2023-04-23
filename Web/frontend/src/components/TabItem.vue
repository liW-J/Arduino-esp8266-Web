<template>
    <div class="tab-item" @click="emit('click')">{{ title }}</div>
</template>

<script setup lang='ts'>
type Props = {
    title: string
}

defineProps<Props>()
const emit = defineEmits(['click'])

</script>

<style lang='scss' scoped>
.tab-item {
    width: 36vh;
    min-width: 280px;
    height: 36vh;
    min-height: 280px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 4vh;
    box-shadow: var(--shadow-5);
    border-radius: 5px;
    cursor: pointer;
    margin: 5vh;
    transition: all .5s;

    &:hover {
        transform: skewX(3deg) rotateX(-3deg);
        box-shadow: var(--shadow-6);
    }
}
</style>