<template>
  <section id="positions" class="mk-jobs" aria-labelledby="jobs-title">
    <div class="head">
      <h2 id="jobs-title">Open positions</h2>
      <p>We’re looking for people who want to build the future of Minecraft infrastructure with us. Everything is remote. Pick an area to see what fits you.</p>
    </div>
    <div class="chips" role="tablist" aria-label="Area">
      <button v-for="g in groups" :key="g" type="button" role="tab" :aria-selected="area === g" @click="area = g">
        {{ g }} <span class="count">{{ g === 'All' ? jobOpenings.length : jobOpenings.filter((j) => groupOf(j.role) === g).length }}</span>
      </button>
    </div>
    <TransitionGroup tag="ul" name="job" class="jobs">
      <li v-for="job in shown" :key="job.role">
        <a :href="job.href || discordLink" target="_blank" rel="noreferrer">
          <div class="tags"><span class="tag level">{{ job.salary }}</span><span class="tag">Remote</span></div>
          <h3>{{ job.role }}</h3>
          <p>{{ job.description }}</p>
          <span class="go">Talk to us on Discord</span>
        </a>
      </li>
    </TransitionGroup>
  </section>
</template>

<script setup lang="ts">
import {discordLink} from "../../../shared";
import {computed, ref} from "vue";
import {Position} from "./types";

const jobOpenings = [
  {
    role: "Full-Stack Engineer",
    description: "Develop and maintain high-quality software. Collaborate with team to design and launch features.",
    salary: "Senior",
    location: "Remote (anywhere)",
  },
  {
    role: "Infrastructure Engineer",
    description: "Design, build and maintain infrastructure. Ensure scalability, security, reliability of systems.",
    salary: "Senior",
    location: "Remote (anywhere)",
  },
  {
    role: "Backend Engineer",
    description: "Design and implement scalable backend systems and APIs. Collaborate with team to design and launch features.",
    salary: "Level 2",
    location: "Remote (anywhere)",
  },
  {
    role: "Frontend Engineer",
    description: "Develop and maintain user-friendly web applications. Collaborate with team to design and launch features.",
    salary: "Level 2",
    location: "Remote (anywhere)",
  },
  {
    role: "Engineering Manager",
    description: "Lead and manage engineering team, set goals and priorities. Collaborate with cross-functional teams on technical and product strategies.",
    salary: "Senior",
    location: "Remote (anywhere)",
  },
  {
    role: "Marketing Manager",
    description: "Lead and execute marketing strategies across channels to increase brand awareness and drive customer acquisition. Collaborate with cross-functional teams on campaigns.",
    salary: "Level 1",
    location: "Remote (anywhere)",
  },
  {
    role: "Customer Success Manager",
    description: "Build and maintain customer relationships, ensuring value from platform. Collaborate with cross-functional teams to address customer needs and drive success.",
    salary: "Senior",
    location: "Remote (anywhere)",
  },
  {
    role: "Developer Relations Engineer",
    description: "Engage with and support developer community, providing technical guidance and resources. Collaborate on improving developer experience and adoption.",
    salary: "Level 2",
    location: "Remote (anywhere)",
  },
  {
    role: "Platform Product Engineer",
    description: "Design and implement platform features and integrations for customer applications. Collaborate with team to design and launch new features.",
    salary: "Level 2",
    location: "Remote (anywhere)",
  },
  {
    role: "Go Developer",
    description: "Develop and maintain high-quality, high-performance software in Go. Collaborate with team to design and launch features.",
    salary: "Level 2",
    location: "Remote (anywhere)",
  },
  {
    role: "Site Reliability Engineer",
    description: "Ensure reliability and performance of platform and services. Design and implement monitoring, alerting, and incident response systems.",
    salary: "Senior",
    location: "Remote (anywhere)",
  },
  {
    role: "Support Engineer",
    description: "Provide technical support and guidance to customers, troubleshoot and resolve issues. Collaborate with cross-functional teams to drive success.",
    salary: "Intern",
    location: "Remote (anywhere)",
  },
  {
    role: "Technical Writer",
    description: "Create and maintain technical documentation, guides, and tutorials. Collaborate with cross-functional teams to ensure accuracy and completeness.",
    salary: "Level 1",
    location: "Remote (anywhere)",
  },
  {
    role: "Community Moderator",
    description: "Manage and moderate online communities, ensuring positive and respectful environment. Collaborate on community engagement and growth.",
    salary: "Intern",
    location: "Remote (anywhere)",
  },
  {
    role: "Content Creator",
    description: "Create and publish engaging content across channels to increase brand awareness and drive engagement. Collaborate on content strategies.",
    salary: "Level 1",
    location: "Remote (anywhere)",
  },
];

// Group roles by area so visitors find what fits them quickly.
const groupOf = (role: string) =>
  /Marketing|Customer|Content/.test(role) ? 'Business'
  : /Community|Developer Relations|Writer|Support/.test(role) ? 'Community'
  : 'Engineering'
const groups = ['All', 'Engineering', 'Community', 'Business']
const area = ref('All')
const shown = computed(() => (area.value === 'All' ? jobOpenings : jobOpenings.filter((j) => groupOf(j.role) === area.value)))
</script>

<style scoped>
.mk-jobs { padding: clamp(4rem, 8vw, 6rem) 0 1rem; }
.head { display: grid; grid-template-columns: auto minmax(16rem, 30rem); justify-content: space-between; align-items: center; gap: 1.25rem clamp(2.5rem, 6vw, 6rem); }
.head h2 { white-space: nowrap; }
@media (max-width: 72rem) { .head h2 { white-space: normal; } }
@media (max-width: 56rem) { .head { grid-template-columns: 1fr; } }
.head h2 { margin: 0; font-family: var(--mk-display); font-weight: 750; font-stretch: 82%; font-size: clamp(2rem, 1.3rem + 2.2vw, 3.25rem); line-height: 1.02; letter-spacing: -0.02em; color: var(--sl-color-white); }
.head p { margin: 0; font-size: 1.0625rem; line-height: 1.6; color: var(--sl-color-gray-2); }
.chips { display: inline-flex; flex-wrap: wrap; gap: 0.3rem; padding: 0.3rem; margin-top: 2rem; border-radius: 999px; background: var(--sl-color-gray-6); }
.chips button { all: unset; cursor: pointer; display: inline-flex; align-items: center; gap: 0.45rem; padding: 0.5rem 1rem; border-radius: 999px; font-size: 0.9375rem; font-weight: 650; color: var(--sl-color-gray-3); transition: background-color 0.2s ease, color 0.2s ease; }
.chips button:hover { color: var(--sl-color-white); }
.chips button[aria-selected='true'] { background: var(--sl-color-black); color: var(--sl-color-white); }
.chips button:focus-visible { outline: 2px solid var(--mk-orange); outline-offset: 2px; }
.count { font-size: 0.75rem; padding: 0.05rem 0.45rem; border-radius: 999px; background: var(--sl-color-gray-5); color: var(--sl-color-gray-2); }
.jobs { list-style: none; padding: 0; margin: 2rem 0 0; display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 1.25rem; }
.jobs a { display: flex; flex-direction: column; gap: 0.6rem; height: 100%; box-sizing: border-box; padding: 1.5rem; border-radius: 1.25rem; background: var(--sl-color-gray-6); text-decoration: none; transition: background-color 0.18s ease, transform 0.18s ease; }
.jobs a:hover { background: var(--sl-color-gray-5); transform: translateY(-2px); }
.tags { display: flex; gap: 0.4rem; }
.tag { font-size: 0.75rem; font-weight: 650; padding: 0.2rem 0.6rem; border-radius: 999px; background: var(--sl-color-gray-5); color: var(--sl-color-gray-2); }
.tag.level { background: color-mix(in srgb, var(--mk-orange) 20%, transparent); color: var(--mk-orange); }
.jobs h3 { margin: 0.2rem 0 0; font-family: var(--sl-font); font-stretch: 100%; font-size: 1.1875rem; font-weight: 700; letter-spacing: -0.01em; color: var(--sl-color-white); }
.jobs p { margin: 0; flex: 1; font-size: 0.9375rem; line-height: 1.6; color: var(--sl-color-gray-2); }
.go { margin-top: 0.4rem; font-size: 0.9375rem; font-weight: 650; color: var(--mk-orange); transition: color 0.2s ease; }
.go::after { content: ' →'; }
.jobs a:hover .go { color: var(--sl-color-white); }
.job-enter-active, .job-leave-active { transition: opacity 0.25s ease, transform 0.25s ease; }
.job-enter-from, .job-leave-to { opacity: 0; transform: translateY(8px); }
.job-leave-active { position: absolute; visibility: hidden; }
.job-move { transition: transform 0.3s ease; }
@media (max-width: 56rem) { .jobs { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
@media (max-width: 36rem) { .jobs { grid-template-columns: 1fr; } }
@media (prefers-reduced-motion: reduce) { .job-enter-active, .job-leave-active, .job-move, .jobs a { transition: none; } }
</style>
