<script setup>
import { ref } from 'vue';

const youtubeVideoId = 'jcN9Dvt1zLk';
const videoLoaded = ref(false);
const iframeSrc = ref('');

// YouTube thumbnail URL - using maxresdefault for best quality
const thumbnailUrl = `https://img.youtube.com/vi/${youtubeVideoId}/maxresdefault.jpg`;

// Embed URL with parameters for better UX:
// - rel=0: Don't show related videos from other channels
// - modestbranding=1: Reduce YouTube branding
// - playsinline=1: Better mobile experience
// - enablejsapi=1: Enable YouTube IFrame API for programmatic control
// - autoplay=1: Start playing immediately when iframe loads
// - mute=1: Muted by default (required for autoplay in most browsers)
const getEmbedUrl = () => {
  return `https://www.youtube.com/embed/${youtubeVideoId}?rel=0&modestbranding=1&playsinline=1&enablejsapi=1&autoplay=1&mute=1`;
};

const loadVideo = () => {
  videoLoaded.value = true;
  // Set src dynamically on click to ensure autoplay works
  iframeSrc.value = getEmbedUrl();
};
</script>

<template>
  <div class="w-full py-24 sm:py-32">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <!-- Optional Section Header (matches Landing.vue style) -->
      <div class="mx-auto max-w-2xl text-center mb-12">
        <h2 class="text-base font-semibold leading-7 text-[--vp-c-brand-2]">Watch & Learn</h2>
        <p class="mt-2 text-3xl font-bold tracking-tight text-[--vp-c-text-1] sm:text-4xl">
          See Minekube Connect in Action
        </p>
        <p class="mt-6 text-lg leading-8 text-[--vp-c-text-2]">
          Watch how easy it is to connect your Minecraft server to the global network
        </p>
      </div>

      <!-- Video Container with Lazy Loading -->
      <div class="mx-auto max-w-5xl">
        <div class="relative w-full rounded-xl overflow-hidden shadow-2xl ring-1 ring-gray-400/10 dark:ring-gray-700/20" style="padding-bottom: 56.25%;">
          <!-- Click-to-play thumbnail (shown before video loads) -->
          <div
            v-if="!videoLoaded"
            @click="loadVideo"
            @keydown.enter="loadVideo"
            @keydown.space.prevent="loadVideo"
            class="absolute inset-0 cursor-pointer group focus:outline-none focus:ring-2 focus:ring-[--vp-c-brand-2] focus:ring-offset-2 rounded-xl"
            role="button"
            tabindex="0"
            aria-label="Play video"
          >
            <img
              :src="thumbnailUrl"
              alt="Video thumbnail"
              class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
              loading="lazy"
            />
            <!-- Play button overlay -->
            <div class="absolute inset-0 flex items-center justify-center bg-black/30 group-hover:bg-black/40 transition-colors duration-300">
              <div class="w-20 h-20 sm:w-24 sm:h-24 rounded-full bg-white/95 group-hover:bg-white flex items-center justify-center shadow-xl transform group-hover:scale-110 transition-all duration-300 group-focus:ring-4 group-focus:ring-[--vp-c-brand-2]/50">
                <svg
                  class="w-8 h-8 sm:w-10 sm:h-10 text-[--vp-c-brand-2] ml-1"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path d="M8 5v14l11-7z"/>
                </svg>
              </div>
            </div>
          </div>

          <!-- YouTube iframe (loaded on click) -->
          <iframe
            v-if="videoLoaded"
            class="absolute top-0 left-0 w-full h-full"
            :src="iframeSrc"
            title="YouTube video player"
            frameborder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
            allowfullscreen
          ></iframe>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Smooth transitions for better UX */
.group {
  transition: all 0.3s ease-in-out;
}

/* Ensure proper focus styles for accessibility */
.group:focus-visible {
  outline: 2px solid var(--vp-c-brand-2);
  outline-offset: 2px;
}
</style>

