<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue';
import type { GlobeInstance } from 'globe.gl';

const props = defineProps<{
  class?: string;
  locations?: Array<{
    lat: number;
    lng: number;
    label?: string;
    color?: string;
    size?: number;
  }>;
  arcs?: Array<{
    startLat: number;
    startLng: number;
    endLat: number;
    endLng: number;
    color?: string;
    label?: string;
  }>;
  width?: number;
  height?: number;
  globeImageUrl?: string;
  atmosphereColor?: string;
  atmosphereAltitude?: number;
  pointColor?: string;
  pointAltitude?: number;
  pointRadius?: number;
  arcColor?: string;
  arcAltitude?: number | null;
  arcAltitudeAutoScale?: number;
  arcStroke?: number | null;
  arcDashLength?: number;
  arcDashGap?: number;
  arcDashAnimateTime?: number;
}>();

const globeContainer = ref<HTMLDivElement | null>(null);
const globeInstance = ref<GlobeInstance | null>(null);

const globeWidth = props.width ?? 600;
const globeHeight = props.height ?? 600;

onMounted(async () => {
  // Dynamically import Globe only on the client-side
  const Globe = (await import('globe.gl')).default;
  if (!Globe) {
    console.error('Failed to load globe.gl');
    return;
  }

  await nextTick(); // Ensure the DOM is ready

  if (globeContainer.value) {
    const instance = Globe()(globeContainer.value)
      .width(globeWidth)
      .height(globeHeight)
      .globeImageUrl(
        props.globeImageUrl ||
          '//unpkg.com/three-globe/example/img/earth-dark.jpg'
      )
      .backgroundColor('rgba(0,0,0,0)') // Transparent background
      .showAtmosphere(true)
      .atmosphereColor(props.atmosphereColor || 'lightskyblue')
      .atmosphereAltitude(props.atmosphereAltitude || 0.15);

    // Auto-rotate
    instance.controls().autoRotate = true;
    instance.controls().autoRotateSpeed = 0.5;

    // Locations from fly platform regions --json
    const connectLocations = [
      {
        lat: 39.02214,
        lng: -77.462555,
        label: 'Ashburn, Virginia (US)',
        color: '#4285F4',
        size: 0.05,
      }, // iad
      {
        lat: 51.516434,
        lng: -0.125656,
        label: 'London, United Kingdom',
        color: '#DB4437',
        size: 0.05,
      }, // lhr
      {
        lat: 20.5217,
        lng: -103.3109,
        label: 'Guadalajara, Mexico',
        color: '#F4B400',
        size: 0.05,
      }, // gdl
      {
        lat: 52.1657,
        lng: 20.9671,
        label: 'Warsaw, Poland',
        color: '#0F9D58',
        size: 0.05,
      }, // waw
      {
        lat: 33.9416,
        lng: -118.4085,
        label: 'Los Angeles, California (US)',
        color: '#4285F4',
        size: 0.05,
      }, // lax
      {
        lat: 52.374344,
        lng: 4.895439,
        label: 'Amsterdam, Netherlands',
        color: '#DB4437',
        size: 0.05,
      }, // ams
      {
        lat: 35.62161,
        lng: 139.74185,
        label: 'Tokyo, Japan',
        color: '#F4B400',
        size: 0.06,
      }, // nrt
      {
        lat: -23.549664,
        lng: -46.65435,
        label: 'Sao Paulo, Brazil',
        color: '#0F9D58',
        size: 0.05,
      }, // gru
      {
        lat: 44.4325,
        lng: 26.1039,
        label: 'Bucharest, Romania',
        color: '#4285F4',
        size: 0.05,
      }, // otp
      {
        lat: 22.25097,
        lng: 114.203224,
        label: 'Hong Kong, Hong Kong',
        color: '#DB4437',
        size: 0.05,
      }, // hkg
      {
        lat: 1.3,
        lng: 103.8,
        label: 'Singapore, Singapore',
        color: '#F4B400',
        size: 0.05,
      }, // sin
    ];

    // Define a more comprehensive set of arcs connecting your locations
    const connectArcs = [
      // North America to Europe & South America
      {
        startLat: 39.02214, // Ashburn (iad)
        startLng: -77.462555,
        endLat: 51.516434, // London (lhr)
        endLng: -0.125656,
        color: 'rgba(76, 175, 80, 0.7)', // Greenish
        label: 'Ashburn to London',
      },
      {
        startLat: 39.02214, // Ashburn (iad)
        startLng: -77.462555,
        endLat: -23.549664, // Sao Paulo (gru)
        endLng: -46.65435,
        color: 'rgba(33, 150, 243, 0.7)', // Bluish
        label: 'Ashburn to Sao Paulo',
      },
      // West Coast US to Asia & Mexico
      {
        startLat: 33.9416, // Los Angeles (lax)
        startLng: -118.4085,
        endLat: 35.62161, // Tokyo (nrt)
        endLng: 139.74185,
        color: 'rgba(255, 193, 7, 0.7)', // Amber
        label: 'LA to Tokyo',
      },
      {
        startLat: 33.9416, // Los Angeles (lax)
        startLng: -118.4085,
        endLat: 20.5217, // Guadalajara (gdl)
        endLng: -103.3109,
        color: 'rgba(255, 87, 34, 0.7)', // Deep Orange
        label: 'LA to Guadalajara',
      },
      // Intra-Europe
      {
        startLat: 51.516434, // London (lhr)
        startLng: -0.125656,
        endLat: 52.374344, // Amsterdam (ams)
        endLng: 4.895439,
        color: 'rgba(156, 39, 176, 0.7)', // Purple
        label: 'London to Amsterdam',
      },
      {
        startLat: 51.516434, // London (lhr)
        startLng: -0.125656,
        endLat: 52.1657, // Warsaw (waw)
        endLng: 20.9671,
        color: 'rgba(233, 30, 99, 0.7)', // Pink
        label: 'London to Warsaw',
      },
      {
        startLat: 52.374344, // Amsterdam (ams)
        startLng: 4.895439,
        endLat: 44.4325, // Bucharest (otp)
        endLng: 26.1039,
        color: 'rgba(0, 188, 212, 0.7)', // Cyan
        label: 'Amsterdam to Bucharest',
      },
      // Intra-Asia & Asia-Europe
      {
        startLat: 35.62161, // Tokyo (nrt)
        startLng: 139.74185,
        endLat: 22.25097, // Hong Kong (hkg)
        endLng: 114.203224,
        color: 'rgba(255, 235, 59, 0.7)', // Yellow
        label: 'Tokyo to Hong Kong',
      },
      {
        startLat: 22.25097, // Hong Kong (hkg)
        startLng: 114.203224,
        endLat: 1.3, // Singapore (sin)
        endLng: 103.8,
        color: 'rgba(121, 85, 72, 0.7)', // Brown
        label: 'Hong Kong to Singapore',
      },
      {
        startLat: 1.3, // Singapore (sin)
        startLng: 103.8,
        endLat: 52.1657, // Warsaw (waw)
        endLng: 20.9671,
        color: 'rgba(96, 125, 139, 0.7)', // Blue Grey
        label: 'Singapore to Warsaw',
      },
      // Intra-South America
      {
        startLat: -23.549664, // Sao Paulo (gru)
        startLng: -46.65435,
        endLat: 20.5217, // Guadalajara (gdl)
        endLng: -103.3109,
        color: 'rgba(255, 152, 0, 0.7)', // Orange
        label: 'Sao Paulo to Guadalajara',
      },
    ];

    instance
      .pointsData(props.locations || connectLocations)
      .pointColor((d) => (d as any).color || props.pointColor || 'white')
      .pointAltitude(props.pointAltitude || 0.01)
      .pointRadius((d) => (d as any).size || props.pointRadius || 0.03)
      .pointLabel((d) => (d as any).label || '');

    instance
      .arcsData(props.arcs || connectArcs)
      .arcColor(
        (d) => (d as any).color || props.arcColor || 'rgba(255, 255, 255, 0.5)'
      )
      .arcAltitude(props.arcAltitude === undefined ? null : props.arcAltitude) // Allow auto-scaling by default
      .arcAltitudeAutoScale(props.arcAltitudeAutoScale || 0.5)
      .arcStroke(props.arcStroke === undefined ? null : props.arcStroke) // Thin line by default
      .arcDashLength(props.arcDashLength || 0.9)
      .arcDashGap(props.arcDashGap || 0.2)
      .arcDashAnimateTime(props.arcDashAnimateTime || 2000)
      .arcLabel((d) => (d as any).label || '');

    globeInstance.value = instance;

    // Adapt to container size
    const resizeObserver = new ResizeObserver(() => {
      const newWidth = globeContainer.value?.clientWidth || globeWidth;
      const newHeight = globeContainer.value?.clientHeight || globeHeight;
      instance.width(newWidth);
      instance.height(newHeight);
    });
    if (globeContainer.value) {
      resizeObserver.observe(globeContainer.value);
    }

    onUnmounted(() => {
      if (instance.controls()) {
        instance.controls().dispose();
      }
      instance._destructor(); // globe.gl's way to clean up
      if (globeContainer.value) {
        resizeObserver.unobserve(globeContainer.value);
      }
    });
  }
});

// Watch for prop changes to update the globe
watch(
  () => props.locations,
  (newLocations) => {
    if (globeInstance.value) {
      globeInstance.value.pointsData(newLocations || []);
    }
  }
);

watch(
  () => props.arcs,
  (newArcs) => {
    if (globeInstance.value) {
      globeInstance.value.arcsData(newArcs || []);
    }
  }
);

watch(
  () => props.globeImageUrl,
  (newVal) => {
    if (globeInstance.value && newVal) {
      globeInstance.value.globeImageUrl(newVal);
    }
  }
);

// Add more watchers for other props as needed, for example:
// watch(() => props.atmosphereColor, (newVal) => {
//   if (globeInstance.value && newVal) {
//     // globe.gl doesn't have a direct method to update atmosphereColor after init
//     // Might need to re-initialize or directly manipulate Three.js objects if critical
//     // For now, we'll note that some visual props might require a re-mount or deeper Three.js manipulation
//     console.warn("Changing atmosphereColor after initialization might require a re-mount or direct Three.js manipulation.");
//   }
// });
</script>

<template>
  <ClientOnly>
    <div
      ref="globeContainer"
      :class="props.class"
      :style="{
        width: globeWidth + 'px',
        height: globeHeight + 'px',
        maxWidth: '100%',
        aspectRatio: '1 / 1',
      }"
    ></div>
  </ClientOnly>
</template>
