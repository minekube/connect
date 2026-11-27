<template>
  <div class="mk-tiers">
    <p class="lead">Use Minekube for free with your whole team. Upgrade to opt out of the Browser ecosystem and unlock more control over your network.</p>
    <div class="grid">
      <section v-for="tier in tiers" :key="tier.id" :class="['tier', tier.featured ? 'plus' : 'free']" :aria-labelledby="tier.id">
        <div class="top">
          <h2 :id="tier.id">{{ tier.featured ? 'Plus' : 'Free' }}</h2>
          <span v-if="tier.featured" class="pill">Recommended</span>
        </div>
        <p class="price">
          <span class="amount" v-if="!(discount.active && tier.featured)">{{ tier.priceMonthly }}</span>
          <span class="amount" v-if="discount.active && tier.featured">{{ discount.price }}</span>
          <s v-if="discount.active && tier.featured" class="was">{{ tier.priceMonthly }}</s>
          <span class="per">/ month</span>
        </p>
        <div v-if="discount.active && tier.featured" class="discount">
          <countdown :end-time="discount.endTime" />
          <p>{{ discount.note }}</p>
        </div>
        <p class="desc">{{ tier.description }}</p>
        <a :href="tier.href" :aria-describedby="tier.id" class="cta">{{ tier.ctaText || 'Get started today' }}</a>
        <p class="includes">{{ tier.featured ? 'Everything in Free, plus:' : 'Included for everyone:' }}</p>
        <ul role="list">
          <li v-for="feature in tier.features" :key="feature"><svg class="ck" viewBox="0 0 20 20" aria-hidden="true"><path fill="currentColor" d="M16.7 4.2a.75.75 0 0 1 .1 1l-8 10.5a.75.75 0 0 1-1.1.1l-4.5-4.5a.75.75 0 1 1 1-1l3.9 3.8 7.5-9.8a.75.75 0 0 1 1.1-.1Z"/></svg>{{ feature }}</li>
        </ul>
      </section>
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
    priceMonthly: plans.free.price,
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

<style scoped>
.mk-tiers { padding-bottom: clamp(4rem, 8vw, 6rem); }
.lead { max-width: 40rem; font-size: 1.25rem; line-height: 1.45; color: var(--sl-color-gray-1); margin: 0 0 1.75rem; }
.grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 1.25rem; align-items: stretch; }
.tier { position: relative; display: flex; flex-direction: column; padding: clamp(1.5rem, 2.5vw, 2rem); border-radius: 1.25rem; background: var(--sl-color-gray-6); color: var(--sl-color-white); }
.tier.plus { background: color-mix(in srgb, var(--mk-orange) 9%, var(--sl-color-gray-6)); box-shadow: inset 0 0 0 2px var(--mk-orange); }
.top { display: flex; align-items: center; justify-content: space-between; gap: 1rem; }
h2 { margin: 0; font-family: var(--sl-font); font-stretch: 100%; font-size: 1.25rem; font-weight: 700; letter-spacing: 0; color: inherit; }
.plus h2 { color: var(--mk-orange-text); }
.pill { font-size: 0.8125rem; font-weight: 700; padding: 0.3rem 0.7rem; border-radius: 999px; background: var(--mk-orange); color: #160a01; }
.price { margin: 0.9rem 0 0; display: flex; align-items: baseline; gap: 0.5rem; }
.amount { font-family: var(--mk-display); font-weight: 800; font-stretch: 80%; font-size: clamp(2.75rem, 2.2rem + 1.8vw, 3.75rem); line-height: 1; letter-spacing: -0.03em; }
.was { font-size: 1.25rem; color: var(--sl-color-gray-3); }
.per { font-size: 1.0625rem; color: var(--sl-color-gray-3); font-weight: 600; }
.discount { margin-top: 0.75rem; font-weight: 600; color: var(--mk-orange-text); }
.desc { margin: 0.75rem 0 0; font-size: 1rem; line-height: 1.5; color: var(--sl-color-gray-2); min-height: 3em; }
.cta { margin-top: 1.25rem; display: flex; align-items: center; justify-content: center; min-height: 3rem; border-radius: 999px; font-weight: 700; text-decoration: none; background: var(--sl-color-gray-5); color: var(--sl-color-white); transition: background-color 0.15s ease, transform 0.15s ease; }
.cta:hover { transform: translateY(-1px); background: var(--sl-color-gray-4); }
.plus .cta { background: var(--mk-orange); color: #160a01; }
.plus .cta:hover { background: #ffad33; }
.includes { margin: 1.4rem 0 0.6rem; font-size: 0.9375rem; font-weight: 700; color: var(--sl-color-white); }
ul { list-style: none; padding: 0; margin: 0; display: grid; grid-template-columns: 1fr 1fr; gap: 0.5rem 1rem; font-size: 0.9375rem; color: var(--sl-color-gray-1); }
li { display: flex; gap: 0.6rem; align-items: flex-start; }
.ck { flex: none; width: 1.1rem; height: 1.1rem; margin-top: 0.2em; color: var(--sl-color-gray-3); }
.plus .ck { color: var(--mk-orange-text); }
@media (max-width: 50rem) { .grid { grid-template-columns: 1fr; } }
@media (max-width: 30rem) { ul { grid-template-columns: 1fr; } }
</style>
