import './tailwind.postcss'
import DefaultTheme from 'vitepress/theme'
import type {Theme} from 'vitepress'
import './styles/vars.css'
import './styles/locationIcons.css'
import VPBadge from 'vitepress/dist/client/theme-default/components/VPBadge.vue'
import MeetTeam from "./components/MeetTeam.vue";
import Layout from "./components/Layout.vue";
import PlansLanding from "./components/plans/PlansLanding.vue";
import PostLayout from "./components/posts/Layout.vue";
import PostHome from "./components/posts/Home.vue";

export default {
  extends: DefaultTheme,
  Layout: Layout,
  enhanceApp({ app }) {
    // app.component('TextAndImage', SvgImage)
    app.component('VPBadge', VPBadge)
    app.component('MeetTeam', MeetTeam)
    app.component('PlansLanding', PlansLanding)
    app.component('PostHome', PostHome)
    app.component('Post', PostLayout)
  }
} satisfies Theme
