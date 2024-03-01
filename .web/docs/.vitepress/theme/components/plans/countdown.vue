<template>
  <div>
    <span>{{ days }}d </span>
    <span>{{ hours }}h </span>
    <span>{{ minutes }}m </span>
    <span>{{ seconds }}s </span>
  </div>
</template>

<script>
import {ref, onMounted} from 'vue';

export default {
  props: {
    endTime: {
      type: String,
      required: true
    }
  },
  setup(props) {
    const endTime = new Date(props.endTime).getTime();
    const now = Date.now();
    const timeLeft = ref((endTime - now) / 1000);

    const days = ref(Math.floor(timeLeft.value / (60 * 60 * 24)));
    const hours = ref(Math.floor((timeLeft.value % (60 * 60 * 24)) / (60 * 60)));
    const minutes = ref(Math.floor((timeLeft.value % (60 * 60)) / 60));
    const seconds = ref(Math.floor(timeLeft.value % 60));

    onMounted(() => {
      const interval = setInterval(() => {
        timeLeft.value -= 1;
        days.value = Math.floor(timeLeft.value / (60 * 60 * 24));
        hours.value = Math.floor((timeLeft.value % (60 * 60 * 24)) / (60 * 60));
        minutes.value = Math.floor((timeLeft.value % (60 * 60)) / 60);
        seconds.value = Math.floor(timeLeft.value % 60);
      }, 1000);
      return () => {
        clearInterval(interval);
      }
    });

    return { days, hours, minutes, seconds };
  }
};
</script>