<script setup lang="ts">
import { useData } from 'vitepress'
import Home from './Home.vue'
import Article from './Article.vue'
import NotFound from './NotFound.vue'
import {gitHubLink} from "../../../../shared";

const { page, frontmatter } = useData()
</script>

<template>
  <div class="antialiased">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 xl:max-w-5xl xl:px-0">
      <nav class="flex flex-col gap-5 py-10 font-bold sm:flex-row sm:items-center sm:justify-between">
        <a class="text-xl" href="/blog/" aria-label="The Minekube Blog">
          <img
            class="inline-block mr-2 w-8 h-8"
            alt="logo"
            src="/minekube-logo.png"
          />
          <span
            v-if="!frontmatter.index"
            class="hidden md:inline dark:text-white"
            >The Minekube Blog</span
          >
        </a>
        <div class="flex flex-wrap items-center gap-3 text-sm text-gray-500 dark:text-white leading-5">
          <a
            class="hover:text-gray-700 dark:hover:text-gray-200"
            :href="gitHubLink"
            target="_blank"
            rel="noopener"
            ><span class="hidden sm:inline">GitHub </span>Source</a
          >
          <a
            class="inline-flex items-center rounded-md border border-orange-500/40 px-3 py-2 font-semibold text-orange-600 hover:border-orange-500 hover:text-orange-700 dark:text-orange-300 dark:hover:text-orange-200"
            href="/feed.rss"
            aria-label="Subscribe to The Minekube Blog RSS feed"
            >RSS<span class="hidden sm:inline">&nbsp;Feed</span></a
          >
          <a
            class="hover:text-gray-700 dark:hover:text-gray-200"
            href="https://minekube.com"
            target="_blank"
            rel="noopener"
            >Minekube.com →</a
          >
        </div>
      </nav>
    </div>
    <main class="max-w-3xl mx-auto px-4 sm:px-6 xl:max-w-5xl xl:px-0">
      <Home v-if="frontmatter.index" />
      <NotFound v-else-if="page.isNotFound" />
      <Article v-else />
    </main>
  </div>
</template>
