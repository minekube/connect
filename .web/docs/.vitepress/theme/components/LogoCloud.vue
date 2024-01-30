<template>
  <div class="py-24 sm:py-32">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div class="grid grid-cols-1 items-center gap-x-8 gap-y-16 lg:grid-cols-2">
        <div class="mx-auto w-full max-w-xl lg:mx-0">
          <h2 class="text-3xl font-bold tracking-tight text-[--vp-c-text-1]">Use your favorite Minecraft software</h2>
          <p class="mt-6 text-lg leading-8 text-[--vp-c-text-2]">
            Connect works with all major Minecraft server software. Whether you are using Spigot, Paper, Gate, Velocity
            or BungeeCord.
            Install the Connect plugin or enable native Connect in the Gate config and you are ready to go.
          </p>
          <div class="mt-8 flex items-center gap-x-6">
            <a class="rounded-md bg-[--vp-button-brand-bg] px-3.5 py-2.5 text-sm font-semibold text-[--vp-button-brand-text] shadow-sm hover:bg-[--vp-button-brand-hover-bg] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[--vp-button-brand-active-bg]"
               href="https://app.minekube.com/">Dashboard
              App</a>
            <a class="text-sm font-semibold text-[--vp-button-brand-text]" href="https://minekube.com/discord">Discord
              <span aria-hidden="true">&rarr;</span></a>
          </div>
        </div>
        <div
            class="mx-auto grid w-full max-w-xl grid-cols-2 items-center gap-y-12 sm:gap-y-14 lg:mx-0 lg:max-w-none lg:pl-8">
          <a href="/guide/connectors/gate">
            <img alt="Gate" class="max-h-24 w-full object-contain object-left" loading="lazy" src="/minekube-logo.png"/>
          </a>
          <a href="/guide/connectors/plugin">
            <img alt="Paper" class="max-h-24 w-full object-contain object-left" loading="lazy"
                 src="https://avatars.githubusercontent.com/u/7608950?s=200&v=4"/>
          </a>
          <a href="/guide/connectors/plugin">
            <img alt="Spigot" class="max-h-24 w-full object-contain object-left" loading="lazy"
                 src="https://avatars.githubusercontent.com/u/4350249?s=48&v=4"/>
          </a>
          <a href="/guide/connectors/plugin">
            <img alt="Velocity" class="max-h-24 w-full object-contain object-left" loading="lazy"
                 src="https://avatars.githubusercontent.com/u/41710604?s=200&v=4"/>
          </a>
        </div>
      </div>
    </div>
  </div>
</template>