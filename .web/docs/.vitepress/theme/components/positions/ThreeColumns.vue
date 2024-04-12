<template>
  <div>
    <div class="mx-auto max-w-7xl px-6 py-16 sm:py-24 lg:px-8">
      <h2 class="text-2xl font-bold leading-10 tracking-tight text-[--vp-c-text-1]">Do you see a position you like?</h2>
      <p class="mt-6 max-w-2xl text-base leading-7 text-[--vp-c-text-1]">
        If you are interested in any of the positions below,
        <a :href="discordLink" class="font-semibold text-[--vp-c-brand-2] hover:text-[--vp-c-brand-2]">send us a message on Discord</a>
        with your short bio and what you would like to build for the Minkube community.
      </p>
      <div class="mt-20">
        <dl class="space-y-16 sm:grid sm:grid-cols-2 sm:gap-x-6 sm:gap-y-16 sm:space-y-0 lg:grid-cols-3 lg:gap-x-10">
          <div  v-for="(opening, index) in openings" :key="index">
            <PositionItem :opening="opening"/>
          </div>
        </dl>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'
import {Position} from "./types.js";
import PositionItem from "./PositionItem.vue";
import {discordLink} from "../../../../shared";

const props = defineProps({
  openings: Array as () => Position[] // not working
})
</script>