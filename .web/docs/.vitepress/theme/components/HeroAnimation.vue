<template>
  <div id="animation-container" class="relative h-128 w-128">
    <div class="absolute left-1/2 animate-logo">
      <img alt="Bottom Middle Node" src="/minekube-logo.png" width="128">
    </div>
    <div class="absolute left-1/2 top-64 animate-logo">
      <img alt="Top Middle Node" src="/minekube-logo.png" width="128">
    </div>
    <div class="line animate-logo "></div>
  </div>
</template>

<style scoped>
@keyframes logo-animation {
  0% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(20px);
  }
  100% {
    transform: translateY(0);
  }
}

.animate-logo {
  animation: logo-animation 4s infinite ease-in-out;
}

@keyframes line-animation {
  0% {
    background-position: 0% 0%;
  }
  100% {
    background-position: 100% 100%;
  }
}

.line {
  position: absolute;
  top: 64px;
  left: 50%;
  width: 2px;
  height: 64px;
  background: linear-gradient(90deg, transparent, #d00808, transparent);
  background-size: 200% 200%;
  animation: line-animation 2s infinite;
}
</style>
<script setup>
</script>