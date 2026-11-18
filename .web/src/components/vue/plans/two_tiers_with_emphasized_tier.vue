<template>
  <div class="mk-tiers">
    <p class="mk-tiers-lead">Use Minekube for free with your whole team. Upgrade to opt out of the Browser ecosystem and enable additional features.</p>
    <div class="mk-tier-grid">
      <section v-for="tier in tiers" :key="tier.id" :class="['mk-tier', tier.featured ? 'is-featured' : '']" :aria-labelledby="tier.id">
        <h2 :id="tier.id" class="mk-tier-name">{{ tier.name }}</h2>
        <p class="mk-tier-price">
          <span v-if="!(discount.active && tier.featured)">{{ tier.priceMonthly }}</span>
          <span v-if="discount.active && tier.featured">{{ discount.price }}</span>
          <s v-if="discount.active && tier.featured" class="mk-tier-was">{{ tier.priceMonthly }}</s>
          <span class="mk-tier-per">per month</span>
        </p>
        <div v-if="discount.active && tier.featured" class="mk-tier-discount">
          <countdown :end-time="discount.endTime" />
          <p>{{ discount.note }}</p>
        </div>
        <p class="mk-tier-desc">{{ tier.description }}</p>
        <ul role="list" class="mk-tier-features">
          <li v-for="feature in tier.features" :key="feature">{{ feature }}</li>
        </ul>
        <a :href="tier.href" :aria-describedby="tier.id" class="mk-tier-cta">{{ tier.ctaText || 'Get started today' }}</a>
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
.mk-tiers { padding: 0 0 clamp(4rem, 8vw, 7rem); }
.mk-tiers-lead { max-width: 38rem; font-size: 1.25rem; line-height: 1.4; color: var(--sl-color-gray-1); margin: 0 0 3rem; }
.mk-tier-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 1rem; }
.mk-tier { border-radius: 1rem; background: var(--sl-color-gray-6); color: var(--sl-color-white); padding: clamp(1.75rem, 4vw, 3.5rem); display: flex; flex-direction: column; }
.mk-tier.is-featured { background: var(--mk-orange); color: #140a02; }
.mk-tier-name { margin: 0; font-size: 1.125rem; font-weight: 700; font-family: var(--sl-font); font-stretch: 100%; letter-spacing: 0; color: inherit; }
.mk-tier-price { margin: 1rem 0 0; display: flex; align-items: baseline; gap: 0.75rem; flex-wrap: wrap; }
.mk-tier-price > span:first-child { font-family: var(--mk-display); font-weight: 800; font-stretch: 80%; font-size: clamp(3.5rem, 2.5rem + 3vw, 5rem); line-height: 0.95; letter-spacing: -0.01em; }
.mk-tier-was { font-size: 1.5rem; font-weight: 700; }
.mk-tier-per { font-size: 1.0625rem; font-weight: 600; }
.mk-tier-discount { margin-top: 1rem; font-weight: 600; }
.mk-tier-desc { margin: 1.5rem 0 0; font-size: 1.125rem; line-height: 1.5; max-width: 30rem; }
.mk-tier-features { list-style: none; padding: 0; margin: 2rem 0 0; display: grid; gap: 0.6rem; font-size: 1.0625rem; }
.mk-tier-features li { padding-left: 1.5rem; position: relative; }
.mk-tier-features li::before { content: ''; position: absolute; left: 0; top: 0.5em; width: 0.55rem; height: 0.55rem; border-radius: 0.15rem; background: currentColor; }
.mk-tier-cta { margin-top: auto; padding-top: 0; align-self: start; display: inline-flex; min-height: 3.25rem; align-items: center; padding-inline: 1.5rem; margin-block-start: 2.5rem; font-weight: 650; text-decoration: none; border-radius: 0.5rem; background: var(--mk-orange); color: #140a02; transition: transform 0.15s ease; }
.is-featured .mk-tier-cta { background: #140a02; color: #fff; }
.mk-tier-cta:hover { transform: translateY(-2px); }
@media (max-width: 50rem) { .mk-tier-grid { grid-template-columns: 1fr; } }
</style>
