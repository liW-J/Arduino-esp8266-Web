<template>
  <div class="finger-print">
    <div class="img">
      <div class="scanLine"></div>
    </div>
    <h3 @click="emit('click')">{{ title }}</h3>
  </div>
</template>

<script setup lang='ts'>
type Props = {
  title: string
}

defineProps<Props>()
const emit = defineEmits(['click'])
</script>

<style lang='scss' scoped>
.finger-print {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  .img {
    width: 30vh;
    height: 41.51vh;
    margin: 0 auto;
    background: url("@/assets/figure.png");
    background-size: cover;
    background-repeat: no-repeat;
    position: relative;

    &::before {
      content: "";
      width: 100%;
      height: 100%;
      background-color: rgb(210, 210, 210);
      position: absolute;
      z-index: -1;
    }

    &::after {
      content: "";
      width: 100%;
      height: 100%;
      background-color: aqua;
      position: absolute;
      z-index: -1;
      animation: scan 4s ease-in-out infinite;
    }

    .scanLine {
      position: absolute;
      height: 8px;
      border-radius: 10px;
      width: 120%;
      left: -10%;
      background-color: aqua;
      filter: drop-shadow(0 0 30px aqua) drop-shadow(0 0 30px aqua);
      animation: scanLine 4s ease-in-out infinite;
    }
  }

  @keyframes scan {

    0%,
    100% {
      height: 0%;
    }

    50% {
      height: 100%;
    }
  }

  @keyframes scanLine {

    0%,
    100% {
      top: -30px;
    }

    50% {
      top: calc(41.52vh + 30px);
    }
  }

  h3 {
    margin-top: 8vh;
    animation: blink 4s ease-in-out infinite;
    cursor: pointer;
    transition: all;

    &:hover {
      color: var(--blue-6);
    }
  }

  @keyframes blink {

    0%,
    100% {
      opacity: 0.6;
    }

    50% {
      opacity: 1;
    }
  }
}
</style>