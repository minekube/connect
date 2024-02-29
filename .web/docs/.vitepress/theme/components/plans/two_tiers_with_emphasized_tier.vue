<template>
  <div class="relative isolate px-6 py-24 sm:py-32 lg:px-8">
    <div class="absolute inset-x-0 -top-3 -z-10 transform-gpu overflow-hidden px-36 blur-3xl" aria-hidden="true">
      <div class="mx-auto aspect-[1155/678] w-[72.1875rem] bg-gradient-to-tr from-[--vp-c-brand-1] to-[#ff80b5] opacity-30" style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)" />
    </div>
    <div class="mx-auto max-w-2xl text-center lg:max-w-4xl">
      <h2 class="text-base font-semibold leading-7 text-[--vp-c-brand-2]">Minekube Connect</h2>
      <p class="mt-2 text-4xl font-bold tracking-tight text-[--vp-c-text-1] sm:text-5xl">Pricing Plans</p>
    </div>
    <p class="mx-auto mt-6 max-w-2xl text-center text-lg leading-8 text-[--vp-c-text-1]">Use Minekube for free with your whole team. Upgrade to opt out of Browser ecosystem, and enable additional features.</p>
    <div class="mx-auto mt-16 grid max-w-lg grid-cols-1 items-center gap-y-6 sm:mt-20 sm:gap-y-0 lg:max-w-4xl lg:grid-cols-2">
      <div v-for="(tier, tierIdx) in tiers" :key="tier.id" :class="[tier.featured ? 'dark relative bg-red-950 shadow-2xl' : 'dark:bg-transparent/60 bg-[--vp-c-default-soft] sm:mx-8 lg:mx-0', tier.featured ? '' : tierIdx === 0 ? 'rounded-t-3xl sm:rounded-b-none lg:rounded-tr-none lg:rounded-bl-3xl' : 'sm:rounded-t-none lg:rounded-tr-3xl lg:rounded-bl-none', 'rounded-3xl p-8 ring-1 ring-gray-900/10 sm:p-10']">
        <h3 :id="tier.id" :class="[tier.featured ? 'text-[--vp-c-brand-2]' : 'text-[--vp-c-brand-2]', 'text-base font-semibold leading-7']">{{ tier.name }}</h3>
        <p class="mt-4 flex items-baseline gap-x-2">
          <span v-if="!(discount.active && tier.featured)" :class="[tier.featured ? 'text-[--vp-c-text-1]' : 'text-[--vp-c-text-1]', 'text-5xl font-bold tracking-tight']">{{ tier.priceMonthly }}</span>
          <span v-if="discount.active && tier.featured" class="text-red-500 text-5xl font-bold tracking-tight">{{ discount.price }}</span>
          <span v-if="discount.active && tier.featured" class="text-[--vp-c-text-1] line-through text-xl font-bold tracking-tight">{{ tier.priceMonthly }}</span>
          <span :class="[tier.featured ? 'text-[--vp-c-text-2]' : 'text-[--vp-c-text-3]', 'text-base']">/month</span>
        </p>
        <div v-if="discount.active && tier.featured" class="mt-4 p-1">
          <countdown :end-time="discount.endTime" class="text-red-400 text-lg font-bold" />
          <p class="text-[--vp-c-text-1]">{{ discount.note }}</p>
        </div>
         <p :class="[tier.featured ? 'text-[--vp-c-text-1]' : 'text-[--vp-c-text-2]', 'mt-6 text-base leading-7']">{{ tier.description }}</p>
        <ul role="list" :class="[tier.featured ? 'text-[--vp-c-text-1]' : 'text-[--vp-c-text-2]', 'mt-8 space-y-3 text-sm leading-6 sm:mt-10']">
          <li v-for="feature in tier.features" :key="feature" class="flex gap-x-3">
            <svg :class="[tier.featured ? 'text-[--vp-c-brand-2]' : 'text-[--vp-c-brand-2]', 'h-6 w-5 flex-none']" class="h-6 w-5 flex-none text-[--vp-c-brand-2]" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
              <path fill-rule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clip-rule="evenodd" />
            </svg>
            {{ feature }}
          </li>
        </ul>
        <a :href="tier.href" :aria-describedby="tier.id" :class="[tier.featured ? 'bg-white text-black shadow-sm hover:bg-red-600 hover:text-white focus-visible:outline-indigo-500' : 'bg-[--vp-button-brand-bg] text-[--vp-button-brand-text] hover:bg-[--vp-button-brand-hover-bg] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[--vp-button-brand-active-bg]', 'mt-8 block rounded-md py-2.5 px-3.5 text-center text-sm font-semibold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 sm:mt-10']">{{ tier.ctaText || 'Get started today' }}</a>
      </div>
    </div>
  </div>
</template>

<script setup>
import Countdown from "./countdown.vue";
import {discount, plans} from "./settings";

const tiers = [
  {
    name: 'Free for everyone forever',
    id: 'tier-free',
    ctaText: 'Speedrun: Connect quick start',
    href: '/guide/quick-start',
    priceMonthly: '$0',
    description: "Great for side projects, or projects with a small team. Build DDoS protected Minecraft networks close to your users.",
    features: [
      'Unlimited endpoints & players',
      'DDoS protected global edge',
      'Unlimited org members',
      'Custom domains',
    ],
    featured: false,
  },
  {
    name: 'Plus',
    id: 'tier-plus',
    href: plans.plus.href,
    ctaText: plans.plus.ctaText,
    priceMonthly: plans.plus.price,
    description: 'Enhance branding and unlock features to scale your project. Everything in Free, plus:',
    features: [
      'Opt-out of Browser Hub fallback',
      'Opt-out of the default tablist',
      'Opt-out of fallback motd',
      'Custom fallback motd',
      'Priority support',
      'Automated chat moderation',
    ],
    featured: true,
  },
]
</script>