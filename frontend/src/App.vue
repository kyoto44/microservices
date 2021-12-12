<script setup>
import { ref } from "vue";

const logsTable = ref({ country: "stub", ip: "stub", useragent: "stub" });
const getAll = () => {
  fetch("http://localhost:8000/api/getall")
    .then((response) => response.json())
    .then((json) => {
      logsTable.value = json;
      console.log(json);
    });
};

const getRandom = () => {
  fetch("http://localhost:8000/api/getrandom")
    .then((response) => response.json())
    .then((json) => {
      logsTable.value = {};
      logsTable.value[0] = json;
      console.log(json);
    });
};

const deleteAll = () => {
  fetch("http://localhost:8000/api/deleteall").then((response) => {
    logsTable.value = { country: "", ip: "", useragent: "" };
    console.log(response);
  });
};
</script>

<template>
  <body>
    <nav>
      <button
        type="button"
        class="bg-sky-600 hover:bg-sky-700 rounded-full"
        @click="getAll()"
      >
        Get All Entries
      </button>
      <button
        type="button"
        class="bg-sky-600 hover:bg-sky-700 rounded-full"
        @click="getRandom()"
      >
        Get Random Entry
      </button>
      <button
        type="button"
        class="bg-sky-600 hover:bg-sky-700 rounded-full"
        @click="deleteAll()"
      >
        Delete All Entries from DB
      </button>
    </nav>

    <div>
      <table class="shadow-lg bg-white table-auto">
        <thead class="bg-gray-50">
          <tr>
            <th class="bg-blue-100 border text-left px-8 py-4">country</th>
            <th class="bg-blue-100 border text-left px-8 py-4">ip</th>
            <th class="bg-blue-100 border text-left px-8 py-4">useragent</th>
          </tr>
        </thead>
        <tbody>
          <tr
            class="odd:bg-white even:bg-gray-100"
            v-for="(row, i) in logsTable"
            :key="i"
          >
            <td class="border px-8 py-4">{{ row.country }}</td>
            <td class="border px-8 py-4">{{ row.ip }}</td>
            <td class="border px-8 py-4">{{ row.useragent }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </body>
</template>

<style scoped>
a {
  color: #42b983;
}

body {
  background: #eee;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

button {
  margin: 10px;
}
</style>
