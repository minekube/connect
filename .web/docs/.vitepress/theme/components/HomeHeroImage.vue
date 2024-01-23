<script setup>
import Graph from "./Graph.vue";
import VPImage from "./VPImage.vue";

const velocity = 'https://avatars.githubusercontent.com/u/41710604?s=200&v=4'
const paper = 'https://avatars.githubusercontent.com/u/7608950?s=200&v=4'
const spigot = 'https://avatars.githubusercontent.com/u/4350249?s=48&v=4'

const center = {x: 100, y: 0};
const space = 100; // Define your space constant here

const nodes = [
  {x: center.x, y: center.y, link: '/guide/', image: '/minekube-logo.png', class: 'w-48'}, // Top center node
  {x: center.x - space * 1.5, y: 200, link: '/guide/connectors/gate', image: '/gate.png', class: 'w-32'}, // First bottom node (Gate)
  {x: center.x - space * 0.5, y: 250, link: '/guide/connectors/plugin', image: velocity}, // Second bottom node
  {x: center.x + space * 0.5, y: 250, link: '/guide/connectors/plugin', image: paper, class: 'w-16'}, // Third bottom node
  {x: center.x + space * 1.5, y: 200, link: '/guide/connectors/plugin', image: spigot}, // Fourth bottom node

  // {x: center.x, y: 0, image: '/mascot.png', class: 'w-12'},

];

// Define connections between nodes by their index
const connections = [
  [0, 1],
  [0, 2],
  [0, 3],
  [0, 4],
];

</script>

<template>

  <div class="hidden lg:block">
    <Graph :connections="connections" :nodes="nodes"/>
  </div>

  <div class="block lg:hidden">
    <VPImage class="image-src" image="/minekube-logo.png"/>
  </div>

</template>

<style scoped>

</style>