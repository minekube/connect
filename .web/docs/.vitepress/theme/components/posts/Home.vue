<script setup lang="ts">
import {data as posts} from './posts.data.js'
import { useData } from 'vitepress'
import JoinUs from "../positions/JoinUs.vue";

const { frontmatter } = useData()

const postsX = [
  {
    id: 1,
    title: 'Boost your conversion rate',
    href: '#',
    description:
      'Illo sint voluptas. Error voluptates culpa eligendi. Hic vel totam vitae illo. Non aliquid explicabo necessitatibus unde. Sed exercitationem placeat consectetur nulla deserunt vel. Iusto corrupti dicta.',
    imageUrl: 'https://images.playground.com/6f527647f6ed480b92df1213fceeaaaf.jpeg',
    date: 'Mar 16, 2020',
    datetime: '2020-03-16',
    category: { title: 'Engineering', href: '#' },
    author: {
      name: 'Robin Brämer',
      role: 'Founder / CTO',
      href: '#',
      imageUrl: 'https://github.com/robinbraemer.png'
    },
  },
  // More posts...
]

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
          <a :href="post.url">
            <div class="relative w-full">
              <img :src="post.imageUrl" alt="" class="aspect-[16/9] w-full rounded-2xl bg-gray-100 object-cover sm:aspect-[2/1] lg:aspect-[3/2]" />
              <div class="absolute inset-0 rounded-2xl ring-1 ring-inset ring-gray-900/10" />
            </div>
            <div class="max-w-xl">
              <div class="mt-8 flex items-center gap-x-4 text-xs">
                <time :datetime="new Date(post.date.time).toISOString()" class="text-[--vp-c-text-2]">{{ post.date.string }}</time>
                <a href="#" class="relative z-10 rounded-full bg-[--vp-button-brand-bg] px-3 py-1.5 font-medium text-[--vp-button-brand-text] hover:bg-[--vp-button-brand-hover-bg]">{{ post.category }}</a>
              </div>
              <div class="group relative">
                <h3 class="mt-3 text-lg font-semibold leading-6 text-[--vp-c-text-1] group-hover:text-[--vp-c-text-2]">
                  <span class="absolute inset-0" />
                  {{ post.title }}
                </h3>
                <p class="mt-5 line-clamp-3 text-sm leading-6 text-[--vp-c-text-2]" >
                  <div v-html="post.excerpt" />
                </p>
              </div>
              <div class="relative mt-8 flex items-center gap-x-4">
                <img :src="post.author.imageUrl" alt="" class="h-10 w-10 rounded-full" />
                <div class="text-sm leading-6">
                  <p class="font-semibold text-[--vp-c-text-1]">
                    <a :href="post.author.href">
                      <span class="absolute inset-0" />
                      {{ post.author.name }}
                    </a>
                  </p>
                  <p class="text-[--vp-c-text-2]">{{ post.author.role }}</p>
                </div>
              </div>
            </div>
          </a>
        </article>
      </div>
    </div>
  </div>
  <JoinUs/>
</template>
