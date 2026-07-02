<script setup lang="ts">
import {data as posts} from './posts.data.js'
import { useData } from 'vitepress'
import JoinUs from "../positions/JoinUs.vue";

const { frontmatter } = useData()
</script>

<template>
  <div class="py-16 sm:py-18">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="text-3xl font-bold tracking-tight text-[--vp-c-text-1] sm:text-4xl">{{ frontmatter.title }}</h2>
        <p class="mt-2 text-lg leading-8 text-[--vp-c-text-2]">{{ frontmatter.subtext }}</p>
      </div>
      <div class="mx-auto mt-16 grid max-w-2xl grid-cols-1 gap-x-8 gap-y-20 lg:mx-0 lg:max-w-none lg:grid-cols-3">
        <article v-for="(post, index) in posts" :key="index" class="flex flex-col items-start justify-between">
          <a :href="post.url" class="relative w-full">
            <img :src="post.imageUrl" :alt="post.imageAlt || ''" class="aspect-[16/9] w-full rounded-2xl bg-gray-100 object-cover" />
            <div class="absolute inset-0 rounded-2xl ring-1 ring-inset ring-gray-900/10" />
          </a>
          <div class="max-w-xl">
            <div class="mt-8 flex items-center gap-x-4 text-xs">
              <time :datetime="new Date(post.date.time).toISOString()" class="text-[--vp-c-text-2]">{{ post.date.string }}</time>
              <span class="rounded-full bg-[--vp-button-brand-bg] px-3 py-1.5 font-medium text-[--vp-button-brand-text]">{{ post.category }}</span>
            </div>
            <div class="group relative">
              <h3 class="mt-3 text-lg font-semibold leading-6 text-[--vp-c-text-1] group-hover:text-[--vp-c-text-2]">
                <a :href="post.url">
                  {{ post.title }}
                </a>
              </h3>
              <p v-if="post.description" class="mt-5 line-clamp-3 text-sm leading-6 text-[--vp-c-text-2]">
                {{ post.description }}
              </p>
              <div v-else class="mt-5 line-clamp-3 text-sm leading-6 text-[--vp-c-text-2]" v-html="post.excerpt" />
            </div>
            <div class="relative mt-8 flex items-center gap-x-4">
              <img :src="post.author.imageUrl" alt="" class="h-10 w-10 rounded-full" />
              <div class="text-sm leading-6">
                <p class="font-semibold text-[--vp-c-text-1]">
                  <a :href="post.author.href" target="_blank" rel="noopener noreferrer">
                    {{ post.author.name }}
                  </a>
                </p>
                <p class="text-[--vp-c-text-2]">{{ post.author.role }}</p>
              </div>
            </div>
          </div>
        </article>
      </div>
    </div>
  </div>
  <JoinUs/>
</template>
