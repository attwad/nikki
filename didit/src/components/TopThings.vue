<script setup>
import { ref, onMounted, defineEmits } from "vue";
defineEmits(['did'])

const oftenDoneThings = ref([])

onMounted(async () => {
  oftenDoneThings.value = await (await fetch("/api/things/oftendone")).json();
})

</script>

<template>
  <div class="title">Often done</div>
  <Button v-for="thing in oftenDoneThings" :key="thing" variant="outlined" class="btn" :label="thing"
    icon="pi pi-angle-double-down" iconPos="bottom" @click="$emit('did', thing)" />
</template>

<style scoped>
.btn {
  margin-right: 0.5em;
  margin-bottom: 0.5em;
}
.title {
  font-weight: bold;
  font-size:larger;
  margin-bottom: 0.5em;
}
</style>