<template>
  <div >
    <div class="mx-auto max-w-7xl px-6 py-16 sm:py-24 lg:px-8">
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="text-2xl font-bold leading-10 tracking-tight text-[--vp-c-text-1]">Frequently asked questions</h2>
        <p class="mt-6 text-base leading-7 text-[--vp-c-text-2]">Have a different question and can’t find the answer you’re looking for? Reach out to us on <a href="https://minekube.com/discord" class="hover:text-indigo-30 font-semibold text-indigo-400 hover:text-indigo-300">Discord</a> and we or the community will get back to you as soon as possible.</p>
      </div>
      <div class="mt-20">
        <dl class="space-y-16 sm:grid sm:grid-cols-2 sm:gap-x-6 sm:gap-y-16 sm:space-y-0 lg:gap-x-10">
          <div v-for="faq in faqs" :key="faq.id">
            <dt class="text-xl font-semibold leading-7 text-[--vp-c-text-1]">{{ faq.question }}</dt>
            <dd class="mt-2 text-base leading-7 text-[--vp-c-text-2]">{{ faq.answer }}</dd>
          </div>
        </dl>
      </div>
    </div>
  </div>
</template>

<script setup>
const faqs = [
  {
    id: 0,
    question: "What is your pricing philosophy?",
    // answer: "We hate limits. At Minekube, our core philosophy centers around empowering developers, ensuring that Connect goes out of your way to be a truly empowering platform for creators. We believe in fostering an environment where users can thrive without the constraints of unnecessary limitations, akin to Cloudflare's model. Our free tier is designed to empower the ecosystem and users to grow and innovate, while our paid plan offers enhanced features and customizations for those wanting to opt out of the ecosystem. This approach ensures that the core functionalities remain accessible to everyone, supporting both burgeoning communities and larger. We're committed to creating a happy user base, driven by a robust ecosystem free from stupid limits, where developers feel supported in their creative and technical endeavors."
    answer: "We hate limits. At Minekube, our philosophy is to empower developers by making Connect a platform that transcends ordinary constraints, inspired by Cloudflare's model. Our free tier enables users and the ecosystem to flourish and innovate, while our paid plan provides advanced features for those seeking further customization outside the standard ecosystem. This strategy ensures core functionalities are universally accessible, nurturing both emerging and established communities. Our commitment is to a satisfied community, upheld by an ecosystem that supports developers' creative and technical projects without unnecessary restrictions."
  },
  {
    id: 1,
    question: "What does 'Prevent Fallback to Browser Hub' mean?",
    answer: "This option controls whether players should be redirected to the Browser Hub when the endpoint players are trying to connect to is offline or when a player is disconnected from the endpoint. Enabling this option could lead to players being unable to connect when the endpoint is offline."
  },
  {
    id: 2,
    question: "What does 'Prevent Browser Endpoint Listing' mean?",
    answer: "This option controls whether the endpoint should be listed on the Browser for advertising purposes. Enabling this option could reduce the visibility of the endpoint to potential players."
  },
  {
    id: 3,
    question: "What does 'Prevent /browser Command' mean?",
    answer: "This option controls whether the /browser command should be available for players on the endpoint. Enabling this option could limit the commands available to players on the endpoint."
  },
  {
    id: 4,
    question: "What does 'Prevent Fallback Motd' mean?",
    answer: "This option controls whether the default fallback Message of the Day (Motd) should be displayed to players when the endpoint is offline. Enabling this option could leave players without important information when the endpoint is offline."
  },
  {
    id: 5,
    question: "What does 'Prevent Default Tablist' mean?",
    answer: "This option controls whether a default tablist should be set for players when they connect to or switch to the endpoint. Enabling this option could affect the user interface for players on the endpoint."
  },
  {
    id: 6,
    question: "What are the MOTD Customizations in the paid plan?",
    answer: "Paid users will have access to advanced dashboard functionalities, including the ability to edit MOTDs directly within the dashboard. This feature enhances control over server messages and announcements."
  },
  {
    id: 7,
    question: "Are there any Player Capacity Considerations?",
    answer: "We hate stupid limits. See our philosophy answer earlier. Hence, any limits introduced would be generous to support community growth and development needs."
  },
  {
    id: 8,
    question: "What is the Bedrock Support feature?",
    answer: "Plans to support Bedrock players by enabling a simple option in the dashboard to integrate GeyserMC without management overhead. This would allow Bedrock players to join Java servers, expanding the server's reach and inclusivity."
  },
  {
    id: 9,
    question: "Can I setup a custom domain?",
    answer: "Yes, you can setup a custom domain for your Connect endpoint. This feature is available to all users, including free tier users."
  },
  {
    id: 10,
    question: "Can you implement a feature for me?",
    answer: "We are always open to feedback and suggestions. If you have a feature request, please reach out to us on Discord and we will consider it for future development if it aligns with our vision and we believe it will benefit the whole community."
  },
];
</script>