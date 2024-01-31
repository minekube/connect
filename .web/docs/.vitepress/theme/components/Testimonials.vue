<template>
  <div class="relative isolate pb-32 pt-24 sm:pt-32">
    <div aria-hidden="true"
         class="absolute inset-x-0 top-1/2 -z-10 -translate-y-1/2 transform-gpu overflow-hidden opacity-30 blur-3xl">
      <div
          class="ml-[max(50%,38rem)] aspect-[1313/771] w-[82.0625rem] bg-gradient-to-tr from-[--vp-c-brand-1] to-[#ff80b5]"
           style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"/>
    </div>
    <div
        aria-hidden="true"
        class="absolute inset-x-0 top-0 -z-10 flex transform-gpu overflow-hidden pt-32 opacity-25 blur-3xl sm:pt-40 xl:justify-end">
      <div
          class="ml-[-22rem] aspect-[1313/771] w-[82.0625rem] flex-none origin-top-right rotate-[30deg] bg-gradient-to-tr from-[#ff80b5] to-[--vp-c-brand-1] xl:ml-0 xl:mr-[calc(50%-12rem)]"
          style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"/>
    </div>
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div class="mx-auto max-w-xl text-center">
        <h2 class="text-lg font-semibold leading-8 tracking-tight text-[--vp-c-brand-2]">Testimonials</h2>
        <p class="mt-2 text-3xl font-bold tracking-tight text-[--vp-c-text-1] sm:text-4xl">See the Minekube experience
          through your amazing words</p>
      </div>
      <div
          class="mx-auto mt-16 grid max-w-2xl grid-cols-1 grid-rows-1 gap-8 text-sm leading-6 text-[--vp-c-text-1] sm:mt-20 sm:grid-cols-2 xl:mx-0 xl:max-w-none xl:grid-flow-col xl:grid-cols-4">
        <!-- This figure is hidden on small screens and visible on screens larger than 'sm' (640px) -->
        <figure
            class="col-span-2 hidden sm:block sm:rounded-2xl sm:bg-[--vp-c-default-soft] sm:shadow-lg sm:ring-1 sm:ring-gray-900/5 xl:col-start-2 xl:row-end-1">
          <blockquote class="p-12 text-xl font-semibold leading-8 tracking-tight text-[--vp-c-text-1]">
            <p>{{ `“${featuredTestimonial.body}”` }}</p>
          </blockquote>
          <figcaption class="flex items-center gap-x-4 border-t border-gray-400/40 px-6 py-4">
            <img :src="featuredTestimonial.author.imageUrl" alt=""
                 class="h-10 w-10 flex-none rounded-full bg-gray-50" loading="lazy"/>
            <div class="flex-auto">
              <div class="font-semibold">{{ featuredTestimonial.author.name }}</div>
              <div class="text-[--vp-c-text-1]">{{ `@${featuredTestimonial.author.handle}` }}</div>
            </div>
            <img :src="featuredTestimonial.author.logoUrl" alt="" class="h-10 w-auto flex-none rounded" loading="lazy"/>
          </figcaption>
        </figure>
        <!-- This figure is visible on small screens and hidden on screens larger than 'sm' (640px) -->
        <figure
            class="col-span-1 block sm:hidden rounded-2xl bg-[--vp-c-default-soft] shadow-lg ring-1 ring-gray-900/5 xl:col-start-2 xl:row-end-1">
          <blockquote class="p-12 text-xl font-semibold leading-8 tracking-tight text-[--vp-c-text-1]">
            <p>{{ `“${featuredTestimonial.body}”` }}</p>
          </blockquote>
          <figcaption class="flex items-center gap-x-4 border-t border-gray-400/40 px-6 py-4">
            <img :src="featuredTestimonial.author.imageUrl" alt=""
                 class="h-10 w-10 flex-none rounded-full bg-gray-50" loading="lazy"/>
            <div class="flex-auto">
              <div class="font-semibold">{{ featuredTestimonial.author.name }}</div>
              <div class="text-[--vp-c-text-1]">{{ `@${featuredTestimonial.author.handle}` }}</div>
            </div>
            <img :src="featuredTestimonial.author.logoUrl" alt="" class="h-10 w-auto flex-none rounded" loading="lazy"/>
          </figcaption>
        </figure>
        <div v-for="(columnGroup, columnGroupIdx) in testimonials" :key="columnGroupIdx"
             class="space-y-8 xl:contents xl:space-y-0">
          <div v-for="(column, columnIdx) in columnGroup" :key="columnIdx"
               :class="[(columnGroupIdx === 0 && columnIdx === 0) || (columnGroupIdx === testimonials.length - 1 && columnIdx === columnGroup.length - 1) ? 'xl:row-span-2' : 'xl:row-start-1', 'space-y-8']">
            <figure v-for="testimonial in column" :key="testimonial.author.handle"
                    class="rounded-2xl bg-[--vp-c-default-soft] p-6 shadow-lg ring-1 ring-gray-900/5">
              <blockquote class="text-[--vp-c-text-1]">
                <p>{{ `“${testimonial.body}”` }}</p>
              </blockquote>
              <figcaption class="mt-6 flex items-center gap-x-4">
                <img :src="testimonial.author.imageUrl" alt="" class="h-10 w-10 rounded-full bg-[--vp-c-text-1]"
                     loading="lazy"/>
                <div>
                  <div class="font-semibold">{{ testimonial.author.name }}</div>
                  <div class="text-[--vp-c-text-1]">{{ `@${testimonial.author.handle}` }}</div>
                </div>
              </figcaption>
            </figure>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const featuredTestimonial = {
  body: "Hey, just wanted to reach out and say hi really! Your Minekube project caught our attention, and it's great to connect with others doing cool stuff in the Minecraft tech space. Would love to have a call to introduce ourselves and hear about your story, plus trade stories of working with Minecraft in weird and wonderful ways",
  author: {
    name: "Zed - Managing Director",
    handle: "cubecraft",
    imageUrl: "https://cdn.discordapp.com/avatars/101763033059041280/e006d30cf2170a60d093730a10ec02b0.webp?size=160",
    logoUrl: "https://avatars.githubusercontent.com/u/6555014?s=200&v=4"
  }
}
const testimonials = [
  [
    [
      {
        body: "It is extremely interesting to see Minekube evolve",
        author: {
          name: "Joshua - CTO",
          handle: "justchunks",
          imageUrl: "https://avatars.githubusercontent.com/u/6044145?s=200&v=4"
        }
      },
      {
        body: "Gradients in Gate 🎉 NBS 😎 added arguments to my NBS command, so I'm just jamming right now. Lol, it even supports server switching mid-song, which I assume wasn't done before. However, Gate has revived my love for Go lol. The SimpleProxy is surprisingly good at covering all bases like scheduling, messages, Brigadier (plus suggestions), and events.",
        author: {
          name: "emortal",
          handle: "emortalmc",
          imageUrl: "https://avatars.githubusercontent.com/u/80778064?s=200&v=4"
        }
      },
      {
        body: "It looks lovely! I have had a good time working on my project so far with Gate and am excited to see its progress.",
        author: {
          name: "Ö for{i}am꧂",
          handle: "bgp179",
          imageUrl: "https://cdn.discordapp.com/avatars/761071532524437565/a764dd33c925c6a2ff8a777c33e50e56.webp?size=160"
        }
      },
      // More testimonials...
    ],
    [
      {
        body: "I’m working on a cloud system in Go and want to support Gate to enhance user experience and make it more extendable. Here's the current version of the Go MC cloud system with a dashboard in action. 😄",
        author: {
          name: "Cerberus - Founder PricelessKeys",
          imageUrl: 'https://cdn.discordapp.com/avatars/539137274554744841/bb326bf193dd2bb3ecfc45980cccf554.webp?size=160',
          handle: 'cerberus.labs'
        }
      }

      // More testimonials...
    ],
  ],
  [
    [
      {
        body: "Thank you, I didn't realize how braindead proxy protocol was, ended up using that. To those looking at a similar situation: it's literally just 'proxyprotocol: true' in paper globals. Enable that on Gate and bada bing bada boom, you're golden.",
        author: {
          name: "Server User",
          imageUrl: 'https://assets-global.website-files.com/6257adef93867e50d84d30e2/636e0a6918e57475a843f59f_icon_clyde_black_RGB.svg',
          handle: 'shit_server'
          // No specific name or image URL provided
        }
      },
      // More testimonials...
    ],
    [
      {
        body: "The instructions for the setup are clear and easy to understand. Both setups I tested are working perfectly, with Gate Lite routing the connection correctly. The custom domain feature is great and works perfectly. I appreciate the ability to manually trigger a refresh for the DNS-check.",
        author: {
          name: "NixNux123 - Admin",
          handle: 'FynnCraft',
          imageUrl: "https://cdn.discordapp.com/avatars/558936206440792076/db90e7439cf81c955d8a8fd0bcd2d109.webp?size=160"
        }
      },
      {
        body: "I like this project so much, I would even be happy to host my own gateway and allow others to use it. I'm sure that people will eventually find Minekube.",
        author: {
          name: "Akos",
          handle: 'theforgtn',
          imageUrl: "https://cdn.discordapp.com/avatars/486530371442180118/9a73f28348432b508160f2a17d0500c5.webp?size=160"
        }
      },
      {
        body: "Thanks for making some cool software Minekube ❤️",
        author: {
          name: "Dylan",
          handle: "dilllxd",
          imageUrl: "https://avatars.githubusercontent.com/u/36942187?v=4"
        }
      }
      // More testimonials...
    ],
  ],
]
</script>