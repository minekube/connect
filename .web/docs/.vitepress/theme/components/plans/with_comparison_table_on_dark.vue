<template>
  <div class="py-24 sm:py-32">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div class="mx-auto max-w-4xl text-center">
        <h2 class="text-base font-semibold leading-7 text-[--vp-c-brand-2]">Plans Comparison</h2>
        <p class="mt-2 text-4xl font-bold tracking-tight text-[--vp-c-text-1] sm:text-5xl">Plans for projects of&nbsp;all&nbsp;sizes</p>
      </div>
      <p class="mx-auto mt-6 max-w-2xl text-center text-lg leading-8 text-[--vp-c-text-1]">Not ready for Plus? Every organization starts on the generous Free forever plan. It’s great for side projects, or learning how to use Minekube</p>

      <!-- xs to lg -->
      <div class="mx-auto mt-12 max-w-md space-y-8 sm:mt-16 lg:hidden">
        <section v-for="tier in tiers" :key="tier.id" :class="[tier.mostPopular ? 'rounded-xl bg-white/5 ring-1 ring-inset ring-white/10' : '', 'p-8']">
          <h3 :id="tier.id" class="text-sm font-semibold leading-6 text-[--vp-c-text-1]">{{ tier.name }}</h3>
          <p class="mt-2 flex items-baseline gap-x-1">
            <span class="text-4xl font-bold text-[--vp-c-text-1]">{{ tier.priceMonthly }}</span>
            <span class="text-sm font-semibold text-[--vp-c-text-1]">/month</span>
          </p>
          <a :href="tier.href" :aria-describedby="tier.id" :class="[tier.mostPopular ? 'bg-indigo-500 text-[--vp-c-text-1] hover:bg-indigo-400 focus-visible:outline-indigo-500' : 'bg-white/10 text-[--vp-c-text-1] hover:bg-white/20 focus-visible:outline-white', 'mt-8 block rounded-md py-2 px-3 text-center text-sm font-semibold leading-6 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2']">Buy plan</a>
          <ul role="list" class="mt-10 space-y-4 text-sm leading-6 text-[--vp-c-text-1]">
            <li v-for="section in sections" :key="section.name">
              <ul role="list" class="space-y-4">
                <template v-for="feature in section.features" :key="feature.name">
                  <li v-if="feature.tiers[tier.name]" class="flex gap-x-3">
<!--                    <CheckIcon class="h-6 w-5 flex-none text-[--vp-c-brand-2]" aria-hidden="true" />-->
                    <svg class="h-6 w-5 flex-none text-[--vp-c-brand-2]" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                      <path fill-rule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clip-rule="evenodd" />
                    </svg>
                    <span>
                      {{ feature.name }}
                      {{ ' ' }}
                      <span v-if="typeof feature.tiers[tier.name] === 'string'" class="text-sm leading-6 text-gray-400">({{ feature.tiers[tier.name] }})</span>
                    </span>
                  </li>
                </template>
              </ul>
            </li>
          </ul>
        </section>
      </div>

      <!-- lg+ -->
      <div class="isolate mt-20 hidden lg:block">
        <div class="relative -mx-8">
          <div v-if="tiers.some((tier) => tier.mostPopular)" class="absolute inset-x-4 inset-y-0 -z-10 flex">
            <div class="flex w-1/4 px-4" aria-hidden="true" :style="{ marginLeft: `${(tiers.findIndex((tier) => tier.mostPopular) + 1) * 25}%` }">
              <div class="w-full rounded-t-xl border-x border-t border-[--vp-button-brand-border]/10 dark:border-white/10 dark:bg-white/5 bg-gray-500/5" />
            </div>
          </div>
          <table class="w-full table-fixed border-separate border-spacing-x-8 text-left">
            <caption class="sr-only">
              Pricing plan comparison
            </caption>
            <colgroup>
              <col class="w-1/4" />
              <col class="w-1/4" />
              <col class="w-1/4" />
              <col class="w-1/4" />
            </colgroup>
            <thead>
              <tr>
                <td />
                <th v-for="tier in tiers" :key="tier.id" scope="col" class="px-6 pt-6 xl:px-8 xl:pt-8">
                  <div class="text-sm font-semibold leading-7 text-[--vp-c-text-1]">{{ tier.name }}</div>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <th scope="row"><span class="sr-only">Price</span></th>
                <td v-for="tier in tiers" :key="tier.id" class="px-6 pt-2 xl:px-8">
                  <div class="flex items-baseline gap-x-1 text-[--vp-c-text-1]">
                    <span class="text-4xl font-bold">{{ tier.priceMonthly }}</span>
                    <span class="text-sm font-semibold leading-6">/month</span>
                  </div>
                  <a :href="tier.href" :class="[tier.mostPopular ? 'bg-black text-white dark:bg-white dark:text-black dark:hover:bg-red-600 hover:bg-red-600 dark:hover:text-[--vp-c-text-1] focus-visible:outline-indigo-600' : 'bg-[--vp-button-brand-bg] text-[--vp-button-brand-text] hover:bg-[--vp-button-brand-hover-bg] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[--vp-button-brand-active-bg]', 'mt-8 block rounded-md py-2 px-3 text-center text-sm font-semibold leading-6 text-[--vp-c-text-1] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2']">{{ tier.ctaText || 'Buy plan' }}</a>
                </td>
              </tr>
              <template v-for="(section, sectionIdx) in sections" :key="section.name">
                <tr>
                  <th scope="colgroup" colspan="4" :class="[sectionIdx === 0 ? 'pt-8' : 'pt-16', 'pb-4 text-sm font-semibold leading-6 text-[--vp-c-text-1]']">
                    {{ section.name }}
                    <div class="absolute inset-x-8 mt-4 h-px dark:bg-white/10 bg-gray-500/10" />
                  </th>
                </tr>
                <tr v-for="feature in section.features" :key="feature.name">
                  <th scope="row" class="py-4 text-sm font-normal leading-6 text-[--vp-c-text-1]">
                    {{ feature.name }}
                    <div class="absolute inset-x-8 mt-4 h-px dark:bg-white/5 bg-gray-500/5" />
                  </th>
                  <td v-for="tier in tiers" :key="tier.id" class="px-6 py-4 xl:px-8">
                    <div v-if="typeof feature.tiers[tier.name] === 'string'" class="text-center text-sm leading-6 text-[--vp-c-text-1]">{{ feature.tiers[tier.name] }}</div>
                    <template v-else>
                      <svg v-if="feature.tiers[tier.name] === true" class="mx-auto h-5 w-5 flex-none text-[--vp-c-brand-2]" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path fill-rule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clip-rule="evenodd" />
                      </svg>
<!--                      <CheckIcon v-if="feature.tiers[tier.name] === true" class="mx-auto h-5 w-5 text-[--vp-c-brand-2]" aria-hidden="true" />-->
                      <svg v-else class="mx-auto h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <path fill-rule="evenodd" d="M4 10a.75.75 0 01.75-.75h10.5a.75.75 0 010 1.5H4.75A.75.75 0 014 10z" clip-rule="evenodd" />
                      </svg>
<!--                      <MinusIcon v-else class="mx-auto h-5 w-5 text-gray-500" aria-hidden="true" />-->
                      <span class="sr-only">{{ feature.tiers[tier.name] === true ? 'Included' : 'Not included' }} in {{ tier.name }}</span>
                    </template>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// import { CheckIcon, MinusIcon } from '@vue-hero-icons/solid/icons'

const tiers = [
  {
    name: 'Free',
    id: 'tier-free-table',
    ctaText: 'Quick start',
    href: '/guide/quick-start',
    priceMonthly: '$0',
    description: 'Quis suspendisse ut fermentum neque vivamus non tellus.',
    mostPopular: false,
  },
  {
    name: 'Plus',
    id: 'tier-plus-table',
    ctaText: 'Coming soon',
    href: 'https://minekube.com/discord',
    priceMonthly: '$5',
    description: 'Quis eleifend a tincidunt pellentesque. A tempor in sed.',
    mostPopular: true,
  },
  // {
  //   name: 'Enterprise',
  //   id: 'tier-enterprise',
  //   href: '#',
  //   priceMonthly: '$59',
  //   description: 'Orci volutpat ut sed sed neque, dui eget. Quis tristique non.',
  //   mostPopular: false,
  // },
]
const sections = [
  {
    name: 'Platform Features',
    features: [
      { name: 'DDoS protected global edge', tiers: { Free: true, Plus: true } },
      { name: 'Custom domains', tiers: { Free: true, Plus: true } },
      { name: 'Unlimited endpoints & connectors', tiers: { Free: true, Plus: true } },
      { name: 'Unlimited players', tiers: { Free: true, Plus: true } },
      { name: 'Unlimited org members', tiers: { Free: true, Plus: true } },
      { name: 'Activity Tracking', tiers: { Free: true, Plus: true } },
      { name: 'Custom fallback motd', tiers: { Free: false, Plus: true } },
      // { name: 'Team members', tiers: { Free: 'Up to 20 users', Plus: 'Up to 50 users' } },
    ],
  },
  {
    name: 'Opt-out ecosystem default options',
    features: [
      { name: 'Prevent Fallback to Browser Hub', tiers: { Free: false, Plus: true } },
      { name: 'Prevent Browser Endpoint Listing', tiers: { Free: false, Plus: true } },
      { name: 'Prevent /browser Command', tiers: { Free: false, Plus: true } },
      { name: 'Prevent Fallback Motd', tiers: { Free: false, Plus: true } },
      { name: 'Prevent Default Tablist', tiers: { Free: false, Plus: true } },
    ],
  },
  {
    name: 'Support',
    features: [
      { name: '24/7 online community support', tiers: { Free: true, Plus: true } },
      { name: 'Developer response within 24 hours', tiers: { Free: false, Plus: true } },
      { name: '1:1 onboarding tour', tiers: { Free: false, Plus: true } },
      { name: 'Early access to new features', tiers: { Free: false, Plus: true } },
    ],
  }
  // {
  //   name: 'Reporting',
  //   features: [
  //     { name: 'Advanced analytics', tiers: { Basic: true, Essential: true, Premium: true } },
  //     { name: 'Basic reports', tiers: { Essential: true, Premium: true } },
  //     { name: 'Professional reports', tiers: { Premium: true } },
  //     { name: 'Custom report builder', tiers: { Premium: true } },
  //   ],
  // },
  // {
  //   name: 'Support',
  //   features: [
  //     { name: '24/7 online support', tiers: { Basic: true, Essential: true, Premium: true } },
  //     { name: 'Quarterly product workshops', tiers: { Essential: true, Premium: true } },
  //     { name: 'Priority phone support', tiers: { Essential: true, Premium: true } },
  //     { name: '1:1 onboarding tour', tiers: { Premium: true } },
  //   ],
  // },
]
</script>