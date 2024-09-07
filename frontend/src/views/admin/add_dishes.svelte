<script lang="ts">
  import { onMount } from "svelte";
  import DishPreview from "./dish_preview.svelte";

  let dish = {
    name: "",
    categories: "",
    url: "",
    price: 0,
  };
  var files;
  function fileInput(filesInput) {
    files = filesInput;
    dish.url = window.URL.createObjectURL(files[0]);
  }

//  let categoriesList = []
onMount(()=>{
//     categoriesList = 
})
</script>

<h3>Добавить блюдо</h3>
<div class="input_box">
  <label for="dish_image">Upload a picture:</label>
  <input
    id="dish_image"
    accept="image/png, image/jpeg"
    on:input={(filesInput) => fileInput(filesInput)}
    type="file"
    class="dish_image_input"
  />
  <label for="dish_name">Название блюда:</label>
  <input id="dish_name" bind:value={dish.name} type="text" />
  <label for="dish_price">Цена блюда:</label>
  <input
    id="dish_price"
    bind:value={dish.price}
    type="number"
    min="1"
    max="9999999"
  />
</div>

<h3>Предосмотр:</h3>
{#if dish.url !== "" && dish.name !== "" && dish.price !== 0 && dish.categories !== ""}
  <DishPreview bind:dish />
{/if}

<style>
  .input_box {
    width: 100vw;
    display: block;
  }
</style>
