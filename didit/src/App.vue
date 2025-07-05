<script setup>
import { onMounted, ref } from "vue";
import Tags from "./components/Tags.vue";
import TopThings from "./components/TopThings.vue";
import Overview from "./components/Overview.vue";

const data = ref([]);

async function fetchData() {
  data.value = await (await fetch("/api/recent")).json();
  console.log('new data fetched')
}

onMounted(async () => {
  fetchData()
});

function padTwoDigits(num) {
  return num.toString().padStart(2, '0');
}

async function onTopThingsDid(what) {
  console.log('onTopThingsDid', what)
  let resp = await fetch(`/api/things`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        what: what,
        year: data.value[0].Year,
        month: data.value[0].Month,
        day: data.value[0].Day
      }),
    });
  fetchData()
}

function onDeleteThing() {
  console.log('onDeleteThing')
  fetchData()
}

function weekday(year, month, day) {
  return new Date(year, month - 1, day).toLocaleDateString(undefined, { weekday: 'long' })
}

</script>

<template>
  <div class="content">
    <div class="overview">
      <Overview :data="data" />
    </div>
    <div class="things">
      <div><TopThings @did="onTopThingsDid" /></div>
      <div class="dates">
        <div class="item" v-for="time of data" :key="time.Year - time.Month - time.Day">
          <div class="date">{{ time.Year }}-{{ padTwoDigits(time.Month) }}-{{ padTwoDigits(time.Day) }} - {{ weekday(time.Year,
            time.Month,
            time.Day) }} </div>
          <Tags :year="time.Year" :month="time.Month" :day="time.Day" :tags="time.Things" @deleted="onDeleteThing"
            @added="fetchData" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.content {
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
}

.things {
  flex-grow: 4;
  margin-left: 0.5em;
}

.date {
  font-weight: bold;
}
.dates {
  margin-top: 1em;
}
.item {
  margin-bottom: 1em;
}
</style>
