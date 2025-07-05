<script setup>
import { ref, defineEmits } from "vue";
import InputText from 'primevue/inputtext';
import Chip from 'primevue/chip';


const emit = defineEmits(["deleted", "added"])

const props = defineProps(["tags", "year", "month", "day"])
const newTag = ref(undefined)

async function addNewTag(what) {
  console.log('adding new tag', what)
  await fetch(`/api/things`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        what: what,
        year: props.year,
        month: props.month,
        day: props.day
      }),
    });
  newTag.value = ''
  emit('added')
}

async function deleteTag(tag) {
  console.log('deleting tag', tag)
  await fetch(`/api/things/${tag.ID}/delete`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });
  emit('deleted')
}
</script>

<template>
  <div v-for="tag in props.tags" :key="tag.ID">
    <Chip class="what" icon :label="tag.What" removable @remove="deleteTag(tag)" />
  </div>
  <InputText v-model.trim="newTag" @keyup.enter="addNewTag(newTag)" placeholder="Did what?" />
</template>

<style scoped>
input {
  max-width: 40em;
}

.what {
  /* margin-bottom: 0.5em; */
}
</style>