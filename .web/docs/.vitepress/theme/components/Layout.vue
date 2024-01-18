<script setup>
import DefaultTheme from 'vitepress/theme'
import {useRouter} from 'vitepress';
import {watch} from 'vue';
import MeetTeam from "./MeetTeam.vue";
import Landing from "./Landing.vue";

const {Layout} = DefaultTheme

const router = useRouter();

// Only run this on the client. Not during build
if (typeof window !== 'undefined' && window.posthog) {
  watch(() => router.route.data.relativePath, (path) => {
    posthog.capture("$pageview");
  }, {immediate: true});
}

</script>

<template>
  <Layout>
    <template #home-features-after>
      <Landing/>
      <MeetTeam/>
    </template>
  </Layout>
</template>