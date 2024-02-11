import './tailwind.postcss'
import DefaultTheme from 'vitepress/theme'
import type {Theme} from 'vitepress'
import './styles/vars.css'
import VPBadge from 'vitepress/dist/client/theme-default/components/VPBadge.vue'
import MeetTeam from "./components/MeetTeam.vue";
import Layout from "./components/Layout.vue";
import PlansLanding from "./components/plans/PlansLanding.vue";

export default {
  extends: DefaultTheme,
  Layout: Layout,
  enhanceApp({ app }) {
    // app.component('TextAndImage', SvgImage)
    app.component('VPBadge', VPBadge)
    app.component('MeetTeam', MeetTeam)
    app.component('PlansLanding', PlansLanding)
  }
} satisfies Theme
