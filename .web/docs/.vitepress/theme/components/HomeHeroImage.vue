<script setup>
import Graph from "./Graph.vue";

const nodes = [
  {x: 100, y: 0, image: '/minekube-logo.png'}, // Top center node
  {x: 0, y: 150, image: 'https://avatars.githubusercontent.com/u/41710604?s=200&v=4'}, // Bottom left node
  {x: 100, y: 150, image: 'https://avatars.githubusercontent.com/u/7608950?s=200&v=4'}, // Bottom center node
  {x: 200, y: 150, image: 'https://avatars.githubusercontent.com/u/4350249?s=48&v=4'}, // Bottom right node
  {x: 200, y: 250, content: '<div class="text-[--vp-c-text-1] font-mono">Gate</div>'}
];

const connections = [
  [0, 1, 10], // Connect top center node to bottom left node
  [0, 2, 100], // Connect top center node to bottom center node
  [0, 3, 10] // Connect top center node to bottom right node
];

</script>

<template>

  <Graph :connections="connections" :nodes="nodes"/>

</template>

<style scoped>

</style>