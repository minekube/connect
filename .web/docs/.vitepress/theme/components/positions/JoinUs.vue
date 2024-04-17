<template>
  <div class="py-32">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div id="positions" class="mx-auto flex max-w-2xl flex-col items-end justify-between gap-16 lg:mx-0 lg:max-w-none lg:flex-row">
        <div class="w-full lg:max-w-lg lg:flex-auto">
          <h2 class="text-3xl font-bold tracking-tight text-[--vp-c-text-1] sm:text-4xl">We’re always looking for awesome people to join us</h2>
          <p class="mt-6 text-xl leading-8 text-[--vp-c-text-2]">We are interested in passionate builders that are excited for innovating the Minecraft space through Minekube's Vision.</p>
          <img src="https://storage.googleapis.com/pai-images/e2f08836cd574c47a19adcbb7b8949de.jpeg" alt="" class="mt-16 aspect-[6/5] w-full rounded-2xl bg-gray-50 object-cover lg:aspect-auto lg:h-[34.5rem]" />
        </div>
        <div class="w-full lg:max-w-xl lg:flex-auto">
          <h3 class="sr-only">Job openings</h3>
          <ul class="-my-8 divide-y divide-gray-100">
            <li v-for="(opening, index) in jobOpenings.slice(0, 3)" :key="index" class="py-8">
              <PositionItem :opening="opening"/>
            </li>
          </ul>
          <div class="mt-8 flex border-t border-gray-100 pt-8">
            <a :href="discordLink" class="text-sm font-semibold leading-6 text-[--vp-c-brand-2] hover:text-[--vp-c-brand-1]">Talk to us <span aria-hidden="true">&rarr;</span></a>
          </div>
        </div>
      </div>
    </div>
    <ThreeColumns :openings="jobOpenings.slice(3)"/>
  </div>
</template>

<script setup lang="ts">
import {discordLink} from "../../../../shared";
import ThreeColumns from "./ThreeColumns.vue";
import PositionItem from "./PositionItem.vue";
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
</script>
