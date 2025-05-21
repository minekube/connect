<script setup lang="ts">
import {discordLink} from "../../../../shared/index";
import {Position} from "./types";

const props = defineProps({
  opening: Object as () => Position,
})
</script>

<template>
  <dl class="relative flex flex-wrap gap-x-3">
    <dt class="sr-only">Role</dt>
    <dd class="w-full flex-none text-lg font-semibold tracking-tight text-[--vp-c-text-1]">
      <a :href="opening.href || discordLink" class="text-[--vp-c-brand-2] hover:text-[--vp-c-brand-1]">
        {{ opening.role }}
        <span class="absolute inset-0" aria-hidden="true"/>
      </a>
    </dd>
    <dt class="sr-only">Description</dt>
    <dd class="mt-2 w-full flex-none text-base leading-7 text-[--vp-c-text-2]">{{ opening.description }}</dd>
    <dt class="sr-only" v-if="opening.salary">Salary</dt>
    <dd class="mt-4 text-base font-semibold leading-7 text-[--vp-c-text-1]" v-if="opening.salary">{{
        opening.salary
      }}
    </dd>
    <dt class="sr-only">Location</dt>
    <dd class="mt-4 flex items-center gap-x-3 text-base leading-7 text-[--vp-c-text-3]">
      <svg viewBox="0 0 2 2" class="h-0.5 w-0.5 flex-none fill-gray-300" aria-hidden="true" v-if="opening.salary">
        <circle cx="1" cy="1" r="1"/>
      </svg>
      {{ opening.location }}
    </dd>
  </dl>
</template>
