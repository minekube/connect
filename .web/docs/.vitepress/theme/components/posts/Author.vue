<script setup lang="ts">
import {useData} from 'vitepress'
import {Post} from "./posts.data";

const {frontmatter}: { frontmatter: Post } = useData()
</script>

<template>
  <dl class="pt-6 pb-10 xl:pt-11 xl:border-b xl:border-gray-200 dark:xl:border-slate-200/5">
    <dt class="sr-only">Authors</dt>
    <dd>
      <ul
          class="flex justify-center xl:block space-x-8 sm:space-x-12 xl:space-x-0 xl:space-y-8"
      >
        <li class="flex items-center space-x-2">
          <img
              v-if="frontmatter.author.imageUrl"
              :src="frontmatter.author.imageUrl"
              alt="author image"
              class="w-10 h-10 rounded-full"
          />
          <dl class="text-sm font-medium leading-5 whitespace-nowrap">
            <dt class="sr-only">Name</dt>
            <dd class="text-[--vp-c-text-1]">
              <a
                  v-if="frontmatter.author.href"
                  :href="frontmatter.author.href"
                  target="_blank"
                  rel="noopnener noreferrer"
                  :class="{ 'underline': frontmatter.author.href }"
              >
                {{ frontmatter.author.name }}
              </a>
              <span v-else>{{ frontmatter.author.name }}</span>
            </dd>
          </dl>
        </li>
      </ul>
    </dd>
  </dl>
</template>
