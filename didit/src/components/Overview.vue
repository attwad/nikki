<script setup>

const props = defineProps(["data"])

function isMonday(year, month, day) {
  return new Date(year, month - 1, day).getDay() === 0
}

</script>

<template>
  <div class="item" v-if="props.data" v-for="time of props.data" :key="time.Year - time.Month - time.Day">
    <hr v-if="isMonday(time.Year, time.Month, time.Day)"></hr>
    <span class="thing" v-if="time.Things" v-for="thing of time.Things" :key="thing.ID"></span>
    <span v-else class="nothing"></span>
  </div>
</template>

<style scoped>
.thing {
  height: 2px;
  width: 2px;
  background-color: #00ff2f;
  display: inline-block;
  margin-right: 0.2em;
  margin-bottom: 0.2em;
}
.nothing {
  height: 2px;
  width: 2px;
  background-color: #ff0000;
  display: inline-block;
  margin-right: 0.2em;
  margin-bottom: 0.2em;
}
</style>