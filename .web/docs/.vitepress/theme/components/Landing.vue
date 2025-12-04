<script setup>

import Testimonials from "./Testimonials.vue";
import LogoCloud from "./LogoCloud.vue";
</script>

<template>


  <div class="">

    <div class="py-24 sm:py-32">
      <div class="mx-auto max-w-7xl px-6 lg:px-8">
        <div class="mx-auto max-w-2xl lg:text-center">
          <h2 class="text-base font-semibold leading-7 text-[--vp-c-brand-2]">Connect securely · Deploy faster</h2>
          <p class="mt-2 text-3xl font-bold tracking-tight text-[--vp-c-text-1] sm:text-4xl">We are changing the way
            developers build Minecraft networks</p>
          <p class="mt-6 text-lg leading-8 text-[--vp-c-text-2]">
            Minekube Connect is a cloud-native Minecraft TCP edge network, built on top of our open source Gate project.
            It's the Cloudflare for Minecraft.
          </p>
        </div>
        <div class="mx-auto mt-16 max-w-2xl sm:mt-20 lg:mt-24 lg:max-w-4xl">
          <dl class="grid max-w-xl grid-cols-1 gap-x-8 gap-y-10 lg:max-w-none lg:grid-cols-2 lg:gap-y-16">
            <div class="relative pl-16">
              <dt class="text-base font-semibold leading-7 text-[--vp-c-text-1]">
                <div
                    class="absolute left-0 top-0 flex h-10 w-10 items-center justify-center rounded-lg bg-[--vp-c-brand-2]">
                  <svg class="h-6 w-6 text-white lucide lucide-radio" fill="none" stroke="currentColor" stroke-linecap="round"
                       stroke-linejoin="round" stroke-width="2" viewBox="0 0 24 24"
                       xmlns="http://www.w3.org/2000/svg">
                    <path d="M4.9 19.1C1 15.2 1 8.8 4.9 4.9"/>
                    <path d="M7.8 16.2c-2.3-2.3-2.3-6.1 0-8.5"/>
                    <circle cx="12" cy="12" r="2"/>
                    <path d="M16.2 7.8c2.3 2.3 2.3 6.1 0 8.5"/>
                    <path d="M19.1 4.9C23 8.8 23 15.1 19.1 19"/>
                  </svg>
                </div>
                Minekube Connect
              </dt>
              <dd class="mt-2 text-base leading-7 text-[--vp-c-text-2]">Self-host your Connector, we manage the
                Superproxy infrastructure.
              </dd>
            </div>
            <div class="relative pl-16">
              <dt class="text-base font-semibold leading-7 text-[--vp-c-text-1]">
                <div
                    class="absolute left-0 top-0 flex h-10 w-10 items-center justify-center rounded-lg bg-[--vp-c-brand-2]">
                  <svg class="lucide lucide-gamepad-2 h-6 w-6 text-white" fill="none" stroke="currentColor"
                       stroke-linecap="round" stroke-linejoin="round" stroke-width="2" viewBox="0 0 24 24">
                    <line x1="6" x2="10" y1="11" y2="11"/>
                    <line x1="8" x2="8" y1="9" y2="13"/>
                    <line x1="15" x2="15.01" y1="12" y2="12"/>
                    <line x1="18" x2="18.01" y1="10" y2="10"/>
                    <path
                        d="M17.32 5H6.68a4 4 0 0 0-3.978 3.59c-.006.052-.01.101-.017.152C2.604 9.416 2 14.456 2 16a3 3 0 0 0 3 3c1 0 1.5-.5 2-1l1.414-1.414A2 2 0 0 1 9.828 16h4.344a2 2 0 0 1 1.414.586L17 18c.5.5 1 1 2 1a3 3 0 0 0 3-3c0-1.545-.604-6.584-.685-7.258-.007-.05-.011-.1-.017-.151A4 4 0 0 0 17.32 5z"/>
                  </svg>
                </div>
                Minekube Games
              </dt>
              <dd class="mt-2 text-base leading-7 text-[--vp-c-text-2]">You create the game experience, we bring the
                players and infrastructure. Matchmaking and scale on demand. (Coming Soon)
              </dd>
            </div>
            <div class="relative pl-16">
              <dt class="text-base font-semibold leading-7 text-[--vp-c-text-1]">
                <div
                    class="absolute left-0 top-0 flex h-10 w-10 items-center justify-center rounded-lg bg-[--vp-c-brand-2]">
                  <svg class="h-6 w-6 text-white lucide lucide-area-chart" fill="none" stroke="currentColor" stroke-linecap="round"
                       stroke-linejoin="round" stroke-width="2" viewBox="0 0 24 24"
                       xmlns="http://www.w3.org/2000/svg">
                    <path d="M3 3v18h18"/>
                    <path d="M7 12v5h12V8l-5 5-4-4Z"/>
                  </svg>
                </div>
                Minekube Analytics
              </dt>
              <dd class="mt-2 text-base leading-7 text-[--vp-c-text-2]">
                Get rich metrics and find insights about your players like, so you can anticipate what players want to
                play next. (Coming Soon)
              </dd>
            </div>
            <div class="relative pl-16">
              <dt class="text-base font-semibold leading-7 text-[--vp-c-text-1]">
                <div
                    class="absolute left-0 top-0 flex h-10 w-10 items-center justify-center rounded-lg bg-[--vp-c-brand-2]">
                  <svg class="h-6 w-6 text-white lucide lucide-blocks" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                       stroke-width="2" viewBox="0 0 24 24">
                    <rect height="7" rx="1" width="7" x="14" y="3"/>
                    <path
                        d="M10 21V8a1 1 0 0 0-1-1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-5a1 1 0 0 0-1-1H3"/>
                  </svg>
                </div>
                Minekube Developer Platform
              </dt>
              <dd class="mt-2 text-base leading-7 text-[--vp-c-text-2]">
                Coming Soon: Search Browser, Reviews & Votes, Marketplace, Ads
              </dd>
            </div>
          </dl>
        </div>
      </div>
    </div>


    <div class="overflow-hidden py-24 sm:py-32">
      <div class="mx-auto max-w-7xl px-6 lg:px-8">
        <div
            class="mx-auto grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 sm:gap-y-20 lg:mx-0 lg:max-w-none lg:grid-cols-2">
          <div class="lg:pr-8 lg:pt-4">
            <div class="lg:max-w-lg">
              <h2 class="text-base font-semibold leading-7 text-[--vp-c-brand-2]">Capable Connectivity</h2>
              <p class="mt-2 text-3xl font-bold tracking-tight text-[--vp-c-text-1] sm:text-4xl">The Cloudflare for
                Minecraft</p>
              <p class="mt-6 text-lg leading-8 text-[--vp-c-text-2]">
                Minekube Connect serves Minecraft TCP for free, where Cloudflare doesn't.
                Connect advances and deploys our well-known open source Gate proxy, globally.
              </p>
              <dl class="mt-10 max-w-xl space-y-8 text-base leading-7 text-[--vp-c-text-2] lg:max-w-none">
                <div class="relative pl-9">
                  <dt class="inline font-semibold text-[--vp-c-text-1]">
                    <svg aria-hidden="true" class="absolute left-1 top-1 h-5 w-5 text-[--vp-c-brand-2]"
                         fill="currentColor" viewBox="0 0 20 20">
                      <path clip-rule="evenodd"
                            d="M5.5 17a4.5 4.5 0 01-1.44-8.765 4.5 4.5 0 018.302-3.046 3.5 3.5 0 014.504 4.272A4 4 0 0115 17H5.5zm3.75-2.75a.75.75 0 001.5 0V9.66l1.95 2.1a.75.75 0 101.1-1.02l-3.25-3.5a.75.75 0 00-1.1 0l-3.25 3.5a.75.75 0 101.1 1.02l1.95-2.1v4.59z"
                            fill-rule="evenodd"/>
                    </svg>
                    Ingress Tunnels
                  </dt>
                  <dd class="inline">A DDoS-protected Minecraft native Superproxy at any edge near your players.
                  </dd>
                </div>
                <div class="relative pl-9">
                  <dt class="inline font-semibold text-[--vp-c-text-1]">
                    <svg class="absolute left-1 top-1 h-5 w-5 text-[--vp-c-brand-2] lucide lucide-globe" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                         stroke-width="2"
                         viewBox="0 0 24 24">
                      <circle cx="12" cy="12" r="10"/>
                      <path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20"/>
                      <path d="M2 12h20"/>
                    </svg>
                    Anycast domains
                  </dt>
                  <dd class="inline">Public IP address routing to your nearest Minecraft endpoints, protected by Connect
                    edge.
                  </dd>
                </div>
                <div class="relative pl-9">
                  <dt class="inline font-semibold text-[--vp-c-text-1]">
                    <svg class="absolute left-1 top-1 h-5 w-5 text-[--vp-c-brand-2] lucide lucide-package-check" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                         stroke-width="2"
                         viewBox="0 0 24 24">
                      <path d="m16 16 2 2 4-4"/>
                      <path
                          d="M21 10V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l2-1.14"/>
                      <path d="m7.5 4.27 9 5.15"/>
                      <polyline points="3.29 7 12 12 20.71 7"/>
                      <line x1="12" x2="12" y1="22" y2="12"/>
                    </svg>
                    <!--                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="absolute left-1 top-1 h-5 w-5 text-[&#45;&#45;vp-c-brand-2] lucide lucide-terminal"><polyline points="4 17 10 11 4 5"/><line x1="12" x2="20" y1="19" y2="19"/></svg>-->
                    Fast setup
                  </dt>
                  <dd class="inline">Build any network architecture you want.
                    One binary / container / plugin, and a few lines of configuration.

                  </dd>
                </div>
              </dl>
            </div>
          </div>
          <img alt="Product screenshot"
               class="w-[48rem] max-w-none rounded-xl shadow-xl ring-1 ring-gray-400/10 sm:w-[57rem] md:-ml-4 lg:-ml-0"
               height="1442"
               src="/images/dark-project-app-screenshot.png" width="2432">
        </div>
      </div>
    </div>

    <Testimonials/>

    <LogoCloud/>

  </div>


</template>

<style scoped>

</style>
