import DefaultTheme from 'vitepress/theme'
import type { Theme } from 'vitepress'
import { h } from 'vue'
import './styles/vars.css'
import VPBadge from 'vitepress/dist/client/theme-default/components/VPBadge.vue'
import MeetTeam from "./components/MeetTeam.vue";

export default {
  extends: DefaultTheme,
  Layout() {
    return h(DefaultTheme.Layout, null, {
      'home-features-after': () => h(MeetTeam),
    })
  },
  enhanceApp({ app }) {
    // app.component('TextAndImage', SvgImage)
    app.component('VPBadge', VPBadge)
    app.component('MeetTeam', MeetTeam)
  }
} satisfies Theme
